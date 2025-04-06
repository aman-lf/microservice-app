package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(c *gin.Context) {
	var requestPayload JSONPayload

	err := app.readJSON(c, &requestPayload)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	// Insert the log entry
	err = app.Models.LogEntry.Insert(LogEntry{
		Name:      requestPayload.Name,
		Data:      requestPayload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		app.errorJSON(c, err)
		return
	}

	// Send a response
	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Logged",
	}

	app.writeJSON(c, http.StatusAccepted, payload)
}

func (app *Config) GetAllLogs(c *gin.Context) {
	logs, err := app.Models.LogEntry.GetAll()
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	payload := struct {
		Error   bool       `json:"error"`
		Message string     `json:"message"`
		Data    []LogEntry `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Logs retrieved",
		Data:    logs,
	}

	app.writeJSON(c, http.StatusOK, payload)
}
