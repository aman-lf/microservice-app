package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) setupRoutes() {
	app.router.GET("/orders", app.GetAllOrders)
	app.router.GET("/orders/:id", app.GetOrder)
	app.router.GET("/orders/customer/:customer_id", app.GetOrdersByCustomer)
	app.router.POST("/orders", app.CreateOrder)
	app.router.PATCH("/orders/:id/status", app.UpdateOrderStatus)
	
	// Health check endpoint
	app.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Order service is running",
		})
	})
}
