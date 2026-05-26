# Pull Request Guidelines

## Title
Short phrase summarizing the scope of changes — no period at the end.

Format: `<area>: <what changed>`

Example: `UI updates: modal handle close, autofill fix, and docs consolidation`

## Body

```
## Summary
- <bullet per logical change>
- <keep each bullet to one sentence>

## Test plan
- [ ] <specific step to verify each change works>
- [ ] <cover both the happy path and any edge cases touched>

Closes #<issue number>
```

Include the `Closes #<issue number>` line only when the PR resolves a GitHub issue. Omit it otherwise.

### Summary rules
- One bullet per logical change — don't group unrelated changes into one line
- Describe what changed and why (e.g. "Fix browser autofill overriding input background color on auth pages"), not just what file changed

### Test plan rules
- Each checkbox maps to a specific, observable action (open X, click Y, expect Z)
- Cover the golden path and any edge cases introduced by the change
- If a change has no UI surface, describe how to verify correctness (e.g. run a specific test, check a log)
