package main

import (
	"net/http"
	"log"
	"time"
)
func main(){

	mux := http.NewServeMux()
	
	server := &http.Server{
		Handler: corsMiddleware(mux),
		Addr: ":8080",
	}

	generalLimit := rateLimit(100, time.Minute)

	mux.HandleFunc("GET /products", generalLimit(listProducts))
	mux.HandleFunc("GET /products/{id}", generalLimit(getProduct))
	mux.HandleFunc("POST /products", generalLimit(adminMiddleware(createProduct)))
	mux.HandleFunc("PUT /products/{id}", generalLimit(adminMiddleware(updateProduct)))
	mux.HandleFunc("DELETE /products/{id}", generalLimit(adminMiddleware(deleteProduct)))

	loginLimit := rateLimit(5, 15*time.Minute)
	registerLimit := rateLimit(3, time.Hour)

	mux.HandleFunc("POST /register", registerLimit(register))
	mux.HandleFunc("POST /login", loginLimit(login))
	mux.HandleFunc("POST /refresh", generalLimit(refresh))
	mux.HandleFunc("POST /logout", generalLimit(logout))
	
	mux.HandleFunc("POST /cart", generalLimit(authMiddleware(addItem)))
	mux.HandleFunc("GET /cart", generalLimit(authMiddleware(viewCart)))
	mux.HandleFunc("PUT /cart", generalLimit(authMiddleware(updateQuantity)))
	mux.HandleFunc("DELETE /cart/{id}", generalLimit(authMiddleware(removeItem)))
	mux.HandleFunc("DELETE /cart", generalLimit(authMiddleware(clearCart)))

	mux.HandleFunc("POST /orders", generalLimit(authMiddleware(createOrder)))
	mux.HandleFunc("GET /orders", generalLimit(authMiddleware(listOrders)))
	mux.HandleFunc("POST /create-payment-intent", generalLimit(authMiddleware(createPaymentIntent)))

	mux.HandleFunc("POST /webhook", handleWebhook)

	log.Println("server running on port " + server.Addr)


	log.Fatal(server.ListenAndServe())
}
