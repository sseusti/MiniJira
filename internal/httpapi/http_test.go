package httpapi

import (
	"MiniJira/internal/logic"
	"fmt"
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
		t.Fatalf("expected ID 1, got %v", resp.ID)
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

func TestGetIssue_HTTP(t *testing.T) {
	handler := newTestHandler()
	createProject(t, handler, "PAY", "Payments")

	body := `{"project_key":"PAY","title":"Fix checkout"}`
	w := performRequest(t, handler, http.MethodPost, "/issues", body)
	if w.Code != http.StatusCreated {
		t.Fatalf(`expected status code 201, got %d`, w.Code)
	}

	var created IssueResponse
	decodeJSON(t, w.Body, &created)

	issueID := created.ID

	path := fmt.Sprintf("/issue?id=%d", issueID)
	w = performRequest(t, handler, http.MethodGet, path, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got IssueResponse
	decodeJSON(t, w.Body, &got)

	if got.ID != created.ID {
		t.Fatalf("expected id %d, got %d", created.ID, got.ID)
	}

	if got.ProjectKey != "PAY" {
		t.Fatalf("expected project key PAY, got %s", got.ProjectKey)
	}

	if got.Title != "Fix checkout" {
		t.Fatalf("expected title Fix checkout, got %s", got.Title)
	}

	if got.Status != logic.StatusOpen {
		t.Fatalf("expected status %s, got %s", logic.StatusOpen, got.Status)
	}
}

func TestListIssues_HTTP(t *testing.T) {
	handler := newTestHandler()

	createProject(t, handler, "PAY", "Payments")

	w := performRequest(t, handler, http.MethodPost, "/issues", `{"project_key":"PAY","title":"Fix checkout"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	w = performRequest(t, handler, http.MethodPost, "/issues", `{"project_key":"PAY","title":"Add retries"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	w = performRequest(t, handler, http.MethodGet, "/issues?project_key=PAY", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var issues []IssueResponse
	decodeJSON(t, w.Body, &issues)

	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}

	for _, issue := range issues {
		if issue.ProjectKey != "PAY" {
			t.Fatalf("expected project key PAY, got %s", issue.ProjectKey)
		}
	}
}

func TestTransitionIssue_HTTP(t *testing.T) {
	handler := newTestHandler()

	createProject(t, handler, "PAY", "Payments")

	w := performRequest(
		t,
		handler,
		http.MethodPost,
		"/issues",
		`{"project_key":"PAY","title":"Fix checkout"}`,
	)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var created IssueResponse
	decodeJSON(t, w.Body, &created)

	body := fmt.Sprintf(`{"issue_id":%d,"to_status":"IN_PROGRESS"}`, created.ID)

	w = performRequest(
		t,
		handler,
		http.MethodPost,
		"/issues/transition",
		body,
	)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var updated IssueResponse
	decodeJSON(t, w.Body, &updated)

	if updated.Status != logic.StatusInProgress {
		t.Fatalf("expected status %s, got %s",
			logic.StatusInProgress,
			updated.Status,
		)
	}
}

func TestCreateIssue_HTTP_ProjectNotFound(t *testing.T) {
	handler := newTestHandler()

	w := performRequest(
		t,
		handler,
		http.MethodPost,
		"/issues",
		`{"project_key":"UNKNOWN","title":"Fix checkout"}`,
	)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}

	var resp ErrorResponse
	decodeJSON(t, w.Body, &resp)

	if resp.Error != "not found" {
		t.Fatalf("expected error %q, got %q", "not found", resp.Error)
	}
}

func TestTransitionIssue_HTTP_InvalidTransition(t *testing.T) {
	handler := newTestHandler()

	createProject(t, handler, "PAY", "Payments")

	w := performRequest(
		t,
		handler,
		http.MethodPost,
		"/issues",
		`{"project_key":"PAY","title":"Fix checkout"}`,
	)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var created IssueResponse
	decodeJSON(t, w.Body, &created)

	body := fmt.Sprintf(`{"issue_id":%d,"to_status":"DONE"}`, created.ID)

	w = performRequest(
		t,
		handler,
		http.MethodPost,
		"/issues/transition",
		body,
	)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w.Code)
	}

	var resp ErrorResponse
	decodeJSON(t, w.Body, &resp)

	if resp.Error != "conflict" {
		t.Fatalf("expected error %q, got %q", "conflict", resp.Error)
	}
}

func TestGetIssue_HTTP_NotFound(t *testing.T) {
	handler := newTestHandler()

	w := performRequest(
		t,
		handler,
		http.MethodGet,
		"/issue?id=999",
		"",
	)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}

	var resp ErrorResponse
	decodeJSON(t, w.Body, &resp)

	if resp.Error != "not found" {
		t.Fatalf("expected error %q, got %q", "not found", resp.Error)
	}
}
