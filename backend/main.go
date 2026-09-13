package main

import (
	"net/http"
	"log"
)
func main(){

	mux := http.NewServeMux()
	
	server := &http.Server{
		Handler: corsMiddleware(mux),
		Addr: ":8080",
	}

	mux.HandleFunc("GET /products", listProducts)
	mux.HandleFunc("GET /products/{id}", getProduct)
	mux.HandleFunc("POST /products", adminMiddleware(createProduct))
	mux.HandleFunc("PUT /products/{id}", adminMiddleware(updateProduct))
	mux.HandleFunc("DELETE /products/{id}", adminMiddleware(deleteProduct))

	mux.HandleFunc("POST /register", register)
	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("POST /refresh", refresh)
	mux.HandleFunc("POST /logout", logout)
	
	mux.HandleFunc("POST /cart", authMiddleware(addItem))
	mux.HandleFunc("GET /cart", authMiddleware(viewCart))
	mux.HandleFunc("PUT /cart", authMiddleware(updateQuantity))
	mux.HandleFunc("DELETE /cart/{id}", authMiddleware(removeItem))
	mux.HandleFunc("DELETE /cart", authMiddleware(clearCart))

	mux.HandleFunc("POST /orders", authMiddleware(createOrder))
	mux.HandleFunc("GET /orders", authMiddleware(listOrders))
	mux.HandleFunc("POST /create-payment-intent", authMiddleware(createPaymentIntent))
	mux.HandleFunc("POST /webhook", handleWebhook)

	log.Println("server running on port " + server.Addr)


	log.Fatal(server.ListenAndServe())
}
