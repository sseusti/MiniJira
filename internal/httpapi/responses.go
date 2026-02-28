package httpapi

import "net/http"

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
	return
}
