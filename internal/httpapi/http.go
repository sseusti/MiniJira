package httpapi

import (
	"MiniJira/internal/httpapi/middleware"
	"MiniJira/internal/logic"
	"MiniJira/internal/store/memory"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)

	return
}

type CreateProjectRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type CreateIssueRequest struct {
	ProjectKey string `json:"project_key"`
	Title      string `json:"title"`
}

type TransitionIssueRequest struct {
	IssueID  int    `json:"issue_id"`
	ToStatus string `json:"to_status"`
}

func NewMux(s *memory.Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req CreateProjectRequest
			err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&req)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "invalid request")
				return
			}

			created, err := logic.CreateProject(s, req.Key, req.Name)
			if errors.Is(err, logic.ErrInvalidProject) {
				WriteError(w, http.StatusBadRequest, "invalid request")
				return
			} else if errors.Is(err, logic.ErrProjectKeyExists) {
				WriteError(w, http.StatusConflict, "conflict")
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			WriteJSON(w, http.StatusCreated, created)
			return
		}

		if r.Method == http.MethodGet {
			p := s.List()
			WriteJSON(w, http.StatusOK, p)
			return
		}

		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	})
	mux.HandleFunc("/issues", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var issue CreateIssueRequest
			err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&issue)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "invalid request")
				return
			}

			created, err := logic.CreateIssue(s, issue.ProjectKey, issue.Title)
			if errors.Is(err, logic.ErrInvalidIssue) {
				WriteError(w, http.StatusBadRequest, "invalid request")
				return
			} else if errors.Is(err, logic.ErrProjectNotFound) {
				WriteError(w, http.StatusNotFound, "not found")
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			WriteJSON(w, http.StatusCreated, created)
			return
		}
		if r.Method == http.MethodGet {
			projectKey := r.URL.Query().Get("project_key")

			projectKey = strings.TrimSpace(projectKey)
			if projectKey == "" {
				WriteError(w, http.StatusBadRequest, "invalid request")
				return
			}

			issues := s.ListIssuesByProjectKey(projectKey)
			WriteJSON(w, http.StatusOK, issues)
			return
		}
	})
	mux.HandleFunc("/issues/transition", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var issue TransitionIssueRequest
		err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&issue)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "invalid request")
			return
		}

		updated, err := logic.TransitionIssue(s, issue.IssueID, issue.ToStatus)
		if errors.Is(err, logic.ErrInvalidIssue) {
			WriteError(w, http.StatusBadRequest, "invalid request")
			return
		} else if errors.Is(err, logic.ErrIssueNotFound) {
			WriteError(w, http.StatusNotFound, "not found")
			return
		} else if errors.Is(err, logic.ErrInvalidTransition) {
			WriteError(w, http.StatusConflict, "conflict")
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		WriteJSON(w, http.StatusOK, updated)
		return
	})
	mux.HandleFunc("/issue", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		idStr := r.URL.Query().Get("id")
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			WriteError(w, http.StatusBadRequest, "invalid request")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			WriteError(w, http.StatusBadRequest, "invalid request")
			return
		}

		issue, err := logic.GetIssue(s, id)
		if errors.Is(err, logic.ErrIssueNotFound) {
			WriteError(w, http.StatusNotFound, "not found")
			return
		} else if errors.Is(err, logic.ErrInvalidID) {
			WriteError(w, http.StatusBadRequest, "invalid request")
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		WriteJSON(w, http.StatusOK, issue)
		return
	})

	return middleware.Logger(mux)
}
