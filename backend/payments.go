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
        clientError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    var totalCents int64
    var status string
    err := db.QueryRow(ctx, "SELECT total_cents, status FROM orders WHERE id = $1 AND user_id = $2", req.OrderID, userID).Scan(&totalCents, &status)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            clientError(w, http.StatusNotFound, "order not found")
            return
        }
        serverError(w, err, "createPaymentIntent query order")
        return
    }

    if status != "pending" {
        clientError(w, http.StatusBadRequest, "order is not payable")
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
        serverError(w, err, "createPaymentIntent stripe create")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{
        "client_secret": pi.ClientSecret,
    })
}

