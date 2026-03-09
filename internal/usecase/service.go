package usecase

import (
	"MiniJira/internal/logic"
	"strings"
)

type Service struct {
	projectStore logic.ProjectStore
	issueStore   logic.IssueStore
	piStore      logic.ProjectIssueStore
}

func NewService(projectStore logic.ProjectStore, issueStore logic.IssueStore, piStore logic.ProjectIssueStore) *Service {
	return &Service{
		projectStore: projectStore,
		issueStore:   issueStore,
		piStore:      piStore,
	}
}

func (s *Service) ListProjects() []logic.Project {
	return s.projectStore.List()
}

func (s *Service) CreateProject(key, name string) (logic.Project, error) {
	return logic.CreateProject(s.projectStore, key, name)
}

func (s *Service) CreateIssue(projectKey, title string) (logic.Issue, error) {
	return logic.CreateIssue(s.piStore, projectKey, title)
}

func (s *Service) ListIssues(projectKey string) ([]logic.Issue, error) {
	projectKey = strings.TrimSpace(projectKey)
	if projectKey == "" {
		return nil, logic.ErrInvalidIssue
	}

	return s.issueStore.ListIssuesByProjectKey(projectKey), nil
}

func (s *Service) GetIssue(id int) (logic.Issue, error) {
	return logic.GetIssue(s.issueStore, id)
}

func (s *Service) TransitionIssue(issueID int, toStatus string) (logic.Issue, error) {
	return logic.TransitionIssue(s.issueStore, issueID, toStatus)
}
