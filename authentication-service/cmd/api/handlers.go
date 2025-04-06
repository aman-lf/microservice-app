package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (app *Config) Authenticate(c *gin.Context) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(c, &requestPayload)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Validate the user against the database
	user, err := app.getByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(c, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Check password
	valid, err := app.passwordMatches(user, requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(c, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Log authentication
	err = app.logRequest("authentication", fmt.Sprintf("User %s logged in", user.Email))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	payload := struct {
		Error   bool  `json:"error"`
		Message string `json:"message"`
		Data    any   `json:"data,omitempty"`
	}{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) CreateUser(c *gin.Context) {
	var user User

	err := app.readJSON(c, &user)
	if err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert the user
	stmt := `insert into users (email, first_name, last_name, password, active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	var newID int
	err = app.DB.QueryRow(stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}

	user.ID = newID
	user.Password = "" // Don't return the hashed password

	payload := struct {
		Error   bool  `json:"error"`
		Message string `json:"message"`
		Data    any   `json:"data,omitempty"`
	}{
		Error:   false,
		Message: fmt.Sprintf("Created user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(c, http.StatusCreated, payload)
}

func (app *Config) GetAllUsers(c *gin.Context) {
	var users []User
	
	query := `select id, email, first_name, last_name, active, created_at, updated_at from users order by last_name`
	
	rows, err := app.DB.Query(query)
	if err != nil {
		app.errorJSON(c, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			app.errorJSON(c, err, http.StatusInternalServerError)
			return
		}
		
		users = append(users, user)
	}
	
	payload := struct {
		Error   bool  `json:"error"`
		Message string `json:"message"`
		Data    any   `json:"data,omitempty"`
	}{
		Error:   false,
		Message: "Users retrieved",
		Data:    users,
	}
	
	app.writeJSON(c, http.StatusOK, payload)
}
