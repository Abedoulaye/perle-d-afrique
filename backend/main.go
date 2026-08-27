package main

import (
	"net/http"
	"log"
)
func main(){

	mux := http.NewServerMux()

	mux.HandleFunc("GET /products", listProducts)
	mux.HandleFunc("GET /products/{id}", getProduct)
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("PUT /products/{id}", updateProduct)
	mux.HandleFunc("DELETE /products/{id}", deleteProduct)

	log.Println("server running on port " + server.Addr)
	
	server := &http.Server(
		Handler: mux,
		Addr: "8080"
	)

	server.ListenAndServe(":8080", corsMiddleWare(mux))
}
