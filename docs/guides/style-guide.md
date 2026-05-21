# UI Style Guide

## Colors
Always use CSS variables — never hardcode hex values for backgrounds, text, or borders.

| Variable         | Use                                      |
|------------------|------------------------------------------|
| `--bg`           | Page background                          |
| `--surface`      | Card / container background             |
| `--surface-2`    | Inputs, nested surfaces                 |
| `--border`       | All borders                              |
| `--text`         | Primary text                             |
| `--text-muted`   | Labels, secondary text                  |
| `--accent`       | Interactive elements, focus states      |

Metric colors (energy, sleep, depression, etc.) are inline styles on data-bound elements — keep them co-located with the data, not in CSS classes.

## Spacing
- Card padding: `clamp(16px, 4vw, 24px)`
- Gap between cards: `14px` margin-bottom
- Form row spacing: `20px` margin-bottom
- Don't introduce new spacing values — reuse existing ones

## Typography
- Font stack: `-apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif` (inherited from body — don't re-declare)
- Use `clamp()` for any font size that should scale with viewport
- Uppercase labels: `font-size: 10–11px`, `letter-spacing: 0.6–0.9px`, `font-weight: 500`
- Never use `font-weight` below 500 for UI text

## Border Radius
- Cards and modals: `16px`
- Buttons and inputs: `10–12px`
- Pill buttons / badges: `20px`
- Don't mix radius values within the same component type

## Components

**Cards** — background `var(--surface)`, `1px solid var(--border)`, `border-radius: 16px`, responsive padding, `margin-bottom: 14px`

**Buttons** — primary uses `var(--accent)` background with dark text; secondary uses transparent background with `var(--border)` border; hover state is `opacity: 0.88` + `translateY(-2px)`

**Inputs / Textareas** — background `var(--surface-2)`, border `var(--border)`, focus border `var(--accent)`, `border-radius: 10px`

**Error messages** — red background at 10% opacity, red border, `font-size: 13px`

## Responsive Approach
- Mobile-first. Use `clamp()` for fluid sizing — avoid hard breakpoints unless layout fundamentally changes.
- Existing breakpoints: `400px` (3-col metric grid), `640px` (6-col), `1024px` (center content with side padding)
- Don't add new breakpoints without a layout reason

## Transitions
- Default interactive elements: `0.15s`
- Card-level emphasis changes: `0.4s`
- Always transition `color`, `background`, `border-color`, `opacity` — not `all`
