package api

import (
	"encoding/json"
	"net/http"
	"time"

	"currency_exchange_app/config"
	"currency_exchange_app/internal/model"
	jwt "currency_exchange_app/pkg/jwt"
)

func NewJWTMiddleware() *jwt.JWTMiddleware {
	return &jwt.JWTMiddleware{
		SecretKey: []byte(config.Config.JwtSecret),
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// In a real application, you would validate the user's credentials
	var user model.User
	user.Username = "admin"
	user.Role = "admin"

	// Generate JWT token
	tokenString, err := jwt.GenerateToken(user, []byte(config.Config.JwtSecret))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Set JWT token in a cookie (optional)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
	})

	// Prepare JSON response with token
	response := map[string]string{"token": tokenString}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Could not marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}