package main

import "sync"

type Store struct {
	issues      []Issue
	projects    []Project
	mu          sync.RWMutex
	nextID      int
	nextIssueID int
}

func NewStore() *Store {
	return &Store{nextID: 1, nextIssueID: 1}
}

func (s *Store) Create(p Project) Project {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextID
	s.nextID++
	s.projects = append(s.projects, p)

	return p
}

func (s *Store) List() []Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projects := make([]Project, len(s.projects))
	copy(projects, s.projects)

	return projects
}

func (s *Store) GetByKey(key string) (Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.projects {
		if p.Key == key {
			return p, true
		}
	}

	return Project{}, false
}

func (s *Store) CreateIssue(i Issue) Issue {
	s.mu.Lock()
	defer s.mu.Unlock()

	i.ID = s.nextIssueID
	s.nextIssueID++
	s.issues = append(s.issues, i)

	return i
}

func (s *Store) GetIssueByID(id int) (Issue, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, i := range s.issues {
		if i.ID == id {
			return i, true
		}
	}

	return Issue{}, false
}

func (s *Store) UpdateIssueStatus(id int, newStatus string) (Issue, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.issues {
		if s.issues[i].ID == id {
			s.issues[i].Status = newStatus
			return s.issues[i], true
		}
	}

	return Issue{}, false
}

func (s *Store) ListIssuesByProjectKey(projectKey string) []Issue {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Issue, 0, len(s.issues))
	for _, i := range s.issues {
		if i.ProjectKey == projectKey {
			res = append(res, i)
		}
	}

	return res
}
