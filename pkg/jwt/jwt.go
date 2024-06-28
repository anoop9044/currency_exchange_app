package jwt

import (
	"context"
	"net/http"
	"time"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"currency_exchange_app/internal/model"
)

type JWTMiddleware struct {
	SecretKey []byte
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{
		SecretKey: []byte(secret),
	}
}

func GenerateToken(user model.User, secretKey []byte) (string, error) {
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
func (m *JWTMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the Authorization header first
		authHeader := r.Header.Get("Authorization")
		var tokenStr string

		if authHeader != "" {
			// Header format should be "Bearer <token>"
			splitToken := strings.Split(authHeader, "Bearer ")
			if len(splitToken) != 2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr = splitToken[1]
		} else {
			// If Authorization header is not present, check the cookie
			tokenCookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr = tokenCookie.Value
		}

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return m.SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store the claims in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}