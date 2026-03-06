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
		{
			name:        "case sensitive",
			key:         "pay",
			projectName: "Payments",
		},
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
