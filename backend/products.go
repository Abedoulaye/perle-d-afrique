package main

import (
	"net/http"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"context"
	"time"
	"errors"
	"strings"
)

func listProducts(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.Query(ctx, "SELECT * FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products = []Product{}
	for rows.Next(){
		var prod Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.PriceInCents, &prod.Stock, &prod.Image); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, prod)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")

	var prod Product

	err := db.QueryRow(ctx, "SELECT id, name, description, price_cents, stock, image FROM products WHERE id = $1", id).Scan(&prod.ID, &prod.Name, &prod.Description, &prod.PriceInCents, &prod.Stock, &prod.Image)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prod)
}

func createProduct(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var prod Product
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prod.Name = strings.TrimSpace(prod.Name)
	prod.Description = strings.TrimSpace(prod.Description)
	prod.Image = strings.TrimSpace(prod.Image)

	if errs := validateProduct(prod); len(errs) > 0 {
		http.Error(w, strings.Join(errs, ","), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(ctx, "INSERT INTO products (name, description, price_cents, stock, image) VALUES($1, $2, $3, $4, $5) RETURNING id", prod.Name, prod.Description, prod.PriceInCents, prod.Stock, prod.Image).Scan(&prod.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prod)
}


func updateProduct(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	id := r.PathValue("id")

	var updated Product
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated.Name = strings.TrimSpace(updated.Name)
	updated.Description = strings.TrimSpace(updated.Description)
	updated.Image = strings.TrimSpace(updated.Image)

	if errs := validateProduct(updated); len(errs) > 0 {
		http.Error(w, strings.Join(errs, ","), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(ctx, "UPDATE products SET name = $1, description = $2, price_cents = $3, stock = $4, image = $5 WHERE id = $6 RETURNING id", updated.Name, updated.Description, updated.PriceInCents, updated.Stock, updated.Image, id).Scan(&updated.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)

}

func deleteProduct(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")

	result, err := db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0{
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}