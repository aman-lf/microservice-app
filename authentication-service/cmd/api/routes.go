package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) setupRoutes() {
	app.router.POST("/authenticate", app.Authenticate)
	app.router.POST("/user", app.CreateUser)
	app.router.GET("/users", app.GetAllUsers)
	
	// Health check endpoint
	app.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication service is running",
		})
	})
}
