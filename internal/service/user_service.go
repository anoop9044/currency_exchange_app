package service

import (
	"errors"
	"currency_exchange_app/internal/model"
)

var users = map[string]model.User{
	"admin": {Username: "admin", Role: "admin"},
	"user":  {Username: "user", Role: "user"},
}

func Authenticate(username, password string) (model.User, error) {
	user, exists := users[username]
	if !exists {
		return model.User{}, errors.New("user not found")
	}

	// In a real application, you would check the password here
	return user, nil
}