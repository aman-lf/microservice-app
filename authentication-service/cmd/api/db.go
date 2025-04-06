package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// openDB creates a new database connection
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// getByEmail returns a user by email
func (app *Config) getByEmail(email string) (*User, error) {
	var user User
	query := `select id, email, first_name, last_name, password, active, created_at, updated_at from users where email = $1`

	row := app.DB.QueryRow(query, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// passwordMatches checks if provided password matches stored hash
func (app *Config) passwordMatches(user *User, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// InsertUser adds a new user to the database
func (app *Config) InsertUser(user User) (int, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

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
		return 0, err
	}

	return newID, nil
}

// logRequest logs a request to the logger service
func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	// In a real application, this would send a request to the logger service
	log.Printf("Log: %s - %s", entry.Name, entry.Data)
	return nil
}
