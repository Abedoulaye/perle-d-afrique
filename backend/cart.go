package main

import (
	"net/http"
	"context"
	"encoding/json"
	"time"
	"github.com/jackc/pgx/v5"
	"strconv"
	"errors"
)


func getOrCreateCart(ctx context.Context, userID int) (int, error) {
    var cartID int
    err := db.QueryRow(ctx, "SELECT id FROM carts WHERE user_id = $1", userID).Scan(&cartID)
    if errors.Is(err, pgx.ErrNoRows) {
        err = db.QueryRow(ctx, "INSERT INTO carts (user_id) VALUES ($1) RETURNING id", userID).Scan(&cartID)
        if err != nil {
            return 0, err
        }
    } else if err != nil {
		return 0, err
	}
    return cartID, nil
}

func addItem(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userID := r.Context().Value(userIDKey).(int)
	
	var input struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	if input.Quantity <= 0 {
	clientError(w, http.StatusBadRequest, "quantity must be positive")
    return
}

	cartID, err := getOrCreateCart(ctx, userID)
	if err != nil {
		serverError(w, err, "addItem getOrCreateCart")
		return
	}

	var existingProductID int
	err = db.QueryRow(ctx, "SELECT product_id FROM cart_items WHERE cart_id = $1 AND product_id = $2", cartID, input.ProductID).Scan(&existingProductID)
	if errors.Is(err, pgx.ErrNoRows) {
		_, err = db.Exec(ctx, "INSERT INTO cart_items (cart_id, product_id, quantity) VALUES($1, $2, $3)", cartID, input.ProductID, input.Quantity)
		if err != nil {
			serverError(w, err, "addItem insert")
			return
		}
	} else if err != nil {
		serverError(w, err, "addItem select")
		return
	} else{
		_, err = db.Exec(ctx, "UPDATE cart_items SET quantity = quantity + $1 WHERE cart_id = $2 AND product_id = $3", input.Quantity, cartID, input.ProductID)
		if err != nil {
			serverError(w, err, "addItem update")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "item added to cart"})
}

func viewCart(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userID := r.Context().Value(userIDKey).(int)
	rows, err := db.Query(ctx, "SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, p.name, p.price_cents, p.description, p.stock, p.image FROM cart_items ci JOIN carts c ON c.id = ci.cart_id JOIN products p ON p.id = ci.product_id WHERE c.user_id = $1", userID)

	if err != nil {
		serverError(w, err, "viewCart query")
		return
	}
	defer rows.Close()


	var c CartItem
	var p Product

	items := []CartDetail{}
	for rows.Next(){

		if err := rows.Scan(&c.ID, &c.CartID, &c.ProductID, &c.Quantity, &p.Name, &p.PriceInCents, &p.Description, &p.Stock, &p.Image); err != nil {
			serverError(w, err, "viewCart scan")
			return
		}
		
		detail := CartDetail{
			Item: c,
			Prod: p,
		}

		items = append(items, detail)

	}

	if err := rows.Err(); err != nil {
		serverError(w, err, "viewCart rows error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)

}

func updateQuantity(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userID := r.Context().Value(userIDKey).(int)

	cartID, err := getOrCreateCart(ctx, userID)
	if err != nil {
		serverError(w, err, "updateQuantity getOrCreateCart")
		return
	}

	var updated struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		clientError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	if updated.Quantity <= 0 {
    clientError(w, http.StatusBadRequest, "quantity must be positive")
    return
	}

	result, err := db.Exec(ctx, "UPDATE cart_items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3", updated.Quantity, cartID, updated.ProductID)
	if err != nil {
		serverError(w, err, "updateQuantity update")
		return
	}
	if result.RowsAffected() == 0{
		clientError(w, http.StatusNotFound, "item not found in cart")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func removeItem(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userID := r.Context().Value(userIDKey).(int)
	productID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		clientError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	cartID, err := getOrCreateCart(ctx, userID)
	if err != nil {
		serverError(w, err, "removeItem getOrCreateCart")
		return
	}

	result, err := db.Exec(ctx, "DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2", cartID, productID)
	if err != nil {
		serverError(w, err, "removeItem delete query")
		return
	}

	if result.RowsAffected() == 0{
		clientError(w, http.StatusNotFound, "item not found in cart")
		return
	}

	w.WriteHeader(http.StatusOK)
}


func clearCart(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userID := r.Context().Value(userIDKey).(int)

	cartID, err := getOrCreateCart(ctx, userID)
	if err != nil {
		serverError(w, err, "clearCart getOrCreateCart")
		return
	}
	result, err := db.Exec(ctx, "DELETE FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		serverError(w, err, "clearCart delete")
		return
	}

	if result.RowsAffected() == 0{
		clientError(w, http.StatusNotFound, "cart is empty")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}