package main

import (
        "errors"
        "fmt"
        "time"
)

// getAllOrders retrieves all orders from the database
func (app *Config) getAllOrders() ([]Order, error) {
        var orders []Order

        // Get all orders
        query := `select id, customer_id, status, total, created_at, updated_at from orders order by created_at desc`

        rows, err := app.DB.Query(query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        for rows.Next() {
                var order Order
                err := rows.Scan(
                        &order.ID,
                        &order.CustomerID,
                        &order.Status,
                        &order.Total,
                        &order.CreatedAt,
                        &order.UpdatedAt,
                )
                if err != nil {
                        return nil, err
                }

                // Get order items for this order
                orderItems, err := app.getOrderItems(order.ID)
                if err != nil {
                        return nil, err
                }

                order.Items = orderItems
                orders = append(orders, order)
        }

        return orders, nil
}

// getOrderByID retrieves an order by its ID with all associated items
func (app *Config) getOrderByID(id int) (Order, error) {
        var order Order

        // Get the order
        query := `select id, customer_id, status, total, created_at, updated_at from orders where id = $1`

        row := app.DB.QueryRow(query, id)
        err := row.Scan(
                &order.ID,
                &order.CustomerID,
                &order.Status,
                &order.Total,
                &order.CreatedAt,
                &order.UpdatedAt,
        )

        if err != nil {
                return Order{}, err
        }

        // Get order items for this order
        orderItems, err := app.getOrderItems(order.ID)
        if err != nil {
                return Order{}, err
        }

        order.Items = orderItems

        return order, nil
}

// getOrdersByCustomer retrieves all orders for a customer
func (app *Config) getOrdersByCustomer(customerID int) ([]Order, error) {
        var orders []Order

        // Get orders for this customer
        query := `select id, customer_id, status, total, created_at, updated_at 
                from orders 
                where customer_id = $1
                order by created_at desc`

        rows, err := app.DB.Query(query, customerID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        for rows.Next() {
                var order Order
                err := rows.Scan(
                        &order.ID,
                        &order.CustomerID,
                        &order.Status,
                        &order.Total,
                        &order.CreatedAt,
                        &order.UpdatedAt,
                )
                if err != nil {
                        return nil, err
                }

                // Get order items for this order
                orderItems, err := app.getOrderItems(order.ID)
                if err != nil {
                        return nil, err
                }

                order.Items = orderItems
                orders = append(orders, order)
        }

        return orders, nil
}

// getOrderItems gets all items for a specific order
func (app *Config) getOrderItems(orderID int) ([]OrderItem, error) {
        var items []OrderItem

        query := `select id, order_id, menu_item_id, quantity, price
                from order_items
                where order_id = $1`

        rows, err := app.DB.Query(query, orderID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        for rows.Next() {
                var item OrderItem
                err := rows.Scan(
                        &item.ID,
                        &item.OrderID,
                        &item.MenuItemID,
                        &item.Quantity,
                        &item.Price,
                )
                if err != nil {
                        return nil, err
                }

                items = append(items, item)
        }

        return items, nil
}

// insertOrder creates a new order with all associated items
func (app *Config) insertOrder(order Order) (int, error) {
        // Begin transaction
        tx, err := app.DB.Begin()
        if err != nil {
                return 0, err
        }
        defer func() {
                if err != nil {
                        tx.Rollback()
                }
        }()

        // Set timestamps
        now := time.Now()

        // Calculate order total from items
        var total float64
        for _, item := range order.Items {
                total += item.Price * float64(item.Quantity)
        }

        // Insert the order
        var newOrderID int
        stmt := `insert into orders (customer_id, status, total, created_at, updated_at)
                values ($1, $2, $3, $4, $5) returning id`

        err = tx.QueryRow(
                stmt,
                order.CustomerID,
                "pending", // Default status
                total,
                now,
                now,
        ).Scan(&newOrderID)

        if err != nil {
                return 0, err
        }

        // Insert order items
        for _, item := range order.Items {
                stmt := `insert into order_items (order_id, menu_item_id, quantity, price)
                        values ($1, $2, $3, $4)`

                _, err := tx.Exec(
                        stmt,
                        newOrderID,
                        item.MenuItemID,
                        item.Quantity,
                        item.Price,
                )

                if err != nil {
                        return 0, err
                }
        }

        // Commit transaction
        err = tx.Commit()
        if err != nil {
                return 0, err
        }

        // Log the new order
        app.logNewOrder(newOrderID, order.CustomerID, total)

        return newOrderID, nil
}

// updateOrderStatus updates the status of an order
func (app *Config) updateOrderStatus(orderID int, status string) error {
        // Check if the order exists
        _, err := app.getOrderByID(orderID)
        if err != nil {
                return errors.New("order not found")
        }

        // Validate status
        validStatuses := map[string]bool{
                "pending":   true,
                "preparing": true,
                "ready":     true,
                "completed": true,
                "cancelled": true,
        }

        if !validStatuses[status] {
                return fmt.Errorf("invalid status: %s", status)
        }

        // Set updated timestamp
        now := time.Now()

        stmt := `update orders set 
                status = $1,
                updated_at = $2
                where id = $3`

        _, err = app.DB.Exec(
                stmt,
                status,
                now,
                orderID,
        )

        if err != nil {
                return err
        }

        // Log the status change
        app.logOrderStatusChange(orderID, status)

        return nil
}

// logNewOrder logs when a new order is created
func (app *Config) logNewOrder(orderID, customerID int, total float64) {
        // In a real application, this would send a request to the logger service
}

// logOrderStatusChange logs when an order status changes
func (app *Config) logOrderStatusChange(orderID int, status string) {
        // In a real application, this would send a request to the logger service
}
