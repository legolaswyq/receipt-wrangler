# Receipt Wrangler — "Modern Fintech" UI Redesign

**Date:** 2026-09-18
**Component:** `desktop/` (Angular 19 web app) — no backend or mobile changes
**Status:** Design approved (pending spec review)

## Goal

Redesign the desktop web UI to a **modern-fintech** visual language — data-forward,
money-first, crisp — and add a **light + dark theme** toggle the app currently lacks.
Delivered at the **design-system level** (tokens + Material theme + app shell + shared-ui
components) so the new look cascades across all 153 components, followed by a **sweep** that
migrates every hardcoded-color SCSS file onto the new tokens so both themes render correctly
everywhere.

This is a **purely visual/theming** effort: no component behavior, routing, form flow, or
data-logic changes.

## Background / Current State

- **Foundation:** `desktop/src/styles.scss` (global + M2 Material theme, **light-only**) and
  `desktop/src/variables.scss` (SCSS token maps: palettes, shadow/spacing/radius scales).
- **Shell:** `src/layout/header/` + `src/layout/sidebar/` define the top bar and drawer nav.
- **Shared-ui:** `src/button`, `src/input`, `src/base-input`, `src/shared-ui/base-table`,
  `src/avatar`, plus global `.mat-mdc-card` styling in `styles.scss`.
- **Scale:** 153 components / 152 component SCSS files. **37** component SCSS files hardcode
  colors (**314** hex + **88** rgba occurrences); reports is the heaviest hardcoder. 21 files
  reference `white`/`black` keywords. 26 files already `@use` `variables.scss`.
- **Stack:** Angular Material M2 (`mat.m2-define-light-theme`, `mat.all-component-themes`),
  Bootstrap SCSS, Inter + Raleway fonts, NGXS (with persistent storage plugin already in use).

## Design

### 1. Token architecture (single source of truth)

Introduce a **semantic CSS custom-property token layer** — the redesign's foundation. Tokens
are defined once on `:root` (light) and overridden under `:root[data-theme="dark"]`, plus a
`@media (prefers-color-scheme: dark)` block that applies dark values when no explicit
`data-theme` is set. Components consume tokens; they never hardcode colors.

New file `desktop/src/_tokens.scss` (imported into `styles.scss`), defining CSS variables:

**Color (light → dark):**
| Token | Light | Dark |
|---|---|---|
| `--rw-bg` | `#F7F8FA` | `#0B0F17` |
| `--rw-surface` | `#FFFFFF` | `#151B26` |
| `--rw-surface-2` | `#F1F3F6` | `#1C2330` |
| `--rw-border` | `#E4E7EC` | `#262E3D` |
| `--rw-text` | `#0B1220` | `#E6EAF2` |
| `--rw-text-muted` | `#667085` | `#9AA4B2` |
| `--rw-accent` | `#4F46E5` | `#6366F1` |
| `--rw-accent-contrast` | `#FFFFFF` | `#FFFFFF` |
| `--rw-accent-hover` | `#4338CA` | `#818CF8` |
| `--rw-positive` | `#16A34A` | `#22C55E` |
| `--rw-negative` | `#DC2626` | `#F87171` |
| `--rw-warning` | `#D97706` | `#F59E0B` |
| `--rw-elevation-1` | `0 1px 2px rgba(16,24,40,.06)` | `0 1px 2px rgba(0,0,0,.4)` |
| `--rw-elevation-2` | `0 4px 12px rgba(16,24,40,.08)` | `0 4px 14px rgba(0,0,0,.5)` |

**Non-color tokens** (theme-independent, kept from current scale): `--rw-radius-sm: .375rem`,
`--rw-radius-md: .625rem`, `--rw-radius-lg: .75rem` (cards), spacing tokens mirror the existing
`$spacing-*` values.

The existing SCSS maps in `variables.scss` are retained for the Material theme definition (M2
needs SCSS palettes), but their **semantic accent/primary values are aligned to the token
palette** so Material output and token output agree.

### 2. Material dual theme

In `styles.scss`:
- Keep `mat.m2-define-light-theme` for the light theme (primary palette re-centered on indigo
  `#4F46E5`).
- Add `mat.m2-define-dark-theme` with the same primaries/typography.
- Emit light component themes globally; emit **dark color overrides** under
  `:root[data-theme="dark"]` via `mat.all-component-colors($dark-theme)` (typography emitted
  once, unscoped). This is the standard M2 dual-theme pattern and keeps CSS size bounded
  (colors-only override, not a full second theme).

### 3. ThemeService + toggle

- New `desktop/src/services/theme.service.ts`: holds a `ThemeMode` signal (`'light' | 'dark' |
  'system'`), writes `data-theme` to `document.documentElement`, and listens to the
  `prefers-color-scheme` media query while in `'system'` mode.
- Persistence via the **existing NGXS persistent storage** (a small `ThemeState`, mirroring how
  auth/preferences persist) so the choice survives reload; initial value is `'system'`.
- **Toggle in the header** (`src/layout/header/`), placed next to the notifications button: a
  single icon button cycling light → dark → system, with a tooltip showing the current mode.

### 4. Fintech visual language

- **Money-first typography:** apply `font-variant-numeric: tabular-nums` app-wide to numeric
  contexts; add a display treatment (larger size, tighter letter-spacing, heavier weight) for
  headline figures — dashboard totals, budget amounts, the spending-table total. Body/UI stays
  Inter; Raleway is retained only where already used for branding.
- **Cards:** `--rw-radius-lg` (12px), 1px `--rw-border`, `--rw-elevation-1` at rest / `-2` on
  hover, `--rw-surface` background. Replace the current translucent-black border + lift.
- **Data surfaces:** tighter table row height, clearer uppercase-muted column headers,
  right-aligned tabular numeric columns (`base-table`).
- **Shell:** header uses `--rw-surface` with a bottom `--rw-border`; sidebar uses `--rw-surface`
  with accent-driven active-group dot and active-nav state; add-FAB uses the accent.
- **Charts:** pie/spending-table slice colors continue to come from the existing category-color
  palette (unchanged data logic); chart chrome (labels, gridlines, empty states) uses tokens so
  it reads in both themes.

### 5. Sweep (dark-mode correctness)

Migrate the **37 component SCSS files that hardcode colors** onto tokens — replace literal
hex/rgba/`white`/`black` used for backgrounds, text, and borders with the matching semantic
token. Purely decorative brand illustrations and the category-color palette are exempt (they are
intentional fixed colors, not theme surfaces). Priority order by occurrence count: reports
cluster first, then receipts/roles/auth, then the rest.

### 6. Scope boundaries

- **No** behavior, routing, form-flow, or data-logic changes.
- **No** edits to generated clients (`desktop/src/open-api/`).
- **No** backend or mobile changes.
- Bootstrap stays; we do not remove or replace the CSS framework.
- Category-color palette and chart slice colors keep their fixed values.

## Testing & Regression Strategy

- `npm run build` stays clean; existing `npm test` specs pass (ThemeService gets its own spec).
- **Playwright screenshot sweep in BOTH themes** across key screens — dashboard, receipts list,
  receipt detail/form, categories, reports, settings, auth/login — using the system Chrome
  harness already established (`channel:'chrome'`, `NODE_PATH` to desktop node_modules,
  `domcontentloaded`). This is the primary guard against contrast failures and missed hardcoded
  colors, which are the main risk of a dual-theme rollout.
- Manual verification checklist: theme toggle persists across reload; `system` mode follows OS;
  no white-on-white or black-on-black; accent/active states legible on both grounds; budget
  over/under and income/spend semantic colors correct in both themes.

## Risks & Mitigations

- **Missed hardcoded color → broken dark screen.** Mitigated by the sweep phase + the two-theme
  screenshot sweep.
- **Material dark override CSS bloat.** Mitigated by emitting colors-only dark override, not a
  full second theme.
- **Selector-specificity collisions** when moving global card/button styles to tokens. Mitigated
  by keeping the existing selectors and only swapping their values to `var(--rw-*)`.
- **Regression in the budget-dashboard work on this branch** (unmerged). The redesign builds on
  the current branch; screenshot sweep includes the dashboard to confirm no regression.

## Out of Scope / Future

- Mobile (Flutter) redesign.
- Restructuring layouts or information architecture.
- Removing Bootstrap or migrating Material M2 → M3.
