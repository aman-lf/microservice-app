package main

import (
	"database/sql"
	"log"
)

// Initialize the database tables if they don't exist
func initDB(db *sql.DB) error {
	// Create inventory items table
	query := `
	CREATE TABLE IF NOT EXISTS inventory_items (
		id SERIAL PRIMARY KEY,
		item_name VARCHAR(255) NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 0,
		unit VARCHAR(50) NOT NULL,
		threshold INTEGER NOT NULL DEFAULT 5,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Inventory service database tables initialized")
	return nil
}
