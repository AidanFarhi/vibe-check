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
| Deployment | Fly.io                            |

---

## Architecture

The app follows a Clean Architecture, organized into concentric layers where dependencies only point inward.

```
┌──────────────────────────────────────────┐
│            middleware (HTTP)             │  Request pipeline: logging, error handling, auth
├──────────────────────────────────────────┤
│     controller (HTTP) + view (DTOs)      │  Handlers, request parsing, view-models
├──────────────────────────────────────────┤
│            service (use cases)           │  Business logic, orchestration, validation
├──────────────────────────────────────────┤
│       domain (entities + interfaces)     │  Entry/User/Session entities, repository interfaces
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
- Handles cross-cutting concerns: request logging, panic recovery
- `RequireAuth` validates the session cookie against the database and injects the user ID into the request context; unauthenticated requests are redirected to `/login`
- Each middleware is a standard `http.Handler` wrapper; composed on the router

**controller**
- HTTP handlers wired to routes
- Parses incoming requests (form values, query params)
- Calls service layer, builds a view-model, and renders the appropriate template
- Returns HTMX partials for dynamic swaps or full-page responses on first load
- Owns a `view/` sub-package of plain Go structs that carry data into HTML templates — no business logic, purely a data-transfer shape for the template renderer

**service**
- Application use cases: `SubmitEntry` (implemented), `GetEntries` (planned)
- Auth use cases: `Register`, `Login`, `Logout`
- Enforces business rules (one entry per day, metric range 1–10, password strength)
- Depends on domain repository interfaces, never on `repo` directly

**domain**
- `Entry`, `User`, `Session` entities with validation rules
- `EntryRepository`, `UserRepository`, `SessionRepository` interfaces — the contracts that `service` consumes and `repo` satisfies
- Sentinel errors: `ErrEmailTaken`, `ErrDuplicateEntry`
- No imports from any other internal package

**repo**
- Postgres implementations of all domain repository interfaces
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

CREATE TABLE session (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    token      TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
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

### 1. Register
1. User navigates to `/register` and submits email + password.
2. Server validates (password ≥ 8 chars, email not already taken) and creates the user.
3. On success, redirects to `/login`. On error, re-renders the register page with an inline message.

### 2. Login / Logout
1. User submits credentials at `/login`; server verifies password hash and creates a 30-day session cookie.
2. Authenticated requests carry the cookie; `RequireAuth` middleware validates it on every protected route.
3. `POST /logout` deletes the session and clears the cookie, redirecting to `/login`.

### 3. Submit Daily Entry
1. User opens the home page; the "Log Today" button opens a modal with sliders for each metric and an optional note.
2. User sets sliders (1–10) and optionally writes a note, then submits.
3. HTMX posts to `POST /entries`; the controller calls `SubmitEntry` on the service.
4. On success, the modal closes (HTMX swaps in the closed modal state).
5. On validation error or duplicate-entry (already logged today), the modal re-renders open with an inline error message.

### 4. View History / Charts *(planned)*
1. User selects a date range (default: last 30 days).
2. HTMX fires `GET /entries?from=...&to=...`; the controller calls `GetEntries` on the service.
3. Server renders a line-chart partial (using Chart.js, driven by data embedded in the HTML).

---

## Project Layout

```
vibe-check/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point, wires dependencies
├── internal/
│   ├── middleware/
│   │   └── middleware.go              # Logging, recovery, RequireAuth
│   ├── controller/
│   │   ├── view/
│   │   │   └── entry.go              # View-model structs for templates
│   │   ├── auth.go                   # Register, Login, Logout handlers
│   │   ├── home.go                   # Home page handler
│   │   └── entry.go                  # Entry submit handler
│   ├── service/
│   │   ├── auth.go                   # Register, Login, Logout use cases
│   │   └── entry.go                  # SubmitEntry use case
│   ├── domain/
│   │   ├── entry.go                  # Entry entity + validation
│   │   ├── user.go                   # User entity + validation
│   │   ├── session.go                # Session entity
│   │   └── repository.go             # EntryRepository, UserRepository, SessionRepository interfaces + sentinel errors
│   └── repo/
│       ├── entry.go                  # Postgres EntryRepository
│       ├── user.go                   # Postgres UserRepository
│       └── session.go                # Postgres SessionRepository
├── db/
│   ├── migrations/
│   │   ├── 001_create_user.up.sql
│   │   ├── 001_create_user.down.sql
│   │   ├── 002_create_session.up.sql
│   │   ├── 002_create_session.down.sql
│   │   ├── 003_create_entry.up.sql
│   │   └── 003_create_entry.down.sql
│   ├── setup/                        # One-time DB role/database creation scripts
│   └── migrate.go                    # Migration runner
├── web/
│   ├── templates/
│   │   ├── base.html
│   │   ├── pages/
│   │   │   ├── home.html
│   │   │   ├── login.html
│   │   │   └── register.html
│   │   └── components/
│   │       ├── chart-card.html
│   │       ├── log-button.html
│   │       ├── log-modal.html
│   │       ├── metric-tiles.html
│   │       ├── metrics-card.html
│   │       ├── navbar.html
│   │       ├── note-card.html
│   │       ├── page-header.html
│   │       ├── recent-strip.html
│   │       └── today-card.html
│   ├── css/
│   │   ├── base.css
│   │   ├── auth.css
│   │   └── home.css
│   └── js/
│       └── home.js
├── design/
│   └── high-level-design.md
├── fly.toml
└── Dockerfile
```
