package main

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "time"
    "github.com/jackc/pgx/v5"
    "github.com/stripe/stripe-go/v85"
    "github.com/stripe/stripe-go/v85/paymentintent"
    "errors"
    "strconv"
)

func init() {
    stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func createPaymentIntent(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

	userID := r.Context().Value(userIDKey).(int)

    var req struct {
        OrderID int64 `json:"order_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var totalCents int64
    err := db.QueryRow(ctx, "SELECT total_cents FROM orders WHERE id = $1 AND user_id = $2", req.OrderID, userID).Scan(&totalCents) 
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            http.Error(w, "Order not found", http.StatusNotFound)
            return
        } 
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    

    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(totalCents),
        Currency: stripe.String("usd"),
        Metadata: map[string]string{
            "order_id": strconv.FormatInt(req.OrderID, 10),
        },
    }

    pi, err := paymentintent.New(params)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "client_secret": pi.ClientSecret,
    })
}

