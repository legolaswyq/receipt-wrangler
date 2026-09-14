# Budgeting & Cash-Flow Dashboard — Design Spec

- **Date:** 2026-09-14
- **Status:** Approved (design), pending implementation plan
- **Components:** `api/` (backend), `desktop/` (Angular), `mobile/` (Flutter — client regen only)

## 1. Problem

The current dashboard cannot answer the questions the owner actually has:

- **Where is my money going, per category, against a plan?** There is a category pie chart, but
  it has no notion of a target/budget, so it can't say "you've used 75% of your Dining Out budget."
- **Income pollutes the expense view.** Income is modelled today as *negative-amount* receipts in
  an "Income" category. The app already reserves negative amounts for **refunds/credits** (the AI
  prompt says so verbatim — `services/prompts.go`), so the system genuinely cannot distinguish a
  salary deposit from a supermarket refund. Every expense chart is therefore computed against a
  total that includes a large negative number, distorting the proportions.

## 2. Goal

A single **Budget** dashboard widget, scoped to the current calendar month, that shows:

- A **summary row**: Income · Spent · Net (income − spent).
- A **progress bar per budgeted category**: `spent / target`, red when over budget.
- Spend in categories with no target rolled into a single **"Untracked"** line, so the Spent total
  reconciles.
- **Inline editing** of targets directly on the widget (click a target to change it; "+ Add a
  category target").

Targets are **recurring monthly** and reset on the 1st of each month.

## 3. Non-goals (v1 / YAGNI)

Explicitly out of scope for this build:

- Per-month target overrides, budget rollover, multi-month trend charts.
- A dedicated "Budgets" management page (editing is inline on the widget).
- Income *targets*/goals — income is measured, not budgeted.
- Overspend notifications/alerts.
- Cleaning up the existing seeded dummy data (see §9 caveat).

## 4. Decisions (locked during brainstorming)

| Decision | Choice |
|---|---|
| Income vs expense | New `Category.IsIncome` flag; amounts stay **positive** |
| Target period | Recurring monthly, resets on the 1st |
| Widget content | One widget: summary row + per-category bars |
| Editing targets | Inline on the widget (stored centrally) |
| Unbudgeted spend | Rolled into a single "Untracked" line |
| Income scope | Flag + show Income/Net in summary; no income targets |
| Seeded-data cleanup | Not in this build |

## 5. Data model (backend)

### 5.1 `Category.IsIncome bool`

- Global flag on the existing global `Category` model (`not null;default:false`). AutoMigrate adds the
  column; no data migration (existing categories default to expense).
- Semantics: a receipt in an income-flagged category counts toward **Income**; it is **excluded**
  from the expense progress bars and the **Spent** total. Amounts are expected positive.
- Rides the existing category create/update contract (`UpsertCategoryCommand` /
  `PUT /category/{categoryId}`), surfaced on `swagger.yml` `Category` + the upsert command.

### 5.2 `CategoryBudget` (new model)

```
CategoryBudget {
  BaseModel
  GroupId    uint            // not null
  CategoryId uint            // not null
  Amount     decimal.Decimal // not null; the recurring monthly target
}
// unique index on (GroupId, CategoryId) — one target per category per group
```

- **Per group.** Each group holds its own targets; categories are global but a budget is a
  per-group slice, mirroring how grants slice the global category pool.
- **The "All" group.** The dashboard in use lives on the virtual "All" aggregate group (id 2). Spend
  for the widget is computed with the **same receipt-scoping the pie chart already uses**
  (`GetPagedReceiptsByGroupId` with grants + paid-by resolver + all-group union), so the All group
  aggregates across the user's groups automatically. Targets set on the All group are stored against
  that group id like any other.
- Deleting a category cascades to its `CategoryBudget` rows (explicit delete inside
  `CategoryRepository.DeleteCategory`'s transaction, matching the existing custom-field/grant cascade
  conventions rather than relying on DB-level FKs). Deleting a group cascades likewise in
  `GroupService.DeleteGroup`.

## 6. Backend endpoints

### 6.1 Widget data — `POST /widget/budget/{groupId}`

Mirrors the pie-chart widget (`handlers/widgets.go`, `services/pie_chart.go`). Returns the current
month's figures:

```json
{
  "month": "2026-09",
  "income": 5000.00,
  "spent": 3240.00,
  "net": 1760.00,
  "categories": [
    { "categoryId": 7, "name": "Dining Out", "target": 400.00, "spent": 300.00, "over": false },
    { "categoryId": 3, "name": "Transport",  "target": 150.00, "spent": 180.00, "over": true  }
  ],
  "untracked": 890.00
}
```

- Reuses `GetPagedReceiptsByGroupId` with the caller's grant + `PaidByListResolver` scoping (so
  category/tag grants and paid-by visibility are honoured exactly as elsewhere) **plus a
  current-month date filter** (`[first-of-month, next-first-of-month)` in the app's timezone).
- `income` = sum of receipts whose category `IsIncome`; `spent` = sum of non-income receipt amounts;
  bars are built only for categories that have a `CategoryBudget`; all other non-income spend sums
  into `untracked`. A receipt with multiple categories fans out the same way the pie chart does
  (documented double-count behaviour) — acceptable and consistent.
- Restricted categories are substituted, not leaked, consistent with `pie_chart.go`.

### 6.2 Budget CRUD (inline edit)

- **Upsert:** `PUT /budget/{groupId}` with `{ categoryId, amount }` (or a small list) → creates or
  updates the `(group, category)` target. `amount <= 0` (or a delete verb) removes the target.
- **List (optional):** the widget-data endpoint already returns targets, so a separate list endpoint
  is not required for v1.
- Gated by the new permission set below.

### 6.3 Permissions

Add a group-scoped `group.budgets.*` set to the permission registry
(`internal/permissions/registry.go`) + the `swagger.yml` `Permission` enum:

- `group.budgets.read` — gates the widget-data endpoint.
- `group.budgets.create` / `group.budgets.update` / `group.budgets.delete` — gate the CRUD.

**Legacy Owner** (full group scope) picks these up automatically via add-only seed reconciliation;
Legacy Editor/Viewer do not (budget management is owner-level, matching member management). No data
migration. `registry_test.go` keeps the registry and swagger enum in sync.

### 6.4 Widget type

Add `BUDGET` to the `WidgetType` enum (`swagger.yml`), rendered by the desktop dashboard.

## 7. Frontend (desktop)

- New `BUDGET` entry in `widgetTypeOptions` (`dashboard/constants/widget-options.ts`) and the
  dashboard-form `ngSwitch` + dashboard renderer `ngSwitch`.
- New `budget/budget.component` (+ spec): renders the summary row and per-category bars from the
  widget-data endpoint; inline-edits a target (calls the CRUD endpoint, then refreshes).
- Progress bar: Material `mat-progress-bar` (or equivalent), red when `over`.
- Income-flag checkbox added to the category create/edit form.
- No new nav/route — everything is on the dashboard and the existing category management screen.

## 8. Mobile (client regen — CRITICAL SAFETY ITEM)

Mobile **does** render dashboards and **does** deserialize a **closed** `widget_type` enum
(`mobile/api/lib/src/model/widget_type.dart`) and the `Category` model. Per the two documented
production login outages, **a new value on a closed enum in a response model fails the whole
payload's deserialization on already-released builds.**

Required in the **same change** as the swagger edits:

1. Regenerate `mobile/api/` (and `desktop/src/open-api/`) from `swagger.yml`; re-apply the two
   documented dart-dio patches; `dart analyze` must report **0 errors**.
2. **Resolve the released-build hazard explicitly** (to be settled in the plan, not assumed):
   - Determine whether the *currently released* mobile build tolerates an unknown `WidgetType`
     (dart-dio `EnumClass` throws on unknown during deserialization → a dashboard containing a
     `BUDGET` widget would fail to load on old builds).
   - Preferred mitigation: give `WidgetType` a serializer fallback (unknown → a sentinel the mobile
     dashboard skips), or defer shipping the enum value until mobile handles it — mirroring the
     "untyped string" resolution used for `AppData` permissions. The budget widget itself is
     **desktop-only** for v1; mobile only needs to *not crash* when a dashboard contains one.

This item gates the change — it is the highest-risk part of the spec.

## 9. Known caveat — seeded income data

Because seeded-data cleanup is out of scope, the **Income** figure will reflect the existing
negative-amount income receipts (reading negative/incorrect) until the owner either re-seeds income
as positive + flags the "Income" category, or requests a one-time cleanup. The widget's *logic* is
correct; only the current dummy data is inconsistent with the new model. A follow-up cleanup script
(convert negative income receipts → positive, set `Income.IsIncome = true`) can be run on request.

## 10. Testing

- **Backend:** `CategoryBudget` repository round-trip + unique constraint + cascade on
  category/group delete; widget-data service (income/spent/net split, untracked rollover,
  current-month boundary, grant/paid-by scoping, over-budget flag); permission registry/enum sync;
  handler auth (read vs CRUD gates). Follow existing `main_test.go` patterns; remove stray `app.db`.
- **Desktop:** budget component unit tests (rendering bars, over state, inline edit calls the
  endpoint); widget-picker/renderer registration.
- **Mobile:** the enum-safety guard test (a dashboard payload carrying `BUDGET` must not crash the
  client), alongside the regen `dart analyze` gate.
- **Client-drift check:** the swagger/mobile timestamp check from the root `CLAUDE.md` must pass.

## 11. Rollout order (respecting the monorepo workflow)

1. Backend: model + migration + endpoints + permissions + swagger + tests.
2. Regenerate desktop **and** mobile clients from swagger in the same change; resolve the mobile
   enum-safety item; re-apply dart patches; analyze.
3. Desktop: widget + category income flag + tests.
4. Update `api/CLAUDE.md`, `desktop/CLAUDE.md`, `mobile/CLAUDE.md`, and root `CLAUDE.md` as needed.
5. Full test suites green before commit/push.
