package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"currency_exchange_app/internal/service"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Currency Exchange Application"))
}


func UpdateExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only authenticated admins can update exchange rates
	username := r.Context().Value("username").(string)
	role := r.Context().Value("role").(string)

	if role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Currency string  `json:"currency"`
		Rate     float64 `json:"rate"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	service.UpdateExchangeRate(r.Context(), username, req.Currency, req.Rate)

	// Respond with a success message
	response := map[string]string{
		"message": "Exchange rate updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}



func GetHistoricalRatesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Extract currency from the request path parameters
	vars := mux.Vars(r)
	currency := vars["currency"]

	// Call service to get historical rates
	rates, err := service.GetHistoricalRates(ctx, currency)
	if err != nil {
		http.Error(w, "Failed to fetch historical rates", http.StatusInternalServerError)
		return
	}

	// Encode response as JSON and send it
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}


