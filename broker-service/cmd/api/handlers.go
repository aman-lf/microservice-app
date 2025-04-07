package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestPayload is the structure that defines the data sent to the broker
type RequestPayload struct {
	Action    string           `json:"action"`
	Auth      AuthPayload      `json:"auth,omitempty"`
	Menu      MenuPayload      `json:"menu,omitempty"`
	Order     OrderPayload     `json:"order,omitempty"`
	Inventory InventoryPayload `json:"inventory,omitempty"`
	Log       LogPayload       `json:"log,omitempty"`
}

// AuthPayload is the data needed for authentication
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// MenuPayload is the data needed for menu operations
type MenuPayload struct {
	ID          int     `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Category    string  `json:"category,omitempty"`
}

// OrderPayload is the data needed for order operations
type OrderPayload struct {
	ID         int         `json:"id,omitempty"`
	CustomerID int         `json:"customer_id,omitempty"`
	Items      []OrderItem `json:"items,omitempty"`
	Status     string      `json:"status,omitempty"`
	Total      float64     `json:"total,omitempty"`
	CreatedAt  string      `json:"created_at,omitempty"`
}

type OrderItem struct {
	MenuItemID int     `json:"menu_item_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

// InventoryPayload is the data needed for inventory operations
type InventoryPayload struct {
	ID        int    `json:"id,omitempty"`
	ItemName  string `json:"item_name,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
	Unit      string `json:"unit,omitempty"`
	Threshold int    `json:"threshold,omitempty"`
}

// LogPayload is the data needed for logging
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Broker handles all incoming requests and routes them to the appropriate service
func (app *Config) Broker(c *gin.Context) {
	payload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "Hit the broker",
	}

	_ = app.writeJSON(c, http.StatusOK, payload)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(c *gin.Context) {
	var requestPayload RequestPayload

	err := app.readJSON(c, &requestPayload)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(c, requestPayload.Auth)
	case "menu":
		app.handleMenuRequest(c, requestPayload.Menu)
	case "order":
		app.handleOrderRequest(c, requestPayload.Order)
	case "inventory":
		app.handleInventoryRequest(c, requestPayload.Inventory)
	case "log":
		app.logItem(c, requestPayload.Log)
	default:
		app.errorJSON(c, errors.New("unknown action"))
	}
}

// authenticate calls the authentication service
func (app *Config) authenticate(c *gin.Context, payload AuthPayload) {
	// Create JSON to send to auth service
	jsonData, _ := json.Marshal(payload)

	// Call the service
	request, err := http.NewRequest("POST", "http://0.0.0.0:8001/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer response.Body.Close()

	// Make sure we get the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(c, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJSON(c, errors.New("error calling auth service"), http.StatusInternalServerError)
		return
	}

	// Read response.Body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	// Send JSON back to the client
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(responseBody)
}

// handleMenuRequest calls the menu service
func (app *Config) handleMenuRequest(c *gin.Context, payload MenuPayload) {
	jsonData, _ := json.Marshal(payload)

	// Call the service
	request, err := http.NewRequest("POST", "http://0.0.0.0:8002/menu", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer response.Body.Close()

	// Make sure we get the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(c, errors.New("error calling menu service"), http.StatusInternalServerError)
		return
	}

	// Read response.Body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	// Send JSON back to the client
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(responseBody)
}

// handleOrderRequest calls the order service
func (app *Config) handleOrderRequest(c *gin.Context, payload OrderPayload) {
	jsonData, _ := json.Marshal(payload)

	// Call the service
	request, err := http.NewRequest("POST", "http://0.0.0.0:8004/order", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer response.Body.Close()

	// Make sure we get the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(c, errors.New("error calling order service"), http.StatusInternalServerError)
		return
	}

	// Read response.Body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	// Send JSON back to the client
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(responseBody)
}

// handleInventoryRequest calls the inventory service
func (app *Config) handleInventoryRequest(c *gin.Context, payload InventoryPayload) {
	jsonData, _ := json.Marshal(payload)

	// Call the service
	request, err := http.NewRequest("POST", "http://0.0.0.0:8003/inventory", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer response.Body.Close()

	// Make sure we get the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(c, errors.New("error calling inventory service"), http.StatusInternalServerError)
		return
	}

	// Read response.Body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	// Send JSON back to the client
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(responseBody)
}

// logItem logs an event using the logger-service (via RPC)
func (app *Config) logItem(c *gin.Context, payload LogPayload) {
	jsonData, _ := json.Marshal(payload)

	// Call the service
	request, err := http.NewRequest("POST", "http://0.0.0.0:8005/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer response.Body.Close()

	// Make sure we get the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(c, errors.New("error calling logger service"), http.StatusInternalServerError)
		return
	}

	// Send response back to the client
	responsePayload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "Log entry created",
	}

	_ = app.writeJSON(c, http.StatusOK, responsePayload)
}
