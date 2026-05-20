# Vibe Check

A personal wellness tracker for logging daily metrics and spotting trends over time. Each day you rate five dimensions of wellbeing on a 1–10 scale and optionally add a note.

**Stack:** Go · PostgreSQL · HTMX · HTML Templates · Docker · Fly.io

## Getting started

**Prerequisites:** Go 1.22+, PostgreSQL

1. Copy the environment file and fill in your database URL:
   ```
   cp .env.example .env
   ```

2. Run the server:
   ```
   go run ./cmd/server
   ```
   The app will automatically run pending migrations on startup and listen on `http://localhost:8088`.

3. For live reload during development (requires [air](https://github.com/air-verse/air)):
   ```
   air
   ```

## Environment

| Variable       | Description                  |
|----------------|------------------------------|
| `DATABASE_URL` | PostgreSQL connection string |

## Project structure

```
cmd/server/        Entry point
internal/
  controller/      HTTP handlers + view models (controller/view/)
  service/         Business logic / use cases
  domain/          Entities, interfaces, sentinel errors
  repo/            Postgres implementations
  middleware/      Logging, recovery, auth
db/
  migrations/      Numbered SQL migrations (run automatically on startup)
  setup/           One-time role/database creation scripts
web/
  templates/       Go HTML templates (base, pages, components)
  css/
  js/
```

## Tracked metrics

| Metric        | Scale |
|---------------|-------|
| Energy        | 1–10  |
| Sleep Quality | 1–10  |
| Happiness     | 1–10  |
| Pain          | 1–10  |
| Depression    | 1–10  |
