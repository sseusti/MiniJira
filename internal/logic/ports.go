package logic

type ProjectStore interface {
	GetByKey(key string) (Project, bool)
	CreateProject(p Project) Project
}

type IssueStore interface {
	CreateIssue(i Issue) Issue
	GetIssueByID(id int) (Issue, bool)
	UpdateIssueStatus(id int, newStatus string) (Issue, bool)
}
