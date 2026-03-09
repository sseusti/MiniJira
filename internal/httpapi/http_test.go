package httpapi

import (
	"MiniJira/internal/logic"
	"net/http"
	"testing"
)

func TestCreateProject_HTTP(t *testing.T) {
	handler := newTestHandler()

	body := `{"key":"PAY","name":"Payments"}`

	w := performRequest(t, handler, http.MethodPost, "/projects", body)
	if w.Code != http.StatusCreated {
		t.Fatalf(`expected status code 201, got %d`, w.Code)
	}

	var resp ProjectResponse
	decodeJSON(t, w.Body, &resp)

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

func TestCreateProject_HTTP_DuplicateKey(t *testing.T) {
	handler := newTestHandler()

	body := `{"key":"PAY","name":"Payments"}`

	createProject(t, handler, "PAY", "Payments")

	w := performRequest(t, handler, http.MethodPost, "/projects", body)
	if w.Code != http.StatusConflict {
		t.Fatalf(`expected status code 201, got %d`, w.Code)
	}
}

func TestGetProjects_HTTP(t *testing.T) {
	handler := newTestHandler()
	createProject(t, handler, "PAY", "Payments")

	w := performRequest(t, handler, http.MethodGet, "/projects", "")
	if w.Code != http.StatusOK {
		t.Fatalf(`expected status code 200, got %d`, w.Code)
	}

	var resp []ProjectResponse
	decodeJSON(t, w.Body, &resp)

	if len(resp) != 1 {
		t.Fatalf("expected 1 project, got %d", len(resp))
	}

	if resp[0].Key != "PAY" {
		t.Fatalf("expected key PAY, got %s", resp[0].Key)
	}
}

func TestCreateIssue_HTTP(t *testing.T) {
	handler := newTestHandler()
	createProject(t, handler, "PAY", "Payments")

	body := `{"project_key":"PAY","title":"Fix checkout"}`
	w := performRequest(t, handler, http.MethodPost, "/issues", body)

	if w.Code != http.StatusCreated {
		t.Fatalf(`expected status code 201, got %d`, w.Code)
	}

	var resp IssueResponse
	decodeJSON(t, w.Body, &resp)

	if resp.ProjectKey != "PAY" {
		t.Fatalf("expected project key PAY, got %s", resp.ProjectKey)
	}

	if resp.Title != "Fix checkout" {
		t.Fatalf("expected title Fix checkout, got %s", resp.Title)
	}

	if resp.Status != logic.StatusOpen {
		t.Fatalf("expected status open, got %s", resp.Status)
	}
}
