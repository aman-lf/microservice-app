package main

import (
	"database/sql"
	"log"
)

// initDB creates necessary tables if they don't exist
func initDB(db *sql.DB) error {
	// Create menu_items table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS menu_items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		category VARCHAR(100) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Menu service database tables initialized")
	return nil
}
