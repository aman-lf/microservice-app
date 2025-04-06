package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *Config) GetAllOrders(c *gin.Context) {
	orders, err := app.getAllOrders()
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool    `json:"error"`
		Message string  `json:"message"`
		Data    []Order `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Orders retrieved",
		Data:    orders,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	order, err := app.getOrderByID(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    Order  `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Order retrieved",
		Data:    order,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) GetOrdersByCustomer(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("customer_id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid customer_id parameter"), http.StatusBadRequest)
		return
	}

	orders, err := app.getOrdersByCustomer(customerID)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool    `json:"error"`
		Message string  `json:"message"`
		Data    []Order `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Customer orders retrieved",
		Data:    orders,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) CreateOrder(c *gin.Context) {
	var order Order

	err := app.readJSON(c, &order)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Validate the order data
	if order.CustomerID == 0 || len(order.Items) == 0 {
		app.errorJSON(c, errors.New("invalid order data"), http.StatusBadRequest)
		return
	}

	// Create the order
	newID, err := app.insertOrder(order)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	// Get the newly created order with items
	newOrder, err := app.getOrderByID(newID)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    Order  `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Order created",
		Data:    newOrder,
	}

	app.writeJSON(c, http.StatusCreated, payload)
}

func (app *Config) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}

	err = app.readJSON(c, &statusUpdate)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Update the order status
	err = app.updateOrderStatus(id, statusUpdate.Status)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	// Get the updated order
	updatedOrder, err := app.getOrderByID(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    Order  `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Order status updated",
		Data:    updatedOrder,
	}

	app.writeJSON(c, http.StatusOK, payload)
}
