package main

import (
	"log"
	"net/http"
)

func serverError(w http.ResponseWriter, err error, context string){
	log.Printf("[Error] %s: %v", context, err)
	http.Error(w, "internal server error", http.StatusInternalServeroError)
}

func clientError(w http.ResponseWriter, status int, message string){
	log.Printtf("[WARN] %d %s", status, message)
	http.Error(w, message, status)
}