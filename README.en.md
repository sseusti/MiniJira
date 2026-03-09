# MiniJira

A lightweight educational issue tracker in Go with REST API, in-memory storage, and Swagger docs.

## Features

- create and list projects
- create and fetch issues
- list issues filtered by `project_key`
- issue status transitions: `OPEN -> IN_PROGRESS -> DONE`
- health-check endpoint

## Requirements

- Go `1.25+`
- `make` (optional)

## Quick start

```bash
cp .env.example .env
make run
```

Or without `make`:

```bash
go run ./cmd/api
```

Default server address: `http://localhost:8080`

## Configuration

Environment variables:

- `HTTP_PORT` — HTTP server port (default: `8080`)
- `LOG_LEVEL` — `debug|info|warn|error` (default: `info`)
- `LOG_FORMAT` — `text|json` (default: `text`)

## API

### Main routes

- `GET /health`
- `GET /projects`
- `POST /projects`
- `GET /issues?project_key=PAY`
- `POST /issues`
- `GET /issue?id=1`
- `POST /issues/transition`

### Swagger

- UI: `http://localhost:8080/swagger/index.html`
- Regenerate docs: `make swag`

## Request examples

Create project:

```bash
curl -X POST http://localhost:8080/projects \
  -H 'Content-Type: application/json' \
  -d '{"key":"PAY","name":"Payments"}'
```

Create issue:

```bash
curl -X POST http://localhost:8080/issues \
  -H 'Content-Type: application/json' \
  -d '{"project_key":"PAY","title":"Fix checkout"}'
```

Transition issue to `IN_PROGRESS`:

```bash
curl -X POST http://localhost:8080/issues/transition \
  -H 'Content-Type: application/json' \
  -d '{"issue_id":1,"to_status":"IN_PROGRESS"}'
```

## Architecture (short)

- `cmd/api` — application entrypoint
- `internal/httpapi` — transport layer (HTTP handlers, middleware, mapping)
- `internal/usecase` — application/use-case layer
- `internal/logic` — domain models, rules, and ports
- `internal/store/memory` — in-memory infrastructure storage
- `internal/config` — config loading and validation

## Development

Useful commands:

```bash
make test      # go test ./...
make fmt       # go fmt ./...
make swag      # generate swagger docs
make check     # fmt + test
```

## Tests

```bash
go test ./...
```
