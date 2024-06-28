package api

import (
	"net/http"

	"golang.org/x/time/rate"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(1), 1) // Allow 5 requests per second

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded, Please try again Later", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}