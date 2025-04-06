package main

import (
        "github.com/gin-gonic/gin"
        "net/http"
)

type jsonResponse struct {
        Error   bool   `json:"error"`
        Message string `json:"message"`
        Data    any    `json:"data,omitempty"`
}

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) readJSON(c *gin.Context, data any) error {
        if err := c.ShouldBindJSON(data); err != nil {
                return err
        }
        return nil
}

// errorJSON takes an error, and optionally a status code, and generates and sends
// a JSON error response
func (app *Config) errorJSON(c *gin.Context, err error, status ...int) {
        statusCode := http.StatusBadRequest

        if len(status) > 0 {
                statusCode = status[0]
        }

        var payload jsonResponse
        payload.Error = true
        payload.Message = err.Error()

        c.JSON(statusCode, payload)
}

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
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
