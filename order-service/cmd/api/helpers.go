package main

import (
        "encoding/json"
        "errors"
        "io"
        "net/http"
)

type jsonResponse struct {
        Error   bool   `json:"error"`
        Message string `json:"message"`
        Data    any    `json:"data,omitempty"`
}

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) readJSONRaw(r *http.Request, data any) error {
        maxBytes := 1048576 // one megabyte

        // Create a new reader with a size limit
        limitedReader := io.LimitReader(r.Body, int64(maxBytes))

        dec := json.NewDecoder(limitedReader)
        err := dec.Decode(data)
        if err != nil {
                return err
        }

        err = dec.Decode(&struct{}{})
        if err != io.EOF {
                return errors.New("body must have only a single JSON value")
        }

        return nil
}

// errorJSONRaw takes an error, and optionally a status code, and generates and sends
// a JSON error response
func (app *Config) errorJSONRaw(w http.ResponseWriter, err error, status ...int) error {
        statusCode := http.StatusBadRequest

        if len(status) > 0 {
                statusCode = status[0]
        }

        var payload jsonResponse
        payload.Error = true
        payload.Message = err.Error()

        return app.writeJSONRaw(w, statusCode, payload)
}

// writeJSONRaw takes a response status code and arbitrary data and writes a json response to the client
func (app *Config) writeJSONRaw(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
        out, err := json.Marshal(data)
        if err != nil {
                return err
        }

        if len(headers) > 0 {
                for key, value := range headers[0] {
                        w.Header()[key] = value
                }
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        _, err = w.Write(out)
        if err != nil {
                return err
        }

        return nil
}

// logRequest sends a log entry to the logger service
func (app *Config) logRequest(name, data string) error {
        // In a real application, this would send a request to the logger service
        // For now, we'll just log to the console
        return nil
}
