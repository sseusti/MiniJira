# MiniJira

[English version](README.en.md)

Лёгкий учебный трекер задач на Go с REST API, in-memory хранилищем и Swagger-документацией.

## Что умеет

- создание и просмотр проектов
- создание и просмотр задач (issues)
- фильтрация задач по `project_key`
- перевод задачи по статусам: `OPEN -> IN_PROGRESS -> DONE`
- health-check endpoint

## Требования

- Go `1.25+`
- `make` (опционально)

## Быстрый старт

```bash
cp .env.example .env
make run
```

Или без `make`:

```bash
go run ./cmd/api
```

Сервер по умолчанию: `http://localhost:8080`

## Конфигурация

Через переменные окружения:

- `HTTP_PORT` — порт HTTP сервера (по умолчанию `8080`)
- `LOG_LEVEL` — `debug|info|warn|error` (по умолчанию `info`)
- `LOG_FORMAT` — `text|json` (по умолчанию `text`)

## API

### Основные маршруты

- `GET /health`
- `GET /projects`
- `POST /projects`
- `GET /issues?project_key=PAY`
- `POST /issues`
- `GET /issue?id=1`
- `POST /issues/transition`

### Swagger

- UI: `http://localhost:8080/swagger/index.html`
- Генерация/обновление: `make swag`

## Примеры запросов

Создать проект:

```bash
curl -X POST http://localhost:8080/projects \
  -H 'Content-Type: application/json' \
  -d '{"key":"PAY","name":"Payments"}'
```

Создать issue:

```bash
curl -X POST http://localhost:8080/issues \
  -H 'Content-Type: application/json' \
  -d '{"project_key":"PAY","title":"Fix checkout"}'
```

Перевести issue в `IN_PROGRESS`:

```bash
curl -X POST http://localhost:8080/issues/transition \
  -H 'Content-Type: application/json' \
  -d '{"issue_id":1,"to_status":"IN_PROGRESS"}'
```

## Архитектура (кратко)

- `cmd/api` — вход в приложение
- `internal/httpapi` — transport слой (HTTP handlers, middleware, mapping)
- `internal/usecase` — application/use-case слой
- `internal/logic` — доменные модели, правила и порты
- `internal/store/memory` — инфраструктурное in-memory хранилище
- `internal/config` — загрузка и валидация конфигурации

## Разработка

Полезные команды:

```bash
make test      # go test ./...
make fmt       # go fmt ./...
make swag      # сгенерировать swagger docs
make check     # fmt + test
```

## Тесты

```bash
go test ./...
```
