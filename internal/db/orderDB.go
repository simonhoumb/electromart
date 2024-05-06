package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"fmt"
	"log"
)

/*
AddOrder adds a single row to the ProductOrder table in the database. Returns the ID if successful, or an error if not.
*/
func AddOrder(order structs.ProductOrder) (string, error) {
	// Generate and retrieve new UUID
	id, err := GenerateUUID(Client)
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return "", err
	}

	// Insert product
	_, err2 := Client.Exec(
		`INSERT INTO ProductOrder (ID, UserID, TotalAmount, OrderDate, ShippedDate, EstimatedDelivery, Status, 
Comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		order.TotalAmount,
		order.OrderDate,
		order.ShippedDate,
		order.EstimatedDelivery,
		order.Status,
		order.Comments,
	)
	if err2 != nil {
		log.Println("Error inserting order: ", err2)
		return "", err2
	}

	return id, nil // Return the UUID
}

func GetAllOrdersByUserID(userID string) ([]structs.ProductOrder, error) {
	rows, err := Client.Query(`SELECT * FROM ProductOrder WHERE UserID = ?`, userID)
	if err != nil {
		log.Println("Error when selecting all orders: ", err)
		return nil, err
	}

	foundOrders, err2 := rowsToProductOrderSlice(rows)
	if err2 != nil {
		log.Println("Error when converting rows to slice: ", err2)
		return nil, err2
	}
	return foundOrders, nil
}

func GetOrderByID(orderID string) (*structs.ProductOrder, error) {
	exists, err := productOrderExists(orderID)
	if err != nil {
		log.Println("Error when checking if order exists: ", err)
		return nil, err
	}
	if exists {
		var order structs.ProductOrder
		err2 := Client.QueryRow(`SELECT * FROM ProductOrder WHERE ID = ?`, orderID).Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.OrderDate,
			&order.ShippedDate,
			&order.EstimatedDelivery,
			&order.Status,
			&order.Comments,
		)
		if err2 != nil {
			log.Println("Error when selecting order by ID: ", err2)
			return nil, err2
		}
		return &order, nil
	} else {
		log.Println("ProductOrder does not exist")
		return nil, fmt.Errorf("order ID does not match any orders in DB")
	}
}

func GetAllOrderItemsByOrderID(orderID string) ([]structs.OrderItem, error) {
	rows, err := Client.Query(`SELECT * FROM OrderItem WHERE OrderID = ?`, orderID)
	if err != nil {
		log.Println("Error when selecting all order items: ", err)
		return nil, err
	}
	foundOrderItems, err2 := rowsToOrderItemSlice(rows)
	if err2 != nil {
		log.Println("Error when converting rows to slice: ", err2)
		return nil, err2
	}
	return foundOrderItems, nil
}

/*
UpdateProductOrder updates a single row in the ProductOrder table in the database based on the ID in the provided
ProductOrder struct.
Returns nil if successful, or an error if not.
*/
func UpdateProductOrder(productOrder structs.ProductOrder) error {
	// Check if order exists
	exists, err := productOrderExists(productOrder.ID)
	if err != nil {
		log.Println("Error when checking if order exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec(
			`UPDATE ProductOrder SET Name = ?, BrandID = ?, CategoryID = ?, Description = ?, QtyInStock = ?, Price = ? WHERE ID = ?`,
		)
		if err2 != nil {
			log.Println("Error when updating order: ", err2)
			return err2
		}
	} else {
		log.Println("Order does not exist")
		return fmt.Errorf("order ID does not match any orders in DB")
	}
	return nil
}

/*
DeleteProductOrderByID deletes a single row from the ProductOrder table in the database based on the ID.
Returns nil if successful, or an error if not.
*/
func DeleteProductOrderByID(id string) error {
	exists, err := productOrderExists(id)
	if err != nil {
		log.Println("Error when checking if order exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec("DELETE FROM ProductOrder WHERE ID = ?", id)
		if err2 != nil {
			log.Println("Error when deleting order: ", err2)
			return err2
		}
	} else {
		log.Println("Order does not exist")
		return fmt.Errorf("order ID does not match any orders in DB")
	}
	return nil
}

/*
rowsToProductOrderSlice converts the rows from a SQL query to a slice of ProductOrder structs.
*/
func rowsToProductOrderSlice(rows *sql.Rows) ([]structs.ProductOrder, error) {
	var orderSlice []structs.ProductOrder
	for rows.Next() {
		var order structs.ProductOrder
		err2 := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.OrderDate,
			&order.ShippedDate,
			&order.EstimatedDelivery,
			&order.Status,
			&order.Comments,
		)
		if err2 != nil {
			return nil, err2
		}
		orderSlice = append(orderSlice, order)
	}
	return orderSlice, nil
}

/*
rowsToOrderItemSlice converts the rows from a SQL query to a slice of OrderItem struct.
*/
func rowsToOrderItemSlice(rows *sql.Rows) ([]structs.OrderItem, error) {
	var orderItemSlice []structs.OrderItem
	for rows.Next() {
		var orderItem structs.OrderItem
		err2 := rows.Scan(
			&orderItem.ID,
			&orderItem.OrderID,
			&orderItem.ProductID,
			&orderItem.Quantity,
		)
		if err2 != nil {
			return nil, err2
		}
		orderItemSlice = append(orderItemSlice, orderItem)
	}
	return orderItemSlice, nil
}

/*
productOrderExists checks if an order with the provided ID exists in the database.
*/
func productOrderExists(id string) (bool, error) {
	var exists bool
	err := Client.QueryRow(`SELECT EXISTS(SELECT * FROM "ProductOrder" WHERE ID = ?)`, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
