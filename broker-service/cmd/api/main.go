package main

import (
        "errors"
        "fmt"
        "log"
        "net/http"

        "github.com/gin-contrib/cors"
        "github.com/gin-gonic/gin"
)

const (
        webPort = "8000"
)

type Config struct {
        router *gin.Engine
}

func main() {
        // Set Gin to release mode in production
        gin.SetMode(gin.ReleaseMode)

        // Create a new Gin router
        router := gin.New()

        // Add middleware
        router.Use(gin.Recovery())
        router.Use(cors.New(cors.Config{
                AllowAllOrigins:  true,
                AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
                AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
                ExposeHeaders:    []string{"Link"},
                AllowCredentials: true,
                MaxAge:           300,
        }))

        // Create app config
        app := Config{
                router: router,
        }

        // Define routes
        app.routes()

        // Start the server
        log.Printf("Starting broker service on port %s\n", webPort)
        srv := &http.Server{
                Addr:    fmt.Sprintf("0.0.0.0:%s", webPort),
                Handler: router,
        }

        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
                log.Fatalf("Failed to listen and serve: %v", err)
        }
}
