package errors

import (
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}