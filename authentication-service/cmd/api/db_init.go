package main

import (
	"database/sql"
	"log"
)

// initDB creates necessary tables if they don't exist
func initDB(db *sql.DB) error {
	// Create users table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		active INT NOT NULL DEFAULT 1,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL
	);
	`
	
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	
	log.Println("Database tables initialized")
	return nil
}