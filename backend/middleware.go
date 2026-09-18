package main

import (
	"context"
	"net/http"
    "os"
    "net"
    "os"
    "strings"
    "time"
)
type contextKey string
var userIDKey contextKey = "userID" // // make your own type like this to prevent collisions with other packages that might also use "userID" as a context key. 

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// Middleware to check if user is logged in
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        // Validate token, get user_id
        userID, _, err := validateToken(token)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        // Add userID to context so handlers can use it
        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        _, role, err := validateToken(token)
        if err != nil || role != "admin" {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    }
}

func rateLimit(limit int, window time.Duration) func(http.HandlerFunc) http.HandlerFunc {
    // store the limiter in closure
    limiter := &rateLimiter{
        visitors: make(map[string]*visitor),
    }
    
    go limiter.cleanup() // a method that loops forever, deleting visitor entries whose resetAt time has passed, so the map doesn't grow forever
    
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            ip := getIP(r)
            
            if !limiter.allow(ip, limit, window) {
                http.Error(w, "Too many requests", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        }
    }
}

func getIP(r *http.Request) string {
    // forwarded might be many ips seperated by commas. if someone contacted the website through a middleman the request doesnt come directly from the visitor, r.RemoteAddr would give the middlemans ip instead of the visitors. this is where X-forwarded-for comes in
    forwarded := r.Header.Get("X-Forwarded-For")
    if forwarded != "" {
        return strings.Split(forwarded, ",")[0]
    }
    
    // this time its direct but r.RemoteAddr is an ip + a port mashed together in one string, we only want the ip so we SplitHostPort. if an error occurs, sometimes no port at all like wierd/raw connections or if jut fails we return the whole thing as is.
    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return r.RemoteAddr
    }
    return host
}