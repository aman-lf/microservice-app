package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *Config) GetAllMenuItems(c *gin.Context) {
	items, err := app.getAllMenuItems()
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool       `json:"error"`
		Message string     `json:"message"`
		Data    []MenuItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Menu items retrieved",
		Data:    items,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) GetMenuItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	item, err := app.getMenuItemByID(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool     `json:"error"`
		Message string   `json:"message"`
		Data    MenuItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Menu item retrieved",
		Data:    item,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) CreateMenuItem(c *gin.Context) {
	var item MenuItem

	err := app.readJSON(c, &item)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Validate the menu item
	if item.Name == "" || item.Price <= 0 {
		app.errorJSON(c, errors.New("invalid menu item data"), http.StatusBadRequest)
		return
	}

	// Create the menu item
	newID, err := app.insertMenuItem(item)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	// Get the newly created item
	newItem, err := app.getMenuItemByID(newID)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool     `json:"error"`
		Message string   `json:"message"`
		Data    MenuItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Menu item created",
		Data:    newItem,
	}

	app.writeJSON(c, http.StatusCreated, payload)
}

func (app *Config) UpdateMenuItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	var item MenuItem
	err = app.readJSON(c, &item)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Make sure ID matches
	item.ID = id

	// Update the menu item
	err = app.updateMenuItem(item)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	// Get the updated item
	updatedItem, err := app.getMenuItemByID(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool     `json:"error"`
		Message string   `json:"message"`
		Data    MenuItem `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Menu item updated",
		Data:    updatedItem,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) DeleteMenuItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.errorJSON(c, errors.New("invalid id parameter"), http.StatusBadRequest)
		return
	}

	err = app.deleteMenuItem(id)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Menu item deleted",
	}

	app.writeJSON(c, http.StatusOK, payload)
}
