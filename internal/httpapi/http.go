package httpapi

import (
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

func NewMux(s *memory.Store) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req CreateProjectRequest
			err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			created, err := logic.CreateProject(s, req.Key, req.Name)
			if errors.Is(err, logic.ErrInvalidProject) {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if errors.Is(err, logic.ErrProjectKeyExists) {
				w.WriteHeader(http.StatusConflict)
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

		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	})
	mux.HandleFunc("/issues", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var issue CreateIssueRequest
			err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&issue)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			created, err := logic.CreateIssue(s, issue.ProjectKey, issue.Title)
			if errors.Is(err, logic.ErrInvalidIssue) {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if errors.Is(err, logic.ErrProjectNotFound) {
				w.WriteHeader(http.StatusNotFound)
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
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			issues := s.ListIssuesByProjectKey(projectKey)
			WriteJSON(w, http.StatusOK, issues)
			return
		}
	})
	mux.HandleFunc("/issues/transition", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var issue TransitionIssueRequest
		err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&issue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updated, err := logic.TransitionIssue(s, issue.IssueID, issue.ToStatus)
		if errors.Is(err, logic.ErrInvalidIssue) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, logic.ErrIssueNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, logic.ErrInvalidTransition) {
			w.WriteHeader(http.StatusConflict)
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
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		issue, err := logic.GetIssue(s, id)
		if errors.Is(err, logic.ErrIssueNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, logic.ErrInvalidID) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		WriteJSON(w, http.StatusOK, issue)
		return
	})

	return mux
}
