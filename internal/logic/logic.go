package logic

import (
	"strings"
)

func CreateProject(store ProjectStore, key, name string) (Project, error) {
	key = strings.TrimSpace(key)
	name = strings.TrimSpace(name)

	if key == "" || name == "" {
		return Project{}, ErrInvalidProject
	}

	_, ok := store.GetByKey(key)
	if ok {
		return Project{}, ErrProjectKeyExists
	}

	created := store.CreateProject(Project{Key: key, Name: name})
	return created, nil
}

type ProjectIssueStore interface {
	ProjectStore
	IssueStore
}

func CreateIssue(store ProjectIssueStore, projectKey, title string) (Issue, error) {
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
		Status:     StatusOpen,
	}

	created := store.CreateIssue(issue)

	return created, nil
}

func TransitionIssue(store IssueStore, issueID int, toStatus string) (Issue, error) {
	toStatus = strings.TrimSpace(toStatus)
	if toStatus == "" || issueID <= 0 {
		return Issue{}, ErrInvalidIssue
	}

	issue, ok := store.GetIssueByID(issueID)
	if !ok {
		return Issue{}, ErrIssueNotFound
	}
	ok = isAllowed(issue.Status, toStatus)
	if !ok {
		return Issue{}, ErrInvalidTransition
	}

	updated, ok := store.UpdateIssueStatus(issue.ID, toStatus)
	if !ok {
		return Issue{}, ErrIssueNotFound
	}

	return updated, nil
}

func isAllowed(status string, toStatus string) bool {
	if status == StatusOpen && toStatus == StatusInProgress {
		return true
	}
	if status == StatusInProgress && toStatus == StatusDone {
		return true
	}

	return false
}

func GetIssue(store IssueStore, id int) (Issue, error) {
	if id <= 0 {
		return Issue{}, ErrInvalidID
	}

	issue, ok := store.GetIssueByID(id)
	if !ok {
		return Issue{}, ErrIssueNotFound
	}

	return issue, nil
}
