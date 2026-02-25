package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var ErrInvalidProject = errors.New("invalid project")
var ErrProjectKeyExists = errors.New("project key already exists")
var ErrInvalidIssue = errors.New("invalid issue")
var ErrProjectNotFound = errors.New("project not found")

type Project struct {
	ID   int
	Key  string
	Name string
}

type Store struct {
	issues      []Issue
	projects    []Project
	mu          sync.RWMutex
	nextID      int
	nextIssueID int
}

func (s *Store) CreateIssue(i Issue) Issue {
	s.mu.Lock()
	defer s.mu.Unlock()

	i.ID = s.nextIssueID
	s.nextIssueID++
	s.issues = append(s.issues, i)

	return i
}

func (s *Store) Create(p Project) Project {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextID
	s.nextID++
	s.projects = append(s.projects, p)

	return p
}

func (s *Store) List() []Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projects := make([]Project, len(s.projects))
	copy(projects, s.projects)

	return projects
}

func (s *Store) GetByKey(key string) (Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.projects {
		if p.Key == key {
			return p, true
		}
	}

	return Project{}, false
}

func (s *Store) ListIssuesByProjectKey(projectKey string) []Issue {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Issue, 0, len(s.issues))
	for _, i := range s.issues {
		if i.ProjectKey == projectKey {
			res = append(res, i)
		}
	}

	return res
}

func NewStore() *Store {
	return &Store{nextID: 1, nextIssueID: 1}
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)

	return
}

func main() {
	s := NewStore()

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

			created, err := CreateProject(s, req.Key, req.Name)
			if errors.Is(err, ErrInvalidProject) {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if errors.Is(err, ErrProjectKeyExists) {
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

			created, err := CreateIssue(s, issue.ProjectKey, issue.Title)
			if errors.Is(err, ErrInvalidIssue) {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if errors.Is(err, ErrProjectNotFound) {
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

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type CreateProjectRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func CreateProject(store *Store, key, name string) (Project, error) {
	key = strings.TrimSpace(key)
	name = strings.TrimSpace(name)

	if key == "" || name == "" {
		return Project{}, ErrInvalidProject
	}

	_, ok := store.GetByKey(key)
	if ok {
		return Project{}, ErrProjectKeyExists
	}

	created := store.Create(Project{Key: key, Name: name})
	return created, nil
}

type Issue struct {
	ID         int
	ProjectKey string
	Title      string
	Status     string
}

func CreateIssue(store *Store, projectKey, title string) (Issue, error) {
	projectKey = strings.TrimSpace(projectKey)
	title = strings.TrimSpace(title)

	if projectKey == "" || title == "" {
		return Issue{}, ErrInvalidIssue
	}

	_, ok := store.GetByKey(projectKey)
	if !ok {
		return Issue{}, ErrProjectNotFound
	}

	issue := Issue{
		ProjectKey: projectKey,
		Title:      title,
		Status:     "OPEN",
	}

	created := store.CreateIssue(issue)

	return created, nil
}

type CreateIssueRequest struct {
	ProjectKey string `json:"project_key"`
	Title      string `json:"title"`
}
