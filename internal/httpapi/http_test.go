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
	if err := json.NewEncoder(w.Body).Encode(&resp); err != nil {
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

}
