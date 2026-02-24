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

type Project struct {
	ID   int
	Key  string
	Name string
}

type Store struct {
	projects []Project
	mu       sync.RWMutex
	nextID   int
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

func NewStore() *Store {
	return &Store{nextID: 1}
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
