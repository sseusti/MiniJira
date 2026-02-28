package httpapi

import (
	_ "MiniJira/docs"
	"MiniJira/internal/httpapi/middleware"
	"MiniJira/internal/logic"
	"encoding/json"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)

	return
}

type CreateProjectRequest struct {
	Key  string `json:"key" example:"PAY"`
	Name string `json:"name" example:"Payments"`
}

type CreateIssueRequest struct {
	ProjectKey string `json:"project_key" example:"PAY"`
	Title      string `json:"title" example:"Fix checkout validation"`
}

type TransitionIssueRequest struct {
	IssueID  int    `json:"issue_id" example:"1"`
	ToStatus string `json:"to_status" example:"IN_PROGRESS"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewMux(projectStore logic.ProjectStore, issueStore logic.IssueStore, piStore logic.ProjectIssueStore) http.Handler {
	h := NewHandler(projectStore, issueStore, piStore)
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/projects", h.Projects)
	mux.HandleFunc("/issues", h.Issues)
	mux.HandleFunc("/issues/transition", h.IssuesTransition)
	mux.HandleFunc("/issue", h.Issue)

	return middleware.Logger(mux)
}
