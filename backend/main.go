package main

import (
	"net/http"
	"log"
)
func main(){

	mux := http.NewServerMux()


	log.Println("server running on port " + server.Addr)
	
	server := &http.Server(
		Handler: mux,
		Addr: "8080"
	)

	server.ListenAndServe(":8080", corsMiddleWare(mux))
}
