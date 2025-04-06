package main

import (
	"errors"
	"time"
)

// getAllMenuItems retrieves all menu items from the database
func (app *Config) getAllMenuItems() ([]MenuItem, error) {
	var items []MenuItem

	query := `select id, name, description, price, category, created_at, updated_at from menu_items order by name`

	rows, err := app.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item MenuItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.Category,
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

// getMenuItemByID retrieves a menu item by its ID
func (app *Config) getMenuItemByID(id int) (MenuItem, error) {
	var item MenuItem

	query := `select id, name, description, price, category, created_at, updated_at from menu_items where id = $1`

	row := app.DB.QueryRow(query, id)
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.Category,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return MenuItem{}, err
	}

	return item, nil
}

// insertMenuItem adds a new menu item to the database
func (app *Config) insertMenuItem(item MenuItem) (int, error) {
	// Set timestamp
	now := time.Now().Format(time.RFC3339)

	var newID int
	stmt := `insert into menu_items (name, description, price, category, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id`

	err := app.DB.QueryRow(
		stmt,
		item.Name,
		item.Description,
		item.Price,
		item.Category,
		now,
		now,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// updateMenuItem updates an existing menu item
func (app *Config) updateMenuItem(item MenuItem) error {
	// Check if the item exists
	_, err := app.getMenuItemByID(item.ID)
	if err != nil {
		return errors.New("menu item not found")
	}

	// Set updated timestamp
	now := time.Now().Format(time.RFC3339)

	stmt := `update menu_items set 
		name = $1, 
		description = $2, 
		price = $3, 
		category = $4,
		updated_at = $5
		where id = $6`

	_, err = app.DB.Exec(
		stmt,
		item.Name,
		item.Description,
		item.Price,
		item.Category,
		now,
		item.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// deleteMenuItem removes a menu item from the database
func (app *Config) deleteMenuItem(id int) error {
	// Check if the item exists
	_, err := app.getMenuItemByID(id)
	if err != nil {
		return errors.New("menu item not found")
	}

	stmt := `delete from menu_items where id = $1`

	_, err = app.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
