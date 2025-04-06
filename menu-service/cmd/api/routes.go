package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) setupRoutes() {
	app.router.GET("/menu", app.GetAllMenuItems)
	app.router.GET("/menu/:id", app.GetMenuItem)
	app.router.POST("/menu", app.CreateMenuItem)
	app.router.PUT("/menu/:id", app.UpdateMenuItem)
	app.router.DELETE("/menu/:id", app.DeleteMenuItem)
	
	// Health check endpoint
	app.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Menu service is running",
		})
	})
}
