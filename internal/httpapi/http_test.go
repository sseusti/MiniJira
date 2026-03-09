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

func testLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	return logger
}

func newTestHandler() http.Handler {
	store := memory.NewStore()
	logger := testLogger()
	return NewMux(store, store, store, logger)
}

func TestCreateProject_HTTP(t *testing.T) {
	handler := newTestHandler()

	body := `{"key":"PAY","name":"Payments"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/projects",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %v", w.Code)
	}

	var resp ProjectResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %v", resp.ID)
	}

	if resp.Key != "PAY" {
		t.Fatalf("expected key PAY, got %s", resp.Key)
	}

	if resp.Name != "Payments" {
		t.Fatalf("expected name Payments, got %s", resp.Name)
	}

	if ct := w.Header().Get("Content-Type"); ct == "" {
		t.Fatal("expected Content-Type header")
	}
}

func TestCreateProjet_HTTP_DuplicateKey(t *testing.T) {
	handler := newTestHandler()

	body := `{"key":"PAY","name":"Payments"}`

	req1 := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %v", w1.Code)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %v", w2.Code)
	}
}

func TestGetProjects_HTTP(t *testing.T) {
	handler := newTestHandler()

	body := `{"key":"PAY","name":"Payments"}`

	reqCreate := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")

	wCreate := httptest.NewRecorder()
	handler.ServeHTTP(wCreate, reqCreate)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %v", w.Code)
	}

	var resp []ProjectResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("expected 1 project, got %d", len(resp))
	}

	if resp[0].Key != "PAY" {
		t.Fatalf("expected key PAY, got %s", resp[0].Key)
	}
}
