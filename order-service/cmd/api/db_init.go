package main

import (
	"database/sql"
	"log"
)

// Initialize the database tables if they don't exist
func initDB(db *sql.DB) error {
	// Create orders table
	orderTableQuery := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		customer_id INTEGER NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		total DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(orderTableQuery)
	if err != nil {
		return err
	}

	// Create order items table
	orderItemsTableQuery := `
	CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		menu_item_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 1,
		price DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(orderItemsTableQuery)
	if err != nil {
		return err
	}

	log.Println("Order service database tables initialized")
	return nil
}