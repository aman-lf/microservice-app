package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) routes() {
	app.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Broker service is running",
		})
	})

	app.router.POST("/", app.Broker)
	app.router.POST("/handle", app.HandleSubmission)
}
