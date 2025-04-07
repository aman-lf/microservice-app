package main

import (
	"errors"
	"time"
)

// getAllInventoryItems retrieves all inventory items from the database
func (app *Config) getAllInventoryItems() ([]InventoryItem, error) {
	var items []InventoryItem

	query := `select id, item_name, quantity, unit, threshold, created_at, updated_at from inventory_items order by item_name`

	rows, err := app.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item InventoryItem
		err := rows.Scan(
			&item.ID,
			&item.ItemName,
			&item.Quantity,
			&item.Unit,
			&item.Threshold,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

// getInventoryItemByID retrieves an inventory item by its ID
func (app *Config) getInventoryItemByID(id int) (InventoryItem, error) {
	var item InventoryItem

	query := `select id, item_name, quantity, unit, threshold, created_at, updated_at from inventory_items where id = $1`

	row := app.DB.QueryRow(query, id)
	err := row.Scan(
		&item.ID,
		&item.ItemName,
		&item.Quantity,
		&item.Unit,
		&item.Threshold,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return InventoryItem{}, err
	}

	return item, nil
}

// insertInventoryItem adds a new inventory item to the database
func (app *Config) insertInventoryItem(item InventoryItem) (int, error) {
	// Set timestamps
	now := time.Now()

	var newID int
	stmt := `insert into inventory_items (item_name, quantity, unit, threshold, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id`

	err := app.DB.QueryRow(
		stmt,
		item.ItemName,
		item.Quantity,
		item.Unit,
		item.Threshold,
		now,
		now,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// updateInventoryItem updates an existing inventory item
func (app *Config) updateInventoryItem(item InventoryItem) error {
	// Check if the item exists
	_, err := app.getInventoryItemByID(item.ID)
	if err != nil {
		return errors.New("inventory item not found")
	}

	// Set updated timestamp
	now := time.Now()

	stmt := `update inventory_items set
		item_name = $1,
		quantity = $2,
		unit = $3,
		threshold = $4,
		updated_at = $5
		where id = $6`

	_, err = app.DB.Exec(
		stmt,
		item.ItemName,
		item.Quantity,
		item.Unit,
		item.Threshold,
		now,
		item.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// deleteInventoryItem removes an inventory item from the database
func (app *Config) deleteInventoryItem(id int) error {
	// Check if the item exists
	_, err := app.getInventoryItemByID(id)
	if err != nil {
		return errors.New("inventory item not found")
	}

	stmt := `delete from inventory_items where id = $1`

	_, err = app.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// adjustInventoryQuantity increases or decreases the quantity of an inventory item
func (app *Config) adjustInventoryQuantity(id int, adjustment int) error {
	// Check if the item exists
	item, err := app.getInventoryItemByID(id)
	if err != nil {
		return errors.New("inventory item not found")
	}

	// Calculate new quantity
	newQuantity := item.Quantity + adjustment

	// Don't allow negative inventory
	if newQuantity < 0 {
		return errors.New("adjustment would result in negative inventory")
	}

	// Set updated timestamp
	now := time.Now()

	stmt := `update inventory_items set
		quantity = $1,
		updated_at = $2
		where id = $3`

	_, err = app.DB.Exec(
		stmt,
		newQuantity,
		now,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

// getLowInventoryItems returns all inventory items that are at or below their threshold
func (app *Config) getLowInventoryItems() ([]InventoryItem, error) {
	var items []InventoryItem

	query := `select id, item_name, quantity, unit, threshold, created_at, updated_at
		from inventory_items
		where quantity <= threshold
		order by item_name`

	rows, err := app.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item InventoryItem
		err := rows.Scan(
			&item.ID,
			&item.ItemName,
			&item.Quantity,
			&item.Unit,
			&item.Threshold,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
