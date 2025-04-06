package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *Config) GetAllInventoryItems(c *gin.Context) {
	items, err := app.getAllInventoryItems()
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool            `json:"error"`
		Message string          `json:"message"`
		Data    []InventoryItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Inventory items retrieved",
		Data:    items,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) GetInventoryItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	item, err := app.getInventoryItemByID(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Data    InventoryItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Inventory item retrieved",
		Data:    item,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) CreateInventoryItem(c *gin.Context) {
	var item InventoryItem

	err := app.readJSON(c, &item)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Validate the inventory item
	if item.ItemName == "" {
		app.errorJSON(c, errors.New("invalid inventory item data"), http.StatusBadRequest)
		return
	}

	// Create the inventory item
	newID, err := app.insertInventoryItem(item)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	// Get the newly created item
	newItem, err := app.getInventoryItemByID(newID)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Data    InventoryItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Inventory item created",
		Data:    newItem,
	}

	app.writeJSON(c, http.StatusCreated, payload)
}

func (app *Config) UpdateInventoryItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	var item InventoryItem
	err = app.readJSON(c, &item)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Make sure ID matches
	item.ID = id

	// Update the inventory item
	err = app.updateInventoryItem(item)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	// Get the updated item
	updatedItem, err := app.getInventoryItemByID(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Data    InventoryItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Inventory item updated",
		Data:    updatedItem,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) DeleteInventoryItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	err = app.deleteInventoryItem(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Inventory item deleted",
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) AdjustInventory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	var adjustmentPayload struct {
		Quantity int `json:"quantity"`
	}

	err = app.readJSON(c, &adjustmentPayload)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Adjust the inventory
	err = app.adjustInventoryQuantity(id, adjustmentPayload.Quantity)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	// Get the updated item
	updatedItem, err := app.getInventoryItemByID(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
		Data    InventoryItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Inventory adjusted",
		Data:    updatedItem,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) CheckLowInventory(c *gin.Context) {
	items, err := app.getLowInventoryItems()
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool            `json:"error"`
		Message string          `json:"message"`
		Data    []InventoryItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Low inventory items retrieved",
		Data:    items,
	}

	app.writeJSON(c, http.StatusOK, payload)
}
