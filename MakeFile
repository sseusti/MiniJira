APP_NAME := mini-jira
CMD := ./cmd/api
SWAG_MAIN := main.go
SWAG_DIRS := cmd/api,internal/httpapi

.PHONY: run
run:
	go run $(CMD)

.PHONY: test
test:
	go test ./...

.PHONY: swag
swag:
	swag init --generalInfo $(SWAG_MAIN) --dir $(SWAG_DIRS)

.PHONY: swag-fmt
swag-fmt:
	swag fmt --dir $(SWAG_DIRS)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: check
check: fmt test

.PHONY: dev
dev: swag run