package logic

import (
	"errors"
	"testing"
)

type projectStore struct {
	projects map[string]Project
	nextID   int
}

func (ps *projectStore) GetByKey(key string) (Project, bool) {
	p, ok := ps.projects[key]
	return p, ok
}

func (ps *projectStore) CreateProject(p Project) Project {
	p.ID = ps.nextID
	ps.nextID++
	ps.projects[p.Key] = p
	return p
}

func (ps *projectStore) List() []Project {
	list := make([]Project, 0, len(ps.projects))
	for _, p := range ps.projects {
		list = append(list, p)
	}
	return list
}

func TestCreateProject_Success(t *testing.T) {
	store := &projectStore{
		projects: make(map[string]Project),
		nextID:   1,
	}

	project, err := CreateProject(store, "PAY", "Payments")
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	if project.ID != 1 {
		t.Fatalf("project id should be 1, got %v", project.ID)
	}

	if project.Key != "PAY" {
		t.Fatalf("project key should be PAY, got %v", project.Key)
	}

	if project.Name != "Payments" {
		t.Fatalf("project name should be Payments, got %v", project.Name)
	}
}

func TestCreateProject_InvalidInput(t *testing.T) {
	store := &projectStore{
		projects: make(map[string]Project),
		nextID:   1,
	}

	tests := []struct {
		name        string
		key         string
		projectName string
	}{
		{
			name:        "empty name",
			key:         "PAY",
			projectName: "",
		},
		{
			name:        "empty key",
			key:         "",
			projectName: "Payments",
		},
		{
			name:        "blank name",
			key:         "PAY",
			projectName: "   ",
		},
		{
			name:        "blank key",
			key:         "    ",
			projectName: "Payments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateProject(store, tt.key, tt.projectName)
			if !errors.Is(err, ErrInvalidProject) {
				t.Fatalf("expected ErrInvalidProject, got %v", err)
			}
		})
	}
}

func TestCreateProject_DuplicateKey(t *testing.T) {
	store := &projectStore{
		projects: map[string]Project{"PAY": Project{}},
		nextID:   1,
	}

	tests := []struct {
		name        string
		key         string
		projectName string
	}{
		{
			name:        "duplicate key",
			key:         "PAY",
			projectName: "Payments",
		},
		//{
		//	name:        "case sensitive",
		//	key:         "pay",
		//	projectName: "Payments",
		//},
		// тут вопрос с регистром
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateProject(store, tt.key, tt.projectName)
			if !errors.Is(err, ErrProjectKeyExists) {
				t.Fatalf("expected ErrProjectKeyExists, got %v", err)
			}
		})
	}
}

type fakeStore struct {
	projects      map[string]Project
	issues        []Issue
	nextProjectID int
	nextIssueId   int
}

func (s *fakeStore) CreateIssue(i Issue) Issue {
	i.ID = s.nextIssueId
	s.nextIssueId++

	s.issues = append(s.issues, i)

	return i
}

func (s *fakeStore) GetIssueByID(id int) (Issue, bool) {
	for _, i := range s.issues {
		if i.ID == id {
			return i, true
		}
	}

	return Issue{}, false
}

func (s *fakeStore) UpdateIssueStatus(id int, newStatus string) (Issue, bool) {
	for i := range s.issues {
		if s.issues[i].ID == id {
			s.issues[i].Status = newStatus
			return s.issues[i], true
		}
	}

	return Issue{}, false
}

func (s *fakeStore) ListIssuesByProjectKey(projectKey string) []Issue {
	res := make([]Issue, 0, len(s.issues))
	for _, i := range s.issues {
		if i.ProjectKey == projectKey {
			res = append(res, i)
		}
	}

	return res
}

func (s *fakeStore) GetByKey(key string) (Project, bool) {
	p, ok := s.projects[key]
	return p, ok
}

func (s *fakeStore) CreateProject(p Project) Project {
	p.ID = s.nextProjectID
	s.nextProjectID++
	s.projects[p.Key] = p
	return p
}

func (s *fakeStore) List() []Project {
	list := make([]Project, 0, len(s.projects))
	for _, p := range s.projects {
		list = append(list, p)
	}
	return list
}

func TestCreateIssue_Success(t *testing.T) {
	store := &fakeStore{
		projects:      map[string]Project{"PAY": Project{}},
		issues:        []Issue{},
		nextProjectID: 1,
		nextIssueId:   1,
	}

	tests := []struct {
		name       string
		id         int
		projectKey string
		title      string
		status     string
	}{
		{
			name:       "project key exists",
			id:         1,
			projectKey: "PAY",
			title:      "Payments",
			status:     StatusOpen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateIssue(store, tt.projectKey, tt.title)
			if err != nil {
				t.Fatalf("expected no errors, got %v", err)
			}

			if tt.status != StatusOpen {
				t.Fatalf("expected status %s, got %s", StatusOpen, tt.status)
			}

			if tt.projectKey != "PAY" {
				t.Fatalf("project key should be PAY, got %v", tt.projectKey)
			}

			if tt.title != "Payments" {
				t.Fatalf("title should be Payments, got %v", tt.title)
			}

			if tt.id != 1 {
				t.Fatalf("id should be 1, got %v", tt.id)
			}
		})
	}
}
