package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"currency_exchange_app/config"
	"currency_exchange_app/internal/api"
	"currency_exchange_app/internal/redis"
	"currency_exchange_app/pkg/jwt"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize Redis
	err := redis.InitRedis(config.Config.RedisAddr)
	if err != nil {
		log.Fatalf("Could not initialize Redis: %v", err)
	}

	// Create a new router
	r := mux.NewRouter()

	// JWT Middleware
	jwtMiddleware := jwt.NewJWTMiddleware(config.Config.JwtSecret)

	// Define routes
	r.HandleFunc("/",api.HomeHandler).Methods("GET")
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")
	r.Handle("/ws", http.HandlerFunc(api.WebSocketHandler)).Methods("GET")
	r.Handle("/historical-rates/{currency}", jwtMiddleware.Handler(http.HandlerFunc(api.GetHistoricalRatesHandler))).Methods("GET")
	r.Handle("/updateExchangeRate", jwtMiddleware.Handler(api.RateLimitMiddleware(http.HandlerFunc(api.UpdateExchangeRateHandler)))).Methods("POST")

	// Start the server
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}