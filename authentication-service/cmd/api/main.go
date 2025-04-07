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

const webPort = "8001"

var counts int64

type Config struct {
	DB     *sql.DB
	router *gin.Engine
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	// Connect to database
	log.Println("Starting authentication service")

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
	log.Printf("Starting authentication service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", webPort),
		Handler: app.router,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

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
