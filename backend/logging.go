package main

import (
	"log"
	"net/http"
)

func serverError(w http.ResponseWriter, err error, context string){
	log.Printf("[ERROR] %s: %v", context, err)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func clientError(w http.ResponseWriter, status int, message string){
	log.Printf("[WARN] %d %s", status, message)
	http.Error(w, message, status)
}