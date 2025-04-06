package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) setupRoutes() {
	app.router.GET("/inventory", app.GetAllInventoryItems)
	app.router.GET("/inventory/:id", app.GetInventoryItem)
	app.router.POST("/inventory", app.CreateInventoryItem)
	app.router.PUT("/inventory/:id", app.UpdateInventoryItem)
	app.router.DELETE("/inventory/:id", app.DeleteInventoryItem)
	app.router.PATCH("/inventory/:id/adjust", app.AdjustInventory)
	app.router.GET("/inventory/low", app.CheckLowInventory)
	
	// Health check endpoint
	app.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Inventory service is running",
		})
	})
}
