package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
	"encoding/base64"
	"crypto/sha256"
	"encoding/hex"
    "strings"
    "log"
)

func register(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        clientError(w, http.StatusBadRequest, "invalid request body")
        return
	}
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

    if err := validateEmail(u.Email); err != nil {
        clientError(w, http.StatusBadRequest, err.Error())
        return
    }
    if err := validatePassword(u.Password); err != nil {
        clientError(w, http.StatusBadRequest, err.Error())
        return
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
        serverError(w, err, "register bcrypt")
        return
	}

    token, err := generateRefreshToken()
    if err != nil {
        serverError(w, err, "register token gen")
        return
    }
    expires := time.Now().Add(24 * time.Hour)

	err = db.QueryRow(ctx, "INSERT INTO users (email, password_hash, verification_token, verification_expires) VALUES ($1, $2, $3, $4) RETURNING id, role, email_verified", u.Email, string(hashedPassword), hashToken(token), expires).Scan(&u.ID, &u.Role, &u.EmailVerified)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505"{
            clientError(w, http.StatusConflict, "email already exists")
            return
		}
        serverError(w, err, "register insert")
        return
	}

    if err := sendVerificationEmail(u.Email, token); err != nil {
        log.Printf("failed to send verification email to %s: %v", u.Email, err)
    }

	u.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func login(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        clientError(w, http.StatusBadRequest, "invalid request body")
        return
	}
    
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))

    if u.Email == "" || u.Password == "" {
        clientError(w, http.StatusUnauthorized, "invalid email and password")
        return
    }


	var storedHash string
	err := db.QueryRow(ctx, "SELECT id, email, password_hash, role, email_verified FROM users WHERE email = $1", u.Email).Scan(&u.ID, &u.Email, &storedHash, &u.Role, &u.EmailVerified)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
            clientError(w, http.StatusUnauthorized, "invalid email or password")
            return
		}
        serverError(w, err, "login select query")
        return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(u.Password)); err != nil {
        clientError(w, http.StatusUnauthorized, "invalid email or password")
        return
	}

    accessToken, err := generateAccessToken(u.ID, u.Role)
    if err != nil {
        serverError(w, err, "generate access token")
        return
    }

    refreshToken, err := generateRefreshToken()
    if err != nil {
        serverError(w, err, "generate refresh token")
        return
    }

    // Store hashed refresh token in DB (30 days)
    _, err = db.Exec(ctx,
        "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)",
        u.ID, hashToken(refreshToken), time.Now().Add(30*24*time.Hour))
    if err != nil {
        serverError(w, err, "insert into refresh token table")
        return
    }

    u.Password = ""
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user":          u,
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}

func validateToken(tokenString string) (int, string, error) {

    // Remove "Bearer " prefix if present
    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        tokenString = tokenString[7:]
    }
    
    // Parse and validate token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    
    if err != nil {
        return 0, "", err
    }
    
    // Extract claims
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := int(claims["user_id"].(float64))
		role := claims["role"].(string)
        return userID, role, nil
    }
    
    return 0, "", fmt.Errorf("invalid token")
}

func generateAccessToken(userID int, role string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "role":    role,
        "exp":     time.Now().Add(15 * time.Minute).Unix(),
        "iat":     time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func generateRefreshToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}

func refresh(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        clientError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    // Look up the hashed token
    var userID int
    var role string
    err := db.QueryRow(ctx, `
        SELECT u.id, u.role 
        FROM refresh_tokens rt
        JOIN users u ON u.id = rt.user_id
        WHERE rt.token_hash = $1 AND rt.expires_at > NOW()
    `, hashToken(req.RefreshToken)).Scan(&userID, &role)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            clientError(w, http.StatusUnauthorized, "invalid or expired refresh token")
            return
        }
            serverError(w, err, "refresh select query")
        return
    }

    // Issue new access token
    accessToken, err := generateAccessToken(userID, role)
    if err != nil {
        serverError(w, err, "refresh generate access token")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": accessToken,
    })
}

func logout(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        clientError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    _, err := db.Exec(ctx, "DELETE FROM refresh_tokens WHERE token_hash = $1", hashToken(req.RefreshToken))
    if err != nil {
        serverError(w, err, "logout delete token")
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func verifyEmail(w http.ResponseWriter, r *http.Request){
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    var req struct {
        Token string `json:"token"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        clientError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    _, err := db.Exec(ctx, `
        UPDATE users SET email_verified = TRUE, verification_token = NULL, verification_expires = NULL WHERE verification_token = $1 AND verification_expires > NOW()
    `, hashToken(req.Token))
    if err != nil {
        serverError(w, err, "verifyEmail update")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "email verified"})
}


func resendVerification(w http.ResponseWriter, r *http.Request){
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    var req struct {
        Email string `json:"email"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        clientError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    req.Email = strings.ToLower(strings.TrimSpace(req.Email))

    var userID int
    var verified bool
    err := db.QueryRow(ctx, "SELECT id, email_verified FROM users WHERE email = $1", req.Email).Scan(&userID, &verified)
    if err != nil {
        w.WriteHeader(http.StatusOK)
        return
    }

    if verified {
        clientError(w, http.StatusBadRequest, "email already verified")
        return
    }

    token, err := generateRefreshToken()
    if err != nil {
        serverError(w, err, "resendVerification token gen")
        return
    }

    expires := time.Now().Add(24 * time.Hour)

    _, err = db.Exec(ctx, "UPDATE users SET verification_token = $1, verification_expires = $2 WHERE id = $3", hashToken(token), expires, userID)
    if err != nil {
        serverError(w, err, "resendVerification update")
        return
    }

    if err := sendVerificationEmail(req.Email, token); err != nil {
        serverError(w, err, "resendVerification send email")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "verification email sent"})

}