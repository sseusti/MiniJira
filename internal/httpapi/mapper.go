package httpapi

import "MiniJira/internal/logic"

func toProjectResponse(p logic.Project) ProjectResponse {
	return ProjectResponse{
		ID:   p.ID,
		Key:  p.Key,
		Name: p.Name,
	}
}

func toProjectResponses(ps []logic.Project) []ProjectResponse {
	res := make([]ProjectResponse, len(ps))
	for i, p := range ps {
		res[i] = toProjectResponse(p)
	}

	return res
}

func toIssueResponse(i logic.Issue) IssueResponse {
	return IssueResponse{
		ID:         i.ID,
		ProjectKey: i.ProjectKey,
		Title:      i.Title,
		Status:     i.Status,
	}
}

func toIssueResponses(is []logic.Issue) []IssueResponse {
	res := make([]IssueResponse, len(is))
	for i, p := range is {
		res[i] = toIssueResponse(p)
	}

	return res
}
