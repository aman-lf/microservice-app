package main

import (
        "database/sql"
        "errors"
        "fmt"
        "log"
        "net/http"
        "os"
        "strings"
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

// setupRoutes is now defined in routes.go

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
        return strings.Contains(s, substr)
}

// Connect to database
func connectToDB() *sql.DB {
        // Get the database URL from environment
        dbURL := os.Getenv("DATABASE_URL")
        
        // Extract the endpoint ID from the hostname for Neon PostgreSQL
        if dbURL != "" {
                // Check if we're using Neon PostgreSQL (contains neon.tech)
                if strings.Contains(dbURL, "neon.tech") {
                        // Extract the hostname from the URL
                        // Format: postgresql://user:pass@ep-name-id.region.aws.neon.tech/dbname
                        startIndex := strings.Index(dbURL, "@") + 1
                        endIndex := strings.Index(dbURL[startIndex:], "/")
                        if startIndex > 0 && endIndex > 0 {
                                hostname := dbURL[startIndex : startIndex+endIndex]
                                
                                // Extract the endpoint ID from the hostname (ep-name-id)
                                epParts := strings.Split(hostname, ".")
                                if len(epParts) > 0 {
                                        epID := epParts[0]
                                        
                                        // Add the endpoint parameter to the URL
                                        if strings.Contains(dbURL, "?") {
                                                dbURL = dbURL + "&options=endpoint%3D" + epID
                                        } else {
                                                dbURL = dbURL + "?options=endpoint%3D" + epID
                                        }
                                        
                                        log.Printf("Using Neon PostgreSQL with endpoint ID: %s", epID)
                                }
                        }
                }
        }
        
        dsn := dbURL

        for {
                connection, err := openDB(dsn)
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
