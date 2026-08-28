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

	log.Println("server running on port " + server.Addr)


	log.Fatal(server.ListenAndServe())
}
