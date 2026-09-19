package main

import (
	"net/http"
	"context"
	"encoding/json"
	"time"
)

func createOrder(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userID := r.Context().Value(userIDKey).(int)
	cartID, err := getOrCreateCart(ctx, userID)
	if err != nil {
		serverError(w, err, "createOrder getOrCreateCart")
		return
	}


	rows, err := db.Query(ctx, "SELECT ci.product_id, ci.quantity, p.price_cents, p.stock FROM cart_items ci JOIN products p ON ci.product_id = p.id WHERE ci.cart_id = $1", cartID)
	if err != nil {
		serverError(w, err, "createOrder query")
		return
	}
	defer rows.Close()

	var items []OrderItem
	var totalPrice int64
	for rows.Next(){
		var item OrderItem
		var stock int
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.PriceInCentsAtPurchase, &stock); err != nil {
			serverError(w, err, "createOrder scan")
			return
		}

		if stock < item.Quantity {
			clientError(w, http.StatusBadRequest, "not enough stock for this item")
			return
		}

		totalPrice += item.PriceInCentsAtPurchase * int64(item.Quantity)
		items = append(items, item)
	}
	if len(items) == 0 {
		clientError(w, http.StatusBadRequest, "cart is empty")
    	return
	}
	
	if err := rows.Err(); err != nil {
		serverError(w, err, "createOrder rows error")
		return
	}
// the upcoming syntax is useful for grouping database operations into one, this way if they all succeed it results in a success but if a single fails the entire order has failed.

	tx, err := db.Begin(ctx)
	if err != nil {
		serverError(w, err, "createOrder begin tx")
		return
	}
	defer tx.Rollback(ctx) // If we return early, rollback everything

	const taxRate = int64(5)
	taxCents := totalPrice * taxRate / 100
	finalTotal := totalPrice + taxCents
	var orderID int
	err = tx.QueryRow(ctx, "INSERT INTO orders (user_id, status, total_cents) VALUES ($1, $2, $3) RETURNING id", userID, "pending", finalTotal).Scan(&orderID)
	if err != nil {
		serverError(w, err, "createOrder insert order")
		return 
	}

	for _, item := range items {
		var currentStock int
		err = tx.QueryRow(ctx, "SELECT stock FROM products WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&currentStock)
		if err != nil {
			serverError(w, err, "createOrder select stock")
			return
		}
		if currentStock < item.Quantity {
			clientError(w, http.StatusBadRequest, "not enough stock for item")
			return
		}

		_, err = tx.Exec(ctx, "INSERT INTO order_items (order_id, product_id, quantity, price_cents_at_purchase) VALUES ($1, $2, $3, $4)", orderID, item.ProductID, item.Quantity, item.PriceInCentsAtPurchase)
		if err != nil {
			serverError(w, err, "createOrder insert order_item")
			return
		}
		_, err = tx.Exec(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			serverError(w, err, "createOrder update stock")
			return
		}
	}
// Everything succeeded, commit the transaction
	if err := tx.Commit(ctx); err != nil {
		serverError(w, err, "createOrder commit")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": orderID,
		"total_cents": finalTotal,
		"status": "pending",
		"items": items,
	})

}

func listOrders(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	userID := r.Context().Value(userIDKey).(int)

	rows, err := db.Query(ctx, "SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		serverError(w, err, "listOrders query")
		return
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next(){
		var order Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalInCents, &order.CreatedAt); err != nil {
			serverError(w, err, "listOrders scan")
			return
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		serverError(w, err, "listOrders rows error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}