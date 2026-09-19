package main

import (
	"encoding/json"
	"net/http"
	"io"
	"github.com/stripe/stripe-go/v85/webhook"
	"github.com/stripe/stripe-go/v85"
	"os"
    "fmt"
)

func handleWebhook(w http.ResponseWriter, r* http.Request){
    const MaxBodyBytes = int64(65536) // 64 kb limit, huge request bodies can overwhelm the server, 64kb is plenty
    r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
    
    payload, err := io.ReadAll(r.Body) // JSON data Stripe sent
    if err != nil {
        clientError(w, http.StatusBadRequest, "error reading request body")
        return
    }
    
    // Verify it's really from Stripe
    event, err := webhook.ConstructEvent(
        payload,
        r.Header.Get("Stripe-Signature"), // something that proves the request came from stripe
        os.Getenv("STRIPE_WEBHOOK_SECRET"),
    )
    if err != nil {
        clientError(w, http.StatusBadRequest, "invalid signature")
        return
    }

    switch event.Type {
    case "payment_intent.succeeded":
        var paymentIntent stripe.PaymentIntent
        if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil { // JSON.urmashall converts JSON bytes into Go struct, opposite of json.Marshal
            clientError(w, http.StatusBadRequest, "error parsing webhook JSON")
            return
        }
        orderID := paymentIntent.Metadata["order_id"]
        if orderID == "" {
            fmt.Println("Webhook received with missing order_id, skipping")
            w.WriteHeader(http.StatusOK)
            return
        }
        // Update order status to paid
        _, err := db.Exec(r.Context(), "UPDATE orders SET status = 'paid' WHERE id = $1", orderID)
        if err != nil {
            serverError(w, err, "webhook update order")
            return
        }

        // Clear the cart for this user
        var userID int
        err = db.QueryRow(r.Context(), "SELECT user_id FROM orders WHERE id = $1", orderID).Scan(&userID)
        if err != nil {
            serverError(w, err, "webhook select user_id")
            return
        }

        // Delete cart items
        _, err = db.Exec(r.Context(), "DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id = $1)", userID)
        if err != nil {
            serverError(w, err, "webhook clear cart")
            return
        }

    case "payment_intent.payment_failed":
        var paymentIntent stripe.PaymentIntent
        if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
            clientError(w, http.StatusBadRequest, "error parsing webhook JSON")
            return
        }
        orderID := paymentIntent.Metadata["order_id"]
        if orderID == "" {
            fmt.Println("Webhook received with missing order_id, skipping")
            w.WriteHeader(http.StatusOK)
            return
        }
        // Update order status to failed
        _, err := db.Exec(r.Context(), "UPDATE orders SET status = 'failed' WHERE id = $1", orderID)
        if err != nil {
            serverError(w, err, "webhook update order")
            return
        }

    default:
        // Ignore other event types
    }
    
    // Always return 200 OK to Stripe
    w.WriteHeader(http.StatusOK)
    
}
