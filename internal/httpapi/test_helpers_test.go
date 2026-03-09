package httpapi

import (
	"MiniJira/internal/store/memory"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func newTestHandler() http.Handler {
	store := memory.NewStore()

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	return NewMux(store, store, store, logger)
}

func performRequest(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	return w
}

func decodeJSON(t *testing.T, body io.Reader, target any) {
	err := json.NewDecoder(body).Decode(target)
	if err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
}

func createProject(t *testing.T, handler http.Handler, key, name string) ProjectResponse {
	body := `{"key":"` + key + `","name":"` + name + `"}`

	w := performRequest(t, handler, http.MethodPost, "/projects", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create project: %v", w.Code)
	}

	var project ProjectResponse
	decodeJSON(t, w.Body, &project)

	return project
}
