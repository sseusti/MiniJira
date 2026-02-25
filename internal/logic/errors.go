package logic

import "errors"

var ErrInvalidProject = errors.New("invalid project")
var ErrProjectKeyExists = errors.New("project key already exists")
var ErrInvalidIssue = errors.New("invalid issue")
var ErrProjectNotFound = errors.New("project not found")
var ErrInvalidTransition = errors.New("invalid transition")
var ErrIssueNotFound = errors.New("issue not found")
var ErrInvalidID = errors.New("invalid id")
