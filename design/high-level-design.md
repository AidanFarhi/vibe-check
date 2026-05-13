# Vibe Check — High-Level Design

## Overview

Vibe Check is a personal health tracking tool that lets users log daily wellness metrics and visualize trends over time. Users submit a simple daily entry rating five dimensions of wellbeing on a 1–10 scale and optionally add a free-text note for the day. The app then surfaces those readings as charts and summaries to help users spot patterns.

### Tracked Metrics

| Metric        | Scale | Description                          |
|---------------|-------|--------------------------------------|
| Depression    | 1–10  | Severity of depressive feelings      |
| Happiness     | 1–10  | Overall sense of joy or positivity   |
| Pain          | 1–10  | Physical pain or discomfort          |
| Energy        | 1–10  | Vitality and motivation level        |
| Sleep Quality | 1–10  | Quality of the previous night's sleep|

---

## Tech Stack

| Layer      | Technology                        |
|------------|-----------------------------------|
| Language   | Go (Golang)                       |
| Database   | PostgreSQL                        |
| Frontend   | HTMX + Go HTML Templates          |
| Charting   | Chart.js                          |
| Container  | Docker                            |

---

## Architecture

The app follows a Clean Architecture, organized into concentric layers where dependencies only point inward.

```
┌──────────────────────────────────────────┐
│            middleware (HTTP)             │  Request pipeline: logging, error handling
├──────────────────────────────────────────┤
│     controller (HTTP) + view (DTOs)      │  Handlers, request parsing, view-models
├──────────────────────────────────────────┤
│            service (use cases)           │  Business logic, orchestration, validation
├──────────────────────────────────────────┤
│       domain (entities + interfaces)     │  Entry entity, EntryRepository interface
├──────────────────────────────────────────┤
│            repo (persistence)            │  Postgres implementation of domain interfaces
└──────────────────────────────────────────┘
```

Dependency graph:

```
repo ← domain ← service ← controller
```

`domain` has no outward dependencies. `service` and `repo` both depend on `domain`. `controller` depends on `service`. `repo` is a plug-in detail that `service` never directly touches.

### Layer Responsibilities

**middleware**
- Sits in front of all controllers in the HTTP pipeline
- Handles cross-cutting concerns: request logging, panic recovery, error formatting
- Each middleware is a standard `http.Handler` wrapper; composed on the router

**controller**
- HTTP handlers wired to routes
- Parses incoming requests (form values, query params)
- Calls service layer, builds a view-model, and renders the appropriate template
- Returns HTMX partials for dynamic swaps or full-page responses on first load
- Owns a `view/` sub-package of plain Go structs that carry data into HTML templates — no business logic, purely a data-transfer shape for the template renderer

**service**
- Application use cases: `SubmitEntry`, `GetEntries`
- Enforces business rules (one entry per day, metric range 1–10)
- Depends on `domain.EntryRepository` interface, never on `repo` directly

**domain**
- `Entry` entity with validation rules (metric values 1–10, date required)
- `EntryRepository` interface — the contract that `service` consumes and `repo` satisfies
- No imports from any other internal package

**repo**
- Postgres implementation of `domain.EntryRepository`
- No business logic — only data access

---

## Data Model

```sql
CREATE TABLE "user" (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT        NOT NULL UNIQUE,
    password_hash   TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE entry (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES "user"(id),
    date       DATE        NOT NULL,
    depression SMALLINT    NOT NULL CHECK (depression  BETWEEN 1 AND 10),
    happiness  SMALLINT    NOT NULL CHECK (happiness   BETWEEN 1 AND 10),
    pain       SMALLINT    NOT NULL CHECK (pain        BETWEEN 1 AND 10),
    energy     SMALLINT    NOT NULL CHECK (energy      BETWEEN 1 AND 10),
    sleep      SMALLINT    NOT NULL CHECK (sleep       BETWEEN 1 AND 10),
    note       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, date)
);
```

---

## Key User Flows

### 1. Submit Daily Entry
1. User opens the home page; a form pre-filled with today's date is shown.
2. User sets sliders/inputs for each metric, optionally writes a note for the day, and submits.
3. HTMX posts to `POST /entries`; the controller calls `SubmitEntry` on the service.
4. On success, the server returns an updated chart partial via HTMX swap.
5. On validation error, the form partial is re-rendered with inline errors.

### 2. View History / Charts
1. User selects a date range (default: last 30 days).
2. HTMX fires `GET /entries?from=...&to=...`; the controller calls `GetEntries` on the service.
3. Server renders a line-chart partial (using Chart.js, driven by data embedded in the HTML).

---

## Project Layout

```
vibe-check/
├── cmd/
│   └── server/
│       └── main.go            # Entry point, wires dependencies
├── internal/
│   ├── middleware/
│   │   └── middleware.go      # Logging, recovery, error handler
│   ├── controller/
│   │   ├── view/
│   │   │   └── entry.go       # View-model structs for templates
│   │   ├── home.go            # Home page handler
│   │   └── entry.go           # Entry submit / history handlers
│   ├── service/
│   │   └── entry.go           # SubmitEntry, GetEntries use cases
│   ├── domain/
│   │   ├── entry.go           # Entry entity + validation rules
│   │   └── repository.go      # EntryRepository interface
│   └── repo/
│       └── entry.go           # Postgres implementation of EntryRepository
├── db/
│   ├── migrations/            # Numbered SQL migration files (up/down)
│   │   ├── 001_create_entries.up.sql
│   │   └── 001_create_entries.down.sql
│   └── migrate.go             # Migration runner (golang-migrate or custom)
├── web/
│   ├── templates/
│   │   ├── base.html
│   │   └── pages/
│   │       └── home.html
│   ├── css/
│   └── js/
├── Dockerfile
└── design-doc/
    └── high-level-design.md
```

---