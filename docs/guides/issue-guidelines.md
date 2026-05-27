# GitHub Issue Guidelines

## Title
Short phrase describing the problem or feature — no period at the end.

Format: `<area>: <what>`

Examples:
- `chart: add date range selector`
- `auth: redirect to original URL after login`
- `entry: show confirmation after successful submission`

## Body

```
## Overview
<One or two sentences describing the problem or goal. No implementation details.>

## Acceptance Criteria
- [ ] <Observable outcome — what the user sees or experiences>
- [ ] <Another outcome, if needed>

## Hints
- <Optional: point to a relevant file, template, or layer without prescribing a solution>
```

## Rules

**Scope** — one issue per logical change. Don't bundle unrelated work.

**Language** — write from the user's perspective. Describe behavior, not implementation. Avoid function names, SQL, or code snippets.

**Acceptance criteria** — each checkbox is something a person can verify by using the app (open X, do Y, expect Z). Not "update the service layer."

**Hints** — optional pointers to help orient whoever picks up the issue (e.g. "likely in the controller or template layer", "see `web/templates/components/`"). One or two at most. Never prescribe a solution.

**No estimates, labels, or assignees** — leave those to be set separately.
