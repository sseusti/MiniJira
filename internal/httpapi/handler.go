package httpapi

import (
	"MiniJira/internal/httpapi/middleware"
	"MiniJira/internal/logic"
	"MiniJira/internal/usecase"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type ProjectResponse struct {
	ID   int    `json:"id" example:"1"`
	Key  string `json:"key" example:"PAY"`
	Name string `json:"name" example:"Payments"`
}

type IssueResponse struct {
	ID         int    `json:"id" example:"10"`
	ProjectKey string `json:"project_key" example:"PAY"`
	Title      string `json:"title" example:"Fix checkout validation"`
	Status     string `json:"status" example:"OPEN" enums:"OPEN,IN_PROGRESS,DONE"`
}

type Handler struct {
	service *usecase.Service
	logger  *logrus.Logger
}

func NewHandler(service *usecase.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Health godoc
// @Summary Health check
// @Description Check service availability
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.ListProjects(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.CreateProject(w, r)
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	return
}

// ListProjects godoc
// @Summary List projects
// @Tags projects
// @Produce json
// @Success 200 {array} ProjectResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects [get]
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	p := h.service.ListProjects()
	WriteJSON(w, http.StatusOK, toProjectResponses(p))
	return
}

// CreateProject godoc
// @Summary Create project
// @Description Create a new project with unique key
// @Tags projects
// @Accept json
// @Produce json
// @Param request body CreateProjectRequest true "Project payload"
// @Success 201 {object} ProjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects [post]
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	created, err := h.service.CreateProject(req.Key, req.Name)
	if errors.Is(err, logic.ErrInvalidProject) {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	} else if errors.Is(err, logic.ErrProjectKeyExists) {
		WriteError(w, http.StatusConflict, "conflict")
		return
	} else if err != nil {
		h.logger.WithFields(logrus.Fields{
			"rid": middleware.GetRequestID(r),
			"op":  "create_project",
		}).WithError(err).Error("operation failed")
		WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	WriteJSON(w, http.StatusCreated, toProjectResponse(created))
	return
}

// CreateIssue godoc
// @Summary Create issue
// @Description Create an issue in existing project
// @Tags issues
// @Accept json
// @Produce json
// @Param request body CreateIssueRequest true "Issue payload"
// @Success 201 {object} IssueResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /issues [post]
func (h *Handler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	var issue CreateIssueRequest
	err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&issue)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	created, err := h.service.CreateIssue(issue.ProjectKey, issue.Title)
	if errors.Is(err, logic.ErrInvalidIssue) {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	} else if errors.Is(err, logic.ErrProjectNotFound) {
		WriteError(w, http.StatusNotFound, "not found")
		return
	} else if err != nil {
		h.logger.WithFields(logrus.Fields{
			"rid": middleware.GetRequestID(r),
			"op":  "create_issue",
		}).WithError(err).Error("operation failed")
		WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	WriteJSON(w, http.StatusCreated, toIssueResponse(created))
	return
}

// ListIssues godoc
// @Summary List issues by project key
// @Description Returns issues for a project (filter is required)
// @Tags issues
// @Produce json
// @Param project_key query string true "Project key"
// @Success 200 {array} IssueResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /issues [get]
func (h *Handler) ListIssues(w http.ResponseWriter, r *http.Request) {
	issues, err := h.service.ListIssues(r.URL.Query().Get("project_key"))
	if errors.Is(err, logic.ErrInvalidIssue) {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	WriteJSON(w, http.StatusOK, toIssueResponses(issues))
	return
}

func (h *Handler) Issues(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.ListIssues(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.CreateIssue(w, r)
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	return
}

// GetIssue godoc
// @Summary Get issue by id
// @Description Returns issue by ID
// @Tags issues
// @Produce json
// @Param id query int true "Issue ID"
// @Success 200 {object} IssueResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /issue [get]
func (h *Handler) GetIssue(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	issue, err := h.service.GetIssue(id)
	if errors.Is(err, logic.ErrIssueNotFound) {
		WriteError(w, http.StatusNotFound, "not found")
		return
	} else if errors.Is(err, logic.ErrInvalidID) {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	} else if err != nil {
		h.logger.WithFields(logrus.Fields{
			"rid": middleware.GetRequestID(r),
			"op":  "get_issue",
		}).WithError(err).Error("operation failed")
		WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	WriteJSON(w, http.StatusOK, toIssueResponse(issue))
	return
}

func (h *Handler) Issue(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetIssue(w, r)
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	return
}

// TransitionIssue godoc
// @Summary Transition issue status
// @Description Change issue status following allowed transitions
// @Tags issues
// @Accept json
// @Produce json
// @Param request body TransitionIssueRequest true "Transition payload"
// @Success 200 {object} IssueResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /issues/transition [post]
func (h *Handler) TransitionIssue(w http.ResponseWriter, r *http.Request) {
	var issue TransitionIssueRequest
	err := json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&issue)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	updated, err := h.service.TransitionIssue(issue.IssueID, issue.ToStatus)
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
		h.logger.WithFields(logrus.Fields{
			"rid": middleware.GetRequestID(r),
			"op":  "transition_issue",
		}).WithError(err).Error("operation failed")
		WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	WriteJSON(w, http.StatusOK, toIssueResponse(updated))
	return
}

func (h *Handler) IssuesTransition(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.TransitionIssue(w, r)
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	return
}
