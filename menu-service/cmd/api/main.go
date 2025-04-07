package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const webPort = "8002"

var counts int64

type Config struct {
	DB     *sql.DB
	router *gin.Engine
}

type MenuItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func main() {
	log.Println("Starting menu service")

	// Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Set up application config
	app := Config{
		DB: conn,
	}

	// Set up Gin router with middleware
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	app.router = router

	// Set up routes
	app.setupRoutes()

	// Start the server
	log.Printf("Starting menu service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", webPort),
		Handler: app.router,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

// writeJSON is a helper to write JSON responses
func (app *Config) writeJSON(c *gin.Context, status int, data any, headers ...http.Header) error {
	// Add headers if they exist
	if len(headers) > 0 {
		for key, value := range headers[0] {
			for _, v := range value {
				c.Header(key, v)
			}
		}
	}

	c.JSON(status, data)
	return nil
}

// readJSON tries to read the body of a request and convert it into JSON
func (app *Config) readJSON(c *gin.Context, data any) error {
	if err := c.ShouldBindJSON(data); err != nil {
		return err
	}

	return nil
}

// errorJSON takes an error, and optionally a status code, and generates and sends a JSON error response
func (app *Config) errorJSON(c *gin.Context, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = err.Error()

	c.JSON(statusCode, payload)
}

// setupRoutes is now defined in routes.go

// Connect to database
func connectToDB() *sql.DB {
	// Get the database URL from environment
	dbURL := os.Getenv("DATABASE_URL")

	for {
		connection, err := openDB(dbURL)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")

			// Initialize the database tables
			err = initDB(connection)
			if err != nil {
				log.Printf("Error initializing database: %v", err)
				return nil
			}

			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
