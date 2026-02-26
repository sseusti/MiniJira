package httpapi

import "net/http"

type ErrorResponse struct {
	Error string
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
	return
}
