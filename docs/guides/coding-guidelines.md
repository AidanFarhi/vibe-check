# Go Coding Guidelines

## Naming
- Receivers: single lowercase letter matching the type (`h` for `*Home`, `r` for `*UserRepo`)
- Constructors: always `New<Type>` (e.g. `NewUserRepo`, `NewAuth`)
- Interfaces: noun + role suffix — `UserRepository`, not `IUserRepository` or `UserRepo`
- Unexported structs for internal view models (`authView`), exported for template-bound ones (`ModalView`)
- Context keys: typed constant using a private string type — `type contextKey string`

## Error Handling
- Define sentinel errors at the package boundary where they originate (`domain` for data errors, `service` for business rule errors)
- Wrap with context using `fmt.Errorf("operation: %w", err)` — always add a short verb phrase describing what failed
- Check domain sentinels with `errors.Is()`, driver-specific errors (e.g. pq) with `errors.As()`
- Controllers handle user-facing messaging: map known errors to specific messages, unknown errors to a generic fallback with a 500
- Never return raw DB errors to the caller above `repo`

## Structs and Constructors
- Structs hold only what they need — `repo` types get `*sql.DB`, services get interfaces, controllers get templates and services
- Depend on interfaces, not concrete types, at the service layer
- Keep constructors to one line when possible: `func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }`
- No public fields on internal structs

## Layering Rules
- Dependencies only point inward: `repo → domain ← service ← controller`
- `domain` has zero internal imports — no `service`, `repo`, or `controller` imports allowed
- `service` never imports `repo` directly — only domain interfaces
- `repo` never contains business logic — data access only

## HTTP Handlers
- Use `http.StatusSeeOther` (303) for all POST-redirect-GET flows
- Parse form values with `r.FormValue()` — no third-party form libraries
- Cookies: always set `HttpOnly: true`, `SameSite: http.SameSiteStrictMode`
- Return HTMX partials for dynamic interactions; full-page renders only on first load

## Validation
- Normalize input in domain constructors (lowercase, trim whitespace) before validating
- Range and presence checks belong in domain; business rules (password strength, uniqueness) belong in service
- Never validate the same constraint in more than one layer

## General
- No `init()` functions
- No global mutable state — wire dependencies in `main.go`
- Prefer table-driven tests when multiple input/output cases exist
- Parameterized SQL queries only (`$1`, `$2`) — no string concatenation in queries
