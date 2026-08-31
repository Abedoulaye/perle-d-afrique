package main

import (
	"encoding/json"
	"context"
	"net/http"
	"io"
	"github.com/stripe/stripe-go/v81/webhook"
	"github.com/stripe/stripe-go/v81"
	"os"
)

func handleWebhook(w http.ResponseWriter, r* http.Request){
    const MaxBodyBytes = int64(65536) // 64 kb limit, huge request bodies can overwhelm the server, 64kb is plenty
    r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
    
    payload, err := io.ReadAll(r.Body) // JSON data Stripe sent
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }
    
    // Verify it's really from Stripe
    event, err := webhook.ConstructEvent(
        payload,
        r.Header.Get("Stripe-Signature"), // something that proves the request came from stripe
        os.Getenv("STRIPE_WEBHOOK_SECRET"),
    )
    if err != nil {
        http.Error(w, "Invalid signature", http.StatusBadRequest)
        return
    }

    switch event.Type {
    case "payment_intent.succeeded":
        var paymentIntent stripe.PaymentIntent
        if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil { // JSON.urmashall converts JSON bytes into Go struct, opposite of json.Marshal
            http.Error(w, "Error parsing webhook JSON", http.StatusBadRequest)
            return
        }

        orderID := paymentIntent.Metadata["order_id"]
        
        // Update order status to paid
        _, err := db.Exec(r.Context(), "UPDATE orders SET status = 'paid' WHERE id = $1", orderID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Clear the cart for this user
        var userID int
        err = db.QueryRow(r.Context(), "SELECT user_id FROM orders WHERE id = $1", orderID).Scan(&userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Delete cart items
        _, err = db.Exec(r.Context(), "DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id = $1)", userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

    case "payment_intent.payment_failed":
        var paymentIntent stripe.PaymentIntent
        if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
            http.Error(w, "Error parsing webhook JSON", http.StatusBadRequest)
            return
        }

        orderID := paymentIntent.Metadata["order_id"]
        
        // Update order status to failed
        _, err := db.Exec(r.Context(), "UPDATE orders SET status = 'failed' WHERE id = $1", orderID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

    default:
        // Ignore other event types
    }
    
    // Always return 200 OK to Stripe
    w.WriteHeader(http.StatusOK)
}
