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
	nextIssueID   int
}

func (s *fakeStore) CreateIssue(i Issue) Issue {
	i.ID = s.nextIssueID
	s.nextIssueID++

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
		projects: map[string]Project{
			"PAY": {
				ID:   1,
				Key:  "PAY",
				Name: "Payments",
			},
		},
		issues:        []Issue{},
		nextProjectID: 2,
		nextIssueID:   1,
	}

	issue, err := CreateIssue(store, "PAY", "Fix checkout")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if issue.ID != 1 {
		t.Fatalf("expected issue ID 1, got %d", issue.ID)
	}

	if issue.ProjectKey != "PAY" {
		t.Fatalf("expected project key PAY, got %s", issue.ProjectKey)
	}

	if issue.Title != "Fix checkout" {
		t.Fatalf("expected title Fix checkout, got %s", issue.Title)
	}

	if issue.Status != StatusOpen {
		t.Fatalf("expected status %s, got %s", StatusOpen, issue.Status)
	}
}

func TestCreateIssue_InvalidInput(t *testing.T) {
	store := &fakeStore{
		projects: map[string]Project{
			"PAY": Project{
				ID:   1,
				Key:  "PAY",
				Name: "Payments",
			},
		},
		issues:        []Issue{},
		nextProjectID: 2,
		nextIssueID:   1,
	}

	tests := []struct {
		name       string
		projectKey string
		title      string
	}{
		{
			name:       "empty project key",
			projectKey: "",
			title:      "Payments",
		},
		{
			name:       "empty title",
			projectKey: "PAY",
			title:      "",
		},
		{
			name:       "blank project key",
			projectKey: "   ",
			title:      "Payments",
		},
		{
			name:       "blank title",
			projectKey: "PAY",
			title:      "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateIssue(store, tt.projectKey, tt.title)

			if !errors.Is(err, ErrInvalidIssue) {
				t.Fatalf("expected ErrInvalidIssue, got %v", err)
			}
		})
	}
}

func TestCreateIssue_ProjectNotFound(t *testing.T) {
	store := &fakeStore{
		projects:      map[string]Project{},
		issues:        []Issue{},
		nextProjectID: 1,
		nextIssueID:   1,
	}

	_, err := CreateIssue(store, "PAY", "Fix checkout")
	if !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestTransitionIssue_Success(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus string
		toStatus   string
	}{
		{
			name:       "open to in progress",
			fromStatus: StatusOpen,
			toStatus:   StatusInProgress,
		},
		{
			name:       "in progress to done",
			fromStatus: StatusInProgress,
			toStatus:   StatusDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStore{
				projects: map[string]Project{
					"PAY": {ID: 1, Key: "PAY", Name: "Payments"},
				},
				issues: []Issue{
					{
						ID:         1,
						ProjectKey: "PAY",
						Title:      "Fix checkout",
						Status:     tt.fromStatus,
					},
				},
				nextIssueID: 2,
			}

			_, err := TransitionIssue(store, 1, tt.toStatus)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if store.issues[0].Status != tt.toStatus {
				t.Fatalf("expected status %s, got %s",
					tt.toStatus,
					store.issues[0].Status)
			}

		})
	}
}

func TestTransitionIssue_InvalidTransition(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus string
		toStatus   string
	}{
		{
			name:       "open to done",
			fromStatus: StatusOpen,
			toStatus:   StatusDone,
		},
		{
			name:       "done to in progress",
			fromStatus: StatusDone,
			toStatus:   StatusInProgress,
		},
		{
			name:       "done to open",
			fromStatus: StatusDone,
			toStatus:   StatusOpen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStore{
				projects: map[string]Project{
					"PAY": {ID: 1, Key: "PAY", Name: "Payments"},
				},
				issues: []Issue{
					{
						ID:         1,
						ProjectKey: "PAY",
						Title:      "Fix checkout",
						Status:     tt.fromStatus,
					},
				},
				nextIssueID:   2,
				nextProjectID: 1,
			}

			_, err := TransitionIssue(store, 1, tt.toStatus)
			if !errors.Is(err, ErrInvalidTransition) {
				t.Fatalf("expected ErrInvalidTransition, got %v", err)
			}
		})
	}
}

func TestTransitionIssue_IssueNotFound(t *testing.T) {
	store := &fakeStore{
		projects: map[string]Project{
			"PAY": {ID: 1, Key: "PAY", Name: "Payments"},
		},
		issues:        []Issue{},
		nextIssueID:   1,
		nextProjectID: 2,
	}

	_, err := TransitionIssue(store, 999, StatusInProgress)
	if !errors.Is(err, ErrIssueNotFound) {
		t.Fatalf("expected ErrIssueNotFound, got %v", err)
	}
}

func TestTransitionIssue_InvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		issueID  int
		toStatus string
	}{
		{
			name:     "zero issue id",
			issueID:  0,
			toStatus: StatusOpen,
		},
		{
			name:     "invalid issue id",
			issueID:  -1,
			toStatus: StatusOpen,
		},
		{
			name:     "empty to status",
			issueID:  1,
			toStatus: "",
		},
		{
			name:     "blank status",
			issueID:  1,
			toStatus: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStore{
				projects: map[string]Project{
					"PAY": {ID: 1, Key: "PAY", Name: "Payments"},
				},
				issues: []Issue{
					{
						ID:         1,
						ProjectKey: "PAY",
						Title:      "Fix checkout",
						Status:     StatusOpen,
					},
				},
				nextProjectID: 2,
				nextIssueID:   2,
			}

			_, err := TransitionIssue(store, tt.issueID, tt.toStatus)
			if !errors.Is(err, ErrInvalidIssue) {
				t.Fatalf("expected ErrInvalidIssue, got %v", err)
			}
		})
	}
}
