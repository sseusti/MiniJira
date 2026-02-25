package logic

type Project struct {
	ID   int
	Key  string
	Name string
}

type Issue struct {
	ID         int
	ProjectKey string
	Title      string
	Status     string
}

const (
	StatusOpen       = "OPEN"
	StatusInProgress = "IN_PROGRESS"
	StatusDone       = "DONE"
)
