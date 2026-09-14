# Budget Dashboard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a monthly Budget dashboard widget with recurring per-category spend targets, progress bars, and income-vs-expense separation.

**Architecture:** A new global `Category.IsIncome` flag distinguishes income from spending. A new per-group `CategoryBudget` model stores recurring monthly targets. A `BudgetService` computes current-month figures by reusing the pie chart's receipt-scoping (grants + paid-by + all-group union) and filtering to the current month in Go. A new `BUDGET` widget type and desktop component render the summary + progress bars with inline target editing.

**Tech Stack:** Go 1.25 (Chi, GORM, shopspring/decimal), Angular 19 (standalone components, signals, generated `WidgetService`), Flutter (client regen only).

## Global Constraints

- **Never hand-edit generated clients** (`desktop/src/open-api/`, `mobile/api/`) except the two documented dart-dio patches.
- **Regenerate `mobile/api/` in the SAME change as any `swagger.yml` edit.** A new closed-enum value on a response model breaks already-released Android builds.
- **Every `swagger.yml` edit** → regenerate desktop + mobile clients; re-apply dart patches (`user_preferences.dart` `..quickScanDefaultStatus = ReceiptStatus.OPEN`, `system_settings.dart` `..currencyDisplay = r'`); `dart analyze` in `mobile/api` must report **0 errors** (73 warnings baseline).
- **macOS client regen:** use the pinned jar directly (`/tmp/openapi-generator-cli-7.10.0.jar`), not `npx` (root-owned npm prefix → EACCES).
- **Permission registry and swagger `Permission` enum must stay in sync** (`permissions/registry_test.go` enforces it).
- **Money is `shopspring/decimal`.** Compare with `decimal.Equal`/`StringFixed`, never `reflect.DeepEqual`. Serialize decimals as strings on the wire (match existing `quantity`/`unitPrice` nullable-string convention).
- **Amounts are positive.** Negative = refund/credit (existing convention). Income is identified by category flag, not sign.
- **Run the relevant test suite after each task.** Remove any stray `app.db` in test dirs before rerunning Go tests. Never touch `sqlite/`.
- **Local dev env:** Go server runs via `air` on `:8081` (SQLite `wrangler.db`, Redis 127.0.0.1:6379, `ENCRYPTION_KEY`/`SECRET_KEY=dev-*`, `CHROMIUM_BINARY_PATH` set). `go test` needs `CHROMIUM_BINARY_PATH` for the two PDF tests.
- **Permissions decided for this feature:** `group.budgets.read` (widget data), `group.budgets.update` (upsert a target), `group.budgets.delete` (remove a target). No separate `create` — upsert covers it. Legacy Owner picks all three up automatically (`LegacyGroupOwnerKeys` = every group-scope key).

---

## File Structure

**Backend (create):**
- `api/internal/models/category_budget.go` — the `CategoryBudget` GORM model.
- `api/internal/commands/upsert_category_budget_command.go` — upsert command + validation.
- `api/internal/repositories/category_budgets.go` — data access.
- `api/internal/services/budget.go` — `BudgetService.GetBudgetData`.
- `api/internal/structs/budget_data.go` — response structs.
- `api/internal/handlers/budgets.go` — HTTP handlers.
- `api/internal/routers/budget.go` — routes.
- Test files alongside each.

**Backend (modify):**
- `api/internal/models/category.go` — add `IsIncome`.
- `api/internal/commands/upsert_category_command.go` — add `IsIncome`.
- `api/internal/handlers/categories.go` — pass `IsIncome` through create/update.
- `api/internal/repositories/categories.go` — persist `IsIncome` on update; cascade-delete budgets.
- `api/internal/permissions/registry.go` — 3 new descriptors + constants.
- `api/internal/repositories/models.go` (or wherever `MakeMigrations` lists models) — register `CategoryBudget`.
- `api/internal/services/groups.go` — cascade-delete budgets on group delete.
- `api/internal/routers/router.go` — mount budget router.
- `api/swagger.yml` — `Category.isIncome`, `UpsertCategoryCommand.isIncome`, `BUDGET` widget type, `BudgetData`/`BudgetCategory`/`UpsertCategoryBudgetCommand` schemas, budget paths, 3 permission enum values.

**Frontend desktop (create):**
- `desktop/src/dashboard/budget/budget.component.{ts,html,scss,spec.ts}`.

**Frontend desktop (modify):**
- `desktop/src/dashboard/constants/widget-options.ts` — add BUDGET option.
- `desktop/src/dashboard/dashboard/dashboard.component.html` + `.ts` — render BUDGET.
- `desktop/src/dashboard/dashboard-form/dashboard-form.component.html` — BUDGET config case (name only).
- The category form component — income checkbox.

**Mobile (regen only):** `mobile/api/**` regenerated; one guard test.

**Docs:** `api/CLAUDE.md`, `desktop/CLAUDE.md`, `mobile/CLAUDE.md`, root `CLAUDE.md`.

---

## Task 1: `Category.IsIncome` flag (model + command + handler + repo)

**Files:**
- Modify: `api/internal/models/category.go`
- Modify: `api/internal/commands/upsert_category_command.go`
- Modify: `api/internal/handlers/categories.go`
- Modify: `api/internal/repositories/categories.go`
- Test: `api/internal/repositories/categories_test.go`

**Interfaces:**
- Produces: `models.Category.IsIncome bool`; `commands.UpsertCategoryCommand.IsIncome bool`.

- [ ] **Step 1: Write the failing test** — append to `api/internal/repositories/categories_test.go`:

```go
func TestUpdateCategoryPersistsIsIncome(t *testing.T) {
	defer repositories.TruncateTestDb()
	repo := repositories.NewCategoryRepository(nil)

	created, err := repo.CreateCategory(models.Category{Name: "Salary"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	created.IsIncome = true
	updated, err := repo.UpdateCategory(&created)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if !updated.IsIncome {
		t.Errorf("expected IsIncome true after update, got false")
	}

	fetched, err := repo.GetCategoryById(utils.UintToString(created.ID), false)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !fetched.IsIncome {
		t.Errorf("expected IsIncome persisted, got false")
	}
}
```

Check the exact repository method names in `api/internal/repositories/categories.go` first (`CreateCategory`, `UpdateCategory`, `GetCategoryById` signatures) and adjust the call sites to match; the assertion (IsIncome round-trips) stays.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd api && go test ./internal/repositories/ -run TestUpdateCategoryPersistsIsIncome -count=1`
Expected: FAIL (`IsIncome` undefined, or false after update).

- [ ] **Step 3: Add the model field** — in `api/internal/models/category.go`, add to the `Category` struct:

```go
IsIncome bool `gorm:"not null;default:false" json:"isIncome"`
```

- [ ] **Step 4: Add the command field** — in `api/internal/commands/upsert_category_command.go`, add to `UpsertCategoryCommand`:

```go
IsIncome bool `json:"isIncome"`
```

- [ ] **Step 5: Wire it through the handler** — in `api/internal/handlers/categories.go`, wherever `CreateCategory` / `UpdateCategory` build the `models.Category` from the request, set `IsIncome`. Read the current handler body; it currently unmarshals into `models.Category` directly (via `LoadDataFromRequest`), so `IsIncome` flows automatically for create. For update, confirm the repository's `UpdateCategory` writes all columns (see next step). If the handler uses `UpsertCategoryCommand`, map `command.IsIncome` onto the model.

- [ ] **Step 6: Persist on update** — in `api/internal/repositories/categories.go`, ensure `UpdateCategory` writes `IsIncome`. If it uses GORM `Save`/`Updates` with a struct, a `false` value is skipped by the struct form; use `Select("*")` or the map form so toggling income **off** sticks. Match the pattern used elsewhere (see `UpdateAppRole`'s map form). Example if it currently does `db.Save(category)`, that writes all fields including false — verify against the actual code and only change if it uses the zero-value-skipping struct `Updates`.

- [ ] **Step 7: Run test to verify it passes**

Run: `cd api && go test ./internal/repositories/ -run TestUpdateCategoryPersistsIsIncome -count=1`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
cd api && git add internal/models/category.go internal/commands/upsert_category_command.go internal/handlers/categories.go internal/repositories/categories.go internal/repositories/categories_test.go
git commit -m "feat(api): add Category.IsIncome flag"
```

---

## Task 2: `CategoryBudget` model + repository

**Files:**
- Create: `api/internal/models/category_budget.go`
- Create: `api/internal/commands/upsert_category_budget_command.go`
- Create: `api/internal/repositories/category_budgets.go`
- Modify: the model-registration list used by `repositories.MakeMigrations()` (grep for `&models.Category{}` in `api/internal/repositories/` to find it)
- Test: `api/internal/repositories/category_budgets_test.go`

**Interfaces:**
- Produces:
  - `models.CategoryBudget{ BaseModel; GroupId uint; CategoryId uint; Amount decimal.Decimal }`
  - `commands.UpsertCategoryBudgetCommand{ CategoryId uint; Amount decimal.Decimal }` + `Validate() structs.ValidatorError`
  - `repositories.CategoryBudgetRepository` with:
    - `NewCategoryBudgetRepository(tx *gorm.DB) CategoryBudgetRepository`
    - `GetBudgetsByGroupId(groupId uint) ([]models.CategoryBudget, error)`
    - `UpsertBudget(groupId, categoryId uint, amount decimal.Decimal) (models.CategoryBudget, error)`
    - `DeleteBudget(groupId, categoryId uint) error`
    - `DeleteBudgetsByCategoryId(tx *gorm.DB, categoryId uint) error`
    - `DeleteBudgetsByGroupId(tx *gorm.DB, groupId uint) error`

- [ ] **Step 1: Create the model** — `api/internal/models/category_budget.go`:

```go
package models

import "github.com/shopspring/decimal"

// CategoryBudget is a recurring monthly spend target for a single category
// within a single group. One target per (GroupId, CategoryId).
type CategoryBudget struct {
	BaseModel
	GroupId    uint            `gorm:"not null;uniqueIndex:idx_group_category_budget" json:"groupId"`
	CategoryId uint            `gorm:"not null;uniqueIndex:idx_group_category_budget" json:"categoryId"`
	Amount     decimal.Decimal `gorm:"not null" json:"amount"`
}
```

- [ ] **Step 2: Create the command** — `api/internal/commands/upsert_category_budget_command.go`:

```go
package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

type UpsertCategoryBudgetCommand struct {
	CategoryId uint            `json:"categoryId"`
	Amount     decimal.Decimal `json:"amount"`
}

func (command *UpsertCategoryBudgetCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &command)
}

func (command *UpsertCategoryBudgetCommand) Validate() structs.ValidatorError {
	vErr := structs.ValidatorError{}
	errorMap := make(map[string]string)

	if command.CategoryId == 0 {
		errorMap["categoryId"] = "Category is required"
	}
	if command.Amount.LessThanOrEqual(decimal.Zero) {
		errorMap["amount"] = "Amount must be greater than zero"
	}

	vErr.Errors = errorMap
	return vErr
}
```

- [ ] **Step 3: Register the model for migration** — add `&models.CategoryBudget{}` to the migration model list (the slice `MakeMigrations` iterates; grep `&models.Category{}` to find it and add the new model beside it).

- [ ] **Step 4: Write the failing repository test** — `api/internal/repositories/category_budgets_test.go`:

```go
package repositories_test

import (
	"receipt-wrangler/api/internal/repositories"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCategoryBudgetUpsertAndList(t *testing.T) {
	defer repositories.TruncateTestDb()
	repo := repositories.NewCategoryBudgetRepository(nil)

	_, err := repo.UpsertBudget(1, 5, decimal.NewFromInt(400))
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	// Upsert again updates in place, does not duplicate.
	_, err = repo.UpsertBudget(1, 5, decimal.NewFromInt(450))
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	budgets, err := repo.GetBudgetsByGroupId(1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(budgets) != 1 {
		t.Fatalf("expected 1 budget, got %d", len(budgets))
	}
	if !budgets[0].Amount.Equal(decimal.NewFromInt(450)) {
		t.Errorf("expected amount 450, got %s", budgets[0].Amount.String())
	}

	if err := repo.DeleteBudget(1, 5); err != nil {
		t.Fatalf("delete: %v", err)
	}
	budgets, _ = repo.GetBudgetsByGroupId(1)
	if len(budgets) != 0 {
		t.Errorf("expected 0 budgets after delete, got %d", len(budgets))
	}
}
```

- [ ] **Step 5: Run test to verify it fails**

Run: `cd api && go test ./internal/repositories/ -run TestCategoryBudgetUpsertAndList -count=1`
Expected: FAIL (`NewCategoryBudgetRepository` undefined).

- [ ] **Step 6: Create the repository** — `api/internal/repositories/category_budgets.go`:

```go
package repositories

import (
	"receipt-wrangler/api/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryBudgetRepository struct {
	BaseRepository
}

func NewCategoryBudgetRepository(tx *gorm.DB) CategoryBudgetRepository {
	return CategoryBudgetRepository{BaseRepository: BaseRepository{DB: GetDB(), TX: tx}}
}

func (repo CategoryBudgetRepository) GetBudgetsByGroupId(groupId uint) ([]models.CategoryBudget, error) {
	db := repo.GetDB()
	var budgets []models.CategoryBudget
	err := db.Where("group_id = ?", groupId).Find(&budgets).Error
	return budgets, err
}

func (repo CategoryBudgetRepository) UpsertBudget(groupId, categoryId uint, amount decimal.Decimal) (models.CategoryBudget, error) {
	db := repo.GetDB()
	budget := models.CategoryBudget{GroupId: groupId, CategoryId: categoryId, Amount: amount}
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "group_id"}, {Name: "category_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"amount", "updated_at"}),
	}).Create(&budget).Error
	return budget, err
}

func (repo CategoryBudgetRepository) DeleteBudget(groupId, categoryId uint) error {
	db := repo.GetDB()
	return db.Where("group_id = ? AND category_id = ?", groupId, categoryId).
		Delete(&models.CategoryBudget{}).Error
}

func (repo CategoryBudgetRepository) DeleteBudgetsByCategoryId(tx *gorm.DB, categoryId uint) error {
	db := repo.DB
	if tx != nil {
		db = tx
	}
	return db.Where("category_id = ?", categoryId).Delete(&models.CategoryBudget{}).Error
}

func (repo CategoryBudgetRepository) DeleteBudgetsByGroupId(tx *gorm.DB, groupId uint) error {
	db := repo.DB
	if tx != nil {
		db = tx
	}
	return db.Where("group_id = ?", groupId).Delete(&models.CategoryBudget{}).Error
}
```

Verify `BaseRepository` has a `GetDB()` that prefers `TX` (grep an existing repo like `categories.go`); if the base helper is named differently, match it. The `OnConflict` composite target must match the `uniqueIndex` columns (`group_id`, `category_id`).

- [ ] **Step 7: Run test to verify it passes**

Run: `cd api && go test ./internal/repositories/ -run TestCategoryBudgetUpsertAndList -count=1`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
cd api && git add internal/models/category_budget.go internal/commands/upsert_category_budget_command.go internal/repositories/category_budgets.go internal/repositories/category_budgets_test.go
# also the migration-list file you edited in Step 3
git add -p   # stage only the migration-list edit
git commit -m "feat(api): add CategoryBudget model and repository"
```

---

## Task 3: Budget permissions

**Files:**
- Modify: `api/internal/permissions/registry.go`
- Modify: `api/swagger.yml` (Permission enum)
- Test: `api/internal/permissions/registry_test.go` (already asserts sync — no new test needed, but run it)

**Interfaces:**
- Produces: `permissions.GroupBudgetsRead`, `permissions.GroupBudgetsUpdate`, `permissions.GroupBudgetsDelete`.

- [ ] **Step 1: Add the constants** — in `api/internal/permissions/registry.go`, in the group-scope constant block (near `GroupWidgetsRead`):

```go
GroupBudgetsRead   = "group.budgets.read"
GroupBudgetsUpdate = "group.budgets.update"
GroupBudgetsDelete = "group.budgets.delete"
```

- [ ] **Step 2: Add the descriptors** — in the `registry` slice, near the Dashboards block:

```go
{GroupBudgetsRead, "Read Budgets", "View category budgets and budget progress.", "Budgets", ScopeGroup},
{GroupBudgetsUpdate, "Set Budgets", "Create or change a category's monthly budget target.", "Budgets", ScopeGroup},
{GroupBudgetsDelete, "Delete Budgets", "Remove a category's budget target.", "Budgets", ScopeGroup},
```

- [ ] **Step 3: Add to the swagger Permission enum** — in `api/swagger.yml`, find the `Permission:` enum block and add the three values:

```yaml
        - "group.budgets.read"
        - "group.budgets.update"
        - "group.budgets.delete"
```

- [ ] **Step 4: Run the sync test to verify it passes**

Run: `cd api && go test ./internal/permissions/ -count=1`
Expected: PASS (registry ↔ swagger enum in sync; `LegacyGroupOwnerKeys` auto-includes the new keys — `legacy_test.go` should still pass since Owner = all group keys).

- [ ] **Step 5: Commit**

```bash
cd api && git add internal/permissions/registry.go swagger.yml
git commit -m "feat(api): add group.budgets.* permissions"
```

---

## Task 4: Budget data service

**Files:**
- Create: `api/internal/structs/budget_data.go`
- Create: `api/internal/services/budget.go`
- Test: `api/internal/services/budget_test.go`

**Interfaces:**
- Consumes: `repositories.CategoryBudgetRepository` (Task 2); the pie chart's receipt-scoping pattern (`GetPagedReceiptsByGroupId`, `PaidByListResolver`, `IntersectReceiptFilterWithGrants`, `SubstituteRestrictedCategoriesTags`).
- Produces:
  - `structs.BudgetData{ Month string; Income float64; Spent float64; Net float64; Untracked float64; Categories []structs.BudgetCategory }`
  - `structs.BudgetCategory{ CategoryId uint; Name string; Target float64; Spent float64; Over bool }`
  - `services.BudgetService` with `NewBudgetService(tx *gorm.DB) BudgetService` and `GetBudgetData(userId uint, groupId string, now time.Time) (structs.BudgetData, error)`.

**Computation rules (documented, tested):**
- **Month window:** `[first-of-month(now), first-of-next-month)` using `now`'s location. `now` is a parameter (not `time.Now()`) so tests are deterministic — the handler passes `time.Now()`.
- **Income receipt:** a receipt with ≥1 category whose `IsIncome` is true. Its full amount adds to `Income`; it is excluded from `Spent`, per-category bars, and `Untracked`.
- **Expense receipt:** every other receipt. `Spent` = Σ expense amounts (each receipt once).
- **Per-category `Spent`:** Σ expense amounts over receipts containing that (budgeted) category — fan-out double-count across categories is accepted, matching the pie chart.
- **`Untracked`:** Σ expense amounts over receipts that contain **no** budgeted category.
- **`Net`** = `Income - Spent`. **`Over`** = category spent > target.

- [ ] **Step 1: Create the structs** — `api/internal/structs/budget_data.go`:

```go
package structs

type BudgetCategory struct {
	CategoryId uint    `json:"categoryId"`
	Name       string  `json:"name"`
	Target     float64 `json:"target"`
	Spent      float64 `json:"spent"`
	Over       bool    `json:"over"`
}

type BudgetData struct {
	Month      string           `json:"month"`
	Income     float64          `json:"income"`
	Spent      float64          `json:"spent"`
	Net        float64          `json:"net"`
	Untracked  float64          `json:"untracked"`
	Categories []BudgetCategory `json:"categories"`
}
```

- [ ] **Step 2: Write the failing service test** — `api/internal/services/budget_test.go`. Seed a group, an income category, two expense categories, one budget, and receipts inside and outside the current month; assert the split. Use the package's existing seeding helpers (grep `seedReceipt` / `CreateTestGroup` / `TruncateTestDb` in `internal/services/*_test.go` for the exact helpers and mirror them). Skeleton:

```go
package services_test

import (
	"receipt-wrangler/api/internal/services"
	"testing"
	"time"
)

func TestGetBudgetData_SplitsIncomeSpendUntracked(t *testing.T) {
	defer repositories.TruncateTestDb()
	// Seed: group G, user U (member of G).
	// Categories: "Salary"(IsIncome=true), "Dining Out"(expense), "Gas"(expense).
	// Budget: Dining Out = 400 for group G.
	// Receipts this month: Salary +5000; Dining Out 300; Gas 90 (no budget).
	// Receipt last month: Dining Out 999 (must be excluded).
	now := time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC)

	svc := services.NewBudgetService(nil)
	data, err := svc.GetBudgetData(userId, groupIdStr, now)
	if err != nil {
		t.Fatalf("GetBudgetData: %v", err)
	}

	if data.Income != 5000 {
		t.Errorf("income = %v, want 5000", data.Income)
	}
	if data.Spent != 390 { // 300 + 90, last-month 999 excluded
		t.Errorf("spent = %v, want 390", data.Spent)
	}
	if data.Untracked != 90 { // Gas has no budget
		t.Errorf("untracked = %v, want 90", data.Untracked)
	}
	if data.Net != 4610 { // 5000 - 390
		t.Errorf("net = %v, want 4610", data.Net)
	}
	if len(data.Categories) != 1 || data.Categories[0].Name != "Dining Out" {
		t.Fatalf("expected 1 budgeted category Dining Out, got %+v", data.Categories)
	}
	if data.Categories[0].Spent != 300 || data.Categories[0].Target != 400 || data.Categories[0].Over {
		t.Errorf("dining bar wrong: %+v", data.Categories[0])
	}
}
```

Fill `userId`, `groupIdStr`, and seeding using the real helpers found in the package.

- [ ] **Step 3: Run test to verify it fails**

Run: `cd api && go test ./internal/services/ -run TestGetBudgetData_SplitsIncomeSpendUntracked -count=1`
Expected: FAIL (`NewBudgetService` undefined).

- [ ] **Step 4: Implement the service** — `api/internal/services/budget.go`:

```go
package services

import (
	"time"

	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

type BudgetService struct {
	BaseService
}

func NewBudgetService(tx *gorm.DB) BudgetService {
	return BudgetService{BaseService: BaseService{DB: repositories.GetDB(), TX: tx}}
}

func (service BudgetService) GetBudgetData(userId uint, groupId string, now time.Time) (structs.BudgetData, error) {
	receiptRepository := repositories.NewReceiptRepository(service.TX)
	budgetRepository := repositories.NewCategoryBudgetRepository(service.TX)
	permissionService := NewPermissionService(service.TX)

	uintGroupId, err := utils.StringToUint(groupId)
	if err != nil {
		return structs.BudgetData{}, err
	}

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	pagedRequest := commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page: -1, PageSize: -1, OrderBy: "date", SortDirection: commands.DESCENDING,
		},
	}
	if err := permissionService.IntersectReceiptFilterWithGrants(userId, uintGroupId, &pagedRequest.Filter); err != nil {
		return structs.BudgetData{}, err
	}

	receipts, _, err := receiptRepository.GetPagedReceiptsByGroupId(
		userId, groupId, pagedRequest,
		[]string{"Categories"},
		permissionService.PaidByListResolver(userId),
	)
	if err != nil {
		return structs.BudgetData{}, err
	}
	if err := permissionService.SubstituteRestrictedCategoriesTags(userId, receipts); err != nil {
		return structs.BudgetData{}, err
	}

	budgets, err := budgetRepository.GetBudgetsByGroupId(uintGroupId)
	if err != nil {
		return structs.BudgetData{}, err
	}
	budgetByCategory := make(map[uint]decimal.Decimal, len(budgets))
	for _, b := range budgets {
		budgetByCategory[b.CategoryId] = b.Amount
	}

	income := decimal.Zero
	spent := decimal.Zero
	untracked := decimal.Zero
	perCategorySpent := make(map[uint]decimal.Decimal)
	categoryNames := make(map[uint]string)

	for _, receipt := range receipts {
		if receipt.Date.Before(monthStart) || !receipt.Date.Before(monthEnd) {
			continue
		}

		isIncome := false
		for _, c := range receipt.Categories {
			if c.IsIncome {
				isIncome = true
				break
			}
		}
		if isIncome {
			income = income.Add(receipt.Amount)
			continue
		}

		spent = spent.Add(receipt.Amount)

		touchesBudget := false
		for _, c := range receipt.Categories {
			if _, ok := budgetByCategory[c.ID]; ok {
				touchesBudget = true
				categoryNames[c.ID] = c.Name
				perCategorySpent[c.ID] = perCategorySpent[c.ID].Add(receipt.Amount)
			}
		}
		if !touchesBudget {
			untracked = untracked.Add(receipt.Amount)
		}
	}

	data := structs.BudgetData{
		Month:      monthStart.Format("2006-01"),
		Income:     toFloat(income),
		Spent:      toFloat(spent),
		Net:        toFloat(income.Sub(spent)),
		Untracked:  toFloat(untracked),
		Categories: []structs.BudgetCategory{},
	}

	for _, b := range budgets {
		catSpent := perCategorySpent[b.CategoryId]
		name := categoryNames[b.CategoryId]
		if name == "" {
			name = service.categoryName(b.CategoryId)
		}
		data.Categories = append(data.Categories, structs.BudgetCategory{
			CategoryId: b.CategoryId,
			Name:       name,
			Target:     toFloat(b.Amount),
			Spent:      toFloat(catSpent),
			Over:       catSpent.GreaterThan(b.Amount),
		})
	}

	return data, nil
}

func (service BudgetService) categoryName(categoryId uint) string {
	categoryRepository := repositories.NewCategoryRepository(service.TX)
	category, err := categoryRepository.GetCategoryById(utils.UintToString(categoryId), false)
	if err != nil {
		return "Unknown"
	}
	return category.Name
}

func toFloat(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}
```

Verify against real signatures: `GetPagedReceiptsByGroupId` param order (copy from `pie_chart.go`), `models.Receipt.Date` type (`time.Time`), `GetCategoryById` signature. If `toFloat` collides with an existing helper in the package, rename to `budgetToFloat`.

- [ ] **Step 5: Run test to verify it passes**

Run: `cd api && go test ./internal/services/ -run TestGetBudgetData_SplitsIncomeSpendUntracked -count=1`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
cd api && git add internal/structs/budget_data.go internal/services/budget.go internal/services/budget_test.go
git commit -m "feat(api): add BudgetService month/income/spend computation"
```

---

## Task 5: Budget handlers, routes, swagger paths, BUDGET widget type

**Files:**
- Create: `api/internal/handlers/budgets.go`
- Create: `api/internal/routers/budget.go`
- Modify: `api/internal/routers/router.go`
- Modify: `api/swagger.yml`
- Test: `api/internal/handlers/budgets_test.go`

**Interfaces:**
- Consumes: `services.BudgetService` (Task 4), `commands.UpsertCategoryBudgetCommand` (Task 2), `permissions.GroupBudgets*` (Task 3).
- Produces: routes `POST /api/budget/{groupId}` (data), `PUT /api/budget/{groupId}` (upsert), `DELETE /api/budget/{groupId}/{categoryId}` (delete).

- [ ] **Step 1: Create the handlers** — `api/internal/handlers/budgets.go`:

```go
package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

func GetBudgetData(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")
	handler := structs.Handler{
		ErrorMessage:     "Error getting budget data",
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupBudgetsRead},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			token := structs.GetClaims(r)
			budgetService := services.NewBudgetService(nil)
			data, err := budgetService.GetBudgetData(token.UserId, groupId, time.Now())
			if err != nil {
				return http.StatusInternalServerError, err
			}
			bytes, err := utils.MarshalResponseData(data)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
			return 0, nil
		},
	}
	HandleRequest(handler)
}

func UpsertBudget(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")
	handler := structs.Handler{
		ErrorMessage:     "Error saving budget",
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupBudgetsUpdate},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			command := commands.UpsertCategoryBudgetCommand{}
			if err := command.LoadDataFromRequest(w, r); err != nil {
				return http.StatusInternalServerError, err
			}
			if vErr := command.Validate(); len(vErr.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
				return 0, nil
			}
			uintGroupId, err := utils.StringToUint(groupId)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			repo := repositories.NewCategoryBudgetRepository(nil)
			budget, err := repo.UpsertBudget(uintGroupId, command.CategoryId, command.Amount)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			bytes, err := utils.MarshalResponseData(budget)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
			return 0, nil
		},
	}
	HandleRequest(handler)
}

func DeleteBudget(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")
	categoryId := chi.URLParam(r, "categoryId")
	handler := structs.Handler{
		ErrorMessage:     "Error deleting budget",
		Writer:           w,
		Request:          r,
		GroupId:          groupId,
		GroupPermissions: []string{permissions.GroupBudgetsDelete},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			uintGroupId, err := utils.StringToUint(groupId)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			uintCategoryId, err := utils.StringToUint(categoryId)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			repo := repositories.NewCategoryBudgetRepository(nil)
			if err := repo.DeleteBudget(uintGroupId, uintCategoryId); err != nil {
				return http.StatusInternalServerError, err
			}
			w.WriteHeader(http.StatusOK)
			return 0, nil
		},
	}
	HandleRequest(handler)
}
```

Confirm helper names against `handlers/widgets.go`: `structs.GetClaims`, `utils.MarshalResponseData`, `structs.WriteValidatorErrorResponse`, `constants.ApplicationJson`.

- [ ] **Step 2: Create the router** — `api/internal/routers/budget.go`:

```go
package routers

import (
	"github.com/go-chi/chi/v5"
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"
)

func BuildBudgetRouter() *chi.Mux {
	budgetRouter := chi.NewRouter()
	budgetRouter.Use(middleware.UnifiedAuthMiddleware)
	budgetRouter.Post("/{groupId}", handlers.GetBudgetData)
	budgetRouter.Put("/{groupId}", handlers.UpsertBudget)
	budgetRouter.Delete("/{groupId}/{categoryId}", handlers.DeleteBudget)
	return budgetRouter
}
```

- [ ] **Step 3: Mount the router** — in `api/internal/routers/router.go`, beside the widget router mount (grep `BuildWidgetRouter`), add:

```go
budgetRouter := BuildBudgetRouter()
rootRouter.Mount("/api/budget", budgetRouter)
```

Match the exact mount style/prefix used by the sibling routers in that file.

- [ ] **Step 4: Write the failing handler test** — `api/internal/handlers/budgets_test.go`. Mirror `widgets_test.go` setup (auth context, seeded group/user with owner role). Assert `GetBudgetData` returns 200 with a `structs.BudgetData` body, and `UpsertBudget` with `amount<=0` returns 400. Use the package's existing test bootstrap (grep `widgets_test.go` for how it builds the request + claims).

- [ ] **Step 5: Run test to verify it fails**

Run: `cd api && go test ./internal/handlers/ -run TestBudget -count=1`
Expected: FAIL (handlers/routes not found or 404).

- [ ] **Step 6: Add swagger paths + schemas + BUDGET widget type** — in `api/swagger.yml`:
  - Add `BUDGET` to the `WidgetType` enum.
  - Add response schemas `BudgetData`, `BudgetCategory` and request schema `UpsertCategoryBudgetCommand` (mirror the Go structs; `amount`/`target`/`spent`/`income`/`net`/`untracked` as `number`, ids as `integer format: uint64`). `UpsertCategoryBudgetCommand.amount` is a `string` on the wire (decimal convention — match `Item.quantity`).
  - Add paths `/budget/{groupId}` (post → `BudgetData`, put → `UpsertCategoryBudgetCommand`) and `/budget/{groupId}/{categoryId}` (delete), each with `bearerAuth`/`apiKeyAuth` security, mirroring `/widget/pieChart/{groupId}`.

- [ ] **Step 7: Run test to verify it passes** (Go side; swagger doesn't affect Go tests)

Run: `cd api && go test ./internal/handlers/ -run TestBudget -count=1`
Expected: PASS

- [ ] **Step 8: Run the full backend suite**

Run: `cd api && CHROMIUM_BINARY_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test -count=1 ./...`
Expected: PASS

- [ ] **Step 9: Commit**

```bash
cd api && git add internal/handlers/budgets.go internal/handlers/budgets_test.go internal/routers/budget.go internal/routers/router.go swagger.yml
git commit -m "feat(api): add budget endpoints and BUDGET widget type"
```

---

## Task 6: Cascade deletes (category → budgets, group → budgets)

**Files:**
- Modify: `api/internal/repositories/categories.go` (inside `DeleteCategory`'s transaction)
- Modify: `api/internal/services/groups.go` (inside `DeleteGroup`)
- Test: `api/internal/repositories/category_budgets_test.go` (extend)

**Interfaces:**
- Consumes: `CategoryBudgetRepository.DeleteBudgetsByCategoryId`, `DeleteBudgetsByGroupId` (Task 2).

- [ ] **Step 1: Write the failing test** — extend `category_budgets_test.go`:

```go
func TestDeleteCategoryRemovesBudgets(t *testing.T) {
	defer repositories.TruncateTestDb()
	catRepo := repositories.NewCategoryRepository(nil)
	budgetRepo := repositories.NewCategoryBudgetRepository(nil)

	cat, err := catRepo.CreateCategory(models.Category{Name: "Dining Out"})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	if _, err := budgetRepo.UpsertBudget(1, cat.ID, decimal.NewFromInt(400)); err != nil {
		t.Fatalf("upsert budget: %v", err)
	}

	if err := catRepo.DeleteCategory(utils.UintToString(cat.ID)); err != nil {
		t.Fatalf("delete category: %v", err)
	}

	budgets, _ := budgetRepo.GetBudgetsByGroupId(1)
	if len(budgets) != 0 {
		t.Errorf("expected budgets removed with category, got %d", len(budgets))
	}
}
```

Match `DeleteCategory`'s real signature.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd api && go test ./internal/repositories/ -run TestDeleteCategoryRemovesBudgets -count=1`
Expected: FAIL (budget survives).

- [ ] **Step 3: Add the cascade in `DeleteCategory`** — inside the existing transaction in `api/internal/repositories/categories.go`, before deleting the category row:

```go
budgetRepository := NewCategoryBudgetRepository(tx)
if err := budgetRepository.DeleteBudgetsByCategoryId(tx, categoryId); err != nil {
	return err
}
```

Use the transaction handle (`tx`) that `DeleteCategory` already uses; `categoryId` is whatever the method holds (may need `utils.StringToUint`).

- [ ] **Step 4: Add the cascade in `DeleteGroup`** — in `api/internal/services/groups.go`, beside the member-grant cleanup, add:

```go
budgetRepository := repositories.NewCategoryBudgetRepository(tx)
if err := budgetRepository.DeleteBudgetsByGroupId(tx, groupId); err != nil {
	return err
}
```

Match the transaction variable and `groupId` type used in `DeleteGroup`.

- [ ] **Step 5: Run test to verify it passes**

Run: `cd api && go test ./internal/repositories/ -run TestDeleteCategoryRemovesBudgets -count=1`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
cd api && git add internal/repositories/categories.go internal/services/groups.go internal/repositories/category_budgets_test.go
git commit -m "feat(api): cascade budget deletes on category/group delete"
```

---

## Task 7: Regenerate clients (desktop + mobile) — enum safety

**Files:**
- Modify (generated): `desktop/src/open-api/**`, `mobile/api/**`
- Test: `mobile/api` analyze; a desktop compile.

**Interfaces:**
- Produces (desktop generated): `WidgetService.getBudgetData(groupId, ...)`, `WidgetService.upsertBudget(...)` OR a new `BudgetService` (depends on swagger `tags` — tag the budget paths `Budget` so a `BudgetService` is generated; tag the data endpoint however you prefer, but keep all three under one tag for one service). Also `BudgetData`, `BudgetCategory`, `UpsertCategoryBudgetCommand` models; `Category.isIncome`; `WidgetType.Budget`.

- [ ] **Step 1: Regenerate desktop client**

```bash
curl -fL -o /tmp/openapi-generator-cli-7.10.0.jar \
  https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.10.0/openapi-generator-cli-7.10.0.jar
cd api
java -jar /tmp/openapi-generator-cli-7.10.0.jar generate -i swagger.yml -g typescript-angular -o ../desktop/src/open-api
```

- [ ] **Step 2: Regenerate mobile client**

```bash
cd /Users/legolaswyq/git/receipt-wrangler/api
java -jar /tmp/openapi-generator-cli-7.10.0.jar generate -i swagger.yml -g dart-dio -o ../mobile/api
export PATH="$(brew --prefix dart 2>/dev/null)/libexec/bin:$PATH"   # or the Homebrew dart-sdk 3.13.3 path
cd ../mobile/api && dart pub get && dart run build_runner build --delete-conflicting-outputs
```

- [ ] **Step 3: Re-apply the two documented dart patches** — in `mobile/api/lib/src/model/user_preferences.dart` restore `..quickScanDefaultStatus = ReceiptStatus.OPEN`; in `system_settings.dart` restore `..currencyDisplay = r'` (the raw-string default). Read `mobile/CLAUDE.md` → "Known dart-dio default-value regressions" for the exact lines. **Revert any churn unrelated to the swagger change** (diff the generated output).

- [ ] **Step 4: Analyze mobile — must be 0 errors**

Run: `cd mobile/api && dart analyze`
Expected: 0 errors (warnings baseline ~73). If `WidgetType.BUDGET` const or `Category.isIncome` is missing, the generation didn't pick up the swagger edit — fix swagger and regenerate.

- [ ] **Step 5: Add the mobile enum-parse guard test** — `mobile/api/test/` (or `mobile/test/` matching where existing client tests live; grep for an existing `*_test.dart` that deserializes a model). Prove the regenerated client deserializes a dashboard/widget payload carrying `"widgetType":"BUDGET"` without throwing:

```dart
// widget_type_budget_test.dart
import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

void main() {
  test('WidgetType deserializes BUDGET', () {
    expect(WidgetType.BUDGET.name, 'BUDGET');
    expect(() => WidgetType.valueOf('BUDGET'), returnsNormally);
  });
}
```

Adjust the import/package name to the generated package (grep an existing mobile/api test import). **This is the enum-safety gate:** it proves the *new* client handles `BUDGET`. Note in the commit body that already-released mobile builds must be rebuilt before viewing a dashboard containing a budget widget (documented released-build limitation; acceptable for this self-hosted, single-operator deployment).

- [ ] **Step 6: Verify desktop still compiles**

Run: `cd desktop && npx tsc -p tsconfig.app.json --noEmit` (or `npm run build`)
Expected: no errors from the generated client.

- [ ] **Step 7: Commit**

```bash
cd /Users/legolaswyq/git/receipt-wrangler
git add desktop/src/open-api mobile/api
git commit -m "chore: regenerate API clients for budget endpoints (desktop + mobile)"
```

---

## Task 8: Desktop budget widget

**Files:**
- Create: `desktop/src/dashboard/budget/budget.component.ts`
- Create: `desktop/src/dashboard/budget/budget.component.html`
- Create: `desktop/src/dashboard/budget/budget.component.scss`
- Create: `desktop/src/dashboard/budget/budget.component.spec.ts`
- Modify: `desktop/src/dashboard/constants/widget-options.ts`
- Modify: `desktop/src/dashboard/dashboard/dashboard.component.html` + `.ts`
- Modify: `desktop/src/dashboard/dashboard-form/dashboard-form.component.html`

**Interfaces:**
- Consumes: generated `WidgetService`/`BudgetService` (`getBudgetData`, `upsertBudget`, `deleteBudget`), `BudgetData`, `BudgetCategory`, `WidgetType.Budget` (Task 7); `Widget` input like `PieChartComponent`.

- [ ] **Step 1: Add the widget option** — in `desktop/src/dashboard/constants/widget-options.ts`, add to `widgetTypeOptions`:

```ts
{
  value: WidgetType.Budget,
  displayValue: "Budget",
},
```

- [ ] **Step 2: Write the failing component spec** — `desktop/src/dashboard/budget/budget.component.spec.ts`. Mock the budget service to return a `BudgetData` and assert the component exposes the categories and computes an `over` flag. Mirror an existing dashboard component spec (grep `pie-chart.component.spec.ts` if present, else `activity.component.spec.ts`) for `TestBed` + service mock setup. Minimum assertions:

```ts
it("loads budget data and exposes categories", () => {
  // arrange: service.getBudgetData returns of({ month:'2026-09', income:5000, spent:390, net:4610, untracked:90,
  //   categories:[{categoryId:5,name:'Dining Out',target:400,spent:300,over:false}] })
  component.ngOnInit();
  expect(component.data()?.categories.length).toBe(1);
  expect(component.data()?.net).toBe(4610);
});

it("progressPercent caps at 100 for display but flags over", () => {
  expect(component.progressPercent({ target: 400, spent: 500 } as any)).toBe(100);
  expect(component.isOver({ target: 400, spent: 500 } as any)).toBe(true);
});
```

- [ ] **Step 3: Run spec to verify it fails**

Run: `cd desktop && npx jest src/dashboard/budget -t "budget" 2>&1 | tail -20`
Expected: FAIL (component missing).

- [ ] **Step 4: Implement the component** — `desktop/src/dashboard/budget/budget.component.ts`:

```ts
import { CommonModule } from "@angular/common";
import { Component, OnInit, OnChanges, SimpleChanges, input, signal } from "@angular/core";
import { take, tap } from "rxjs";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { PipesModule } from "../../pipes/pipes.module";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { BudgetCategory, BudgetData, Widget, WidgetService } from "../../open-api";

@Component({
  selector: "app-budget",
  templateUrl: "./budget.component.html",
  styleUrls: ["./budget.component.scss"],
  standalone: true,
  imports: [CommonModule, SharedUiModule, PipesModule],
  providers: [CustomCurrencyPipe],
})
export class BudgetComponent implements OnInit, OnChanges {
  public readonly widget = input.required<Widget>();
  public readonly groupId = input<number>();

  public readonly data = signal<BudgetData | undefined>(undefined);
  public readonly isLoading = signal(true);
  public readonly editingCategoryId = signal<number | null>(null);

  constructor(private widgetService: WidgetService) {}

  public ngOnInit(): void {
    this.loadData();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["groupId"] && !changes["groupId"].firstChange) {
      this.loadData();
    }
  }

  public progressPercent(category: BudgetCategory): number {
    if (!category.target) {
      return 0;
    }
    return Math.min(100, Math.round((category.spent / category.target) * 100));
  }

  public isOver(category: BudgetCategory): boolean {
    return !!category.target && category.spent > category.target;
  }

  public saveTarget(category: BudgetCategory, amount: number): void {
    const groupId = this.groupId();
    if (!groupId) {
      return;
    }
    this.widgetService
      .upsertBudget(groupId, { categoryId: category.categoryId, amount: String(amount) })
      .pipe(take(1), tap(() => { this.editingCategoryId.set(null); this.loadData(); }))
      .subscribe();
  }

  private loadData(): void {
    const groupId = this.groupId();
    if (!groupId) {
      this.isLoading.set(false);
      return;
    }
    this.isLoading.set(true);
    this.widgetService
      .getBudgetData(groupId)
      .pipe(take(1), tap((response: BudgetData) => {
        this.data.set(response);
        this.isLoading.set(false);
      }))
      .subscribe();
  }
}
```

Verify the generated method names/signatures after Task 7 (`getBudgetData(groupId)`, `upsertBudget(groupId, command)`). If the generator produced a `BudgetService` instead of methods on `WidgetService`, inject that. `amount` is sent as a string (decimal wire convention).

- [ ] **Step 5: Implement the template** — `desktop/src/dashboard/budget/budget.component.html`:

```html
<app-card>
  <ng-container header>
    <div class="d-flex justify-content-between align-items-center">
      <span>{{ widget().name ?? "Budget" }}</span>
      <span class="badge bg-secondary">{{ data()?.month }}</span>
    </div>
  </ng-container>
  <ng-container content>
    @if (isLoading()) {
      <div class="text-center p-3">Loading…</div>
    } @else if (data()) {
      <div class="d-flex justify-content-between summary mb-3">
        <div><small>Income</small><div class="income">{{ data()!.income | customCurrency }}</div></div>
        <div><small>Spent</small><div class="spent">{{ data()!.spent | customCurrency }}</div></div>
        <div><small>Net</small><div [class.negative]="data()!.net < 0">{{ data()!.net | customCurrency }}</div></div>
      </div>

      @for (category of data()!.categories; track category.categoryId) {
        <div class="budget-row mb-2">
          <div class="d-flex justify-content-between">
            <span>{{ category.name }}</span>
            <span [class.over]="isOver(category)">
              {{ category.spent | customCurrency }} / {{ category.target | customCurrency }}
            </span>
          </div>
          <div class="progress">
            <div class="progress-bar" [class.over]="isOver(category)"
                 [style.width.%]="progressPercent(category)"></div>
          </div>
        </div>
      }

      @if (data()!.untracked > 0) {
        <div class="budget-row untracked d-flex justify-content-between mt-2">
          <span>Untracked</span>
          <span>{{ data()!.untracked | customCurrency }}</span>
        </div>
      }
    }
  </ng-container>
</app-card>
```

Confirm `app-card` and the `customCurrency` pipe are the same ones `pie-chart.component.html` uses; copy its wrapper if different.

- [ ] **Step 6: Style it** — `desktop/src/dashboard/budget/budget.component.scss`:

```scss
:host {
  .summary .negative { color: var(--bs-danger, #dc3545); }
  .progress { height: 8px; }
  .progress-bar.over { background-color: var(--bs-danger, #dc3545); }
  .over { color: var(--bs-danger, #dc3545); font-weight: 600; }
  .untracked { color: var(--bs-secondary, #6c757d); }
}
```

- [ ] **Step 7: Register in the dashboard renderer** — in `desktop/src/dashboard/dashboard/dashboard.component.html`, add a case beside the others:

```html
<ng-container *ngSwitchCase="widgetType.Budget">
  <app-budget
    [widget]="widget"
    [groupId]="selectedDashboard()?.groupId"
  ></app-budget>
</ng-container>
```

And in `dashboard.component.ts`, add the standalone import `BudgetComponent` to the component's `imports` array (grep how `PieChartComponent` is imported/registered there and mirror exactly).

- [ ] **Step 8: Add the config case in the dashboard form** — in `desktop/src/dashboard/dashboard-form/dashboard-form.component.html`, add an empty `ngSwitchCase` for `WidgetType.Budget` (like `GroupSummary`, no extra config):

```html
<ng-container *ngSwitchCase="WidgetType.Budget">
</ng-container>
```

- [ ] **Step 9: Run spec to verify it passes**

Run: `cd desktop && npx jest src/dashboard/budget 2>&1 | tail -20`
Expected: PASS

- [ ] **Step 10: Build to confirm wiring**

Run: `cd desktop && npm run build`
Expected: success.

- [ ] **Step 11: Commit**

```bash
cd desktop && git add src/dashboard/budget src/dashboard/constants/widget-options.ts src/dashboard/dashboard/dashboard.component.html src/dashboard/dashboard/dashboard.component.ts src/dashboard/dashboard-form/dashboard-form.component.html
git commit -m "feat(desktop): budget dashboard widget"
```

---

## Task 9: Category income flag in the category form

**Files:**
- Modify: the category form component + template (grep `desktop/src` for the category create/edit form — likely `src/**/category-form` or a dialog).
- Test: that component's spec.

**Interfaces:**
- Consumes: `Category.isIncome` (Task 7).

- [ ] **Step 1: Locate the form** — `cd desktop && grep -rln "UpsertCategoryCommand\|categoryService" src --include=*.ts | grep -i form`. Open the form component; find where the reactive form is built.

- [ ] **Step 2: Write/extend the failing spec** — assert the form includes an `isIncome` control defaulting to the category's value:

```ts
it("includes an isIncome control", () => {
  expect(component.form.get("isIncome")).toBeTruthy();
});
```

- [ ] **Step 3: Run spec to verify it fails**

Run: `cd desktop && npx jest <category-form path> 2>&1 | tail -20`
Expected: FAIL.

- [ ] **Step 4: Add the control** — add `isIncome: new FormControl(category?.isIncome ?? false)` to the form group, include it in the submit payload, and add a checkbox to the template using the app's existing checkbox component (`app-checkbox`, as used in the quick-scan dialog) with label "This is an income category".

- [ ] **Step 5: Run spec to verify it passes**

Run: `cd desktop && npx jest <category-form path> 2>&1 | tail -20`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
cd desktop && git add <category form files>
git commit -m "feat(desktop): income flag on category form"
```

---

## Task 10: Documentation + final verification

**Files:**
- Modify: `api/CLAUDE.md`, `desktop/CLAUDE.md`, `mobile/CLAUDE.md`, root `CLAUDE.md`.

- [ ] **Step 1: Document the feature** — add sections:
  - `api/CLAUDE.md`: "Budgets" — `Category.IsIncome`, `CategoryBudget` model + unique index, `BudgetService` computation rules (month window, income = any income-flagged category, per-category fan-out, untracked = no-budgeted-category), `group.budgets.*` permissions, cascade deletes, the decimal-string wire convention for `amount`.
  - `desktop/CLAUDE.md`: the budget widget (summary + bars, inline edit, over state), the income checkbox on the category form.
  - `mobile/CLAUDE.md`: the `WidgetType.BUDGET` enum addition + the released-build caveat (rebuild before viewing a dashboard containing a budget widget).
  - Root `CLAUDE.md`: a short cross-component note under a new "Budgets" heading (three-component feature: backend owns data, desktop renders/edits, mobile regen-only + enum safety).

- [ ] **Step 2: Full backend suite**

Run: `cd api && CHROMIUM_BINARY_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test -count=1 ./...`
Expected: PASS

- [ ] **Step 3: Full desktop suite**

Run: `cd desktop && npm test -- --watch=false 2>&1 | tail -30`
Expected: all suites green.

- [ ] **Step 4: Mobile analyze**

Run: `cd mobile/api && dart analyze`
Expected: 0 errors.

- [ ] **Step 5: Client-drift check** (root CLAUDE.md convention)

```bash
cd /Users/legolaswyq/git/receipt-wrangler
swagger=$(git log -1 --format=%ct -- api/swagger.yml)
client=$(git log -1 --format=%ct -- mobile/api)
[ "$swagger" -le "$client" ] && echo "mobile/api is current" || echo "mobile/api is STALE — regenerate"
```

Expected: "mobile/api is current".

- [ ] **Step 6: Commit docs**

```bash
cd /Users/legolaswyq/git/receipt-wrangler
git add api/CLAUDE.md desktop/CLAUDE.md mobile/CLAUDE.md CLAUDE.md
git commit -m "docs: document budget dashboard feature"
```

- [ ] **Step 7: Live smoke test** — with `air` (`:8081`) and `ng serve` (`:4200`) running, log in (admin/admin), open a dashboard, add a Budget widget, set a target (e.g. Dining Out = 400), confirm the bar + summary render and the target persists (`SELECT * FROM category_budgets;`). Flag the "Income" category as income in the category form and confirm income leaves the expense bars.

---

## Self-Review Notes

- **Spec §5.1 `Category.IsIncome`** → Task 1. **§5.2 `CategoryBudget`** → Task 2. **§6.1 widget data** → Tasks 4+5. **§6.2 CRUD** → Task 5. **§6.3 permissions** → Task 3 (refinement: `read`/`update`/`delete`, dropped separate `create` since upsert covers it — noted in Global Constraints). **§6.4 widget type** → Task 5. **§7 desktop** → Tasks 8+9. **§8 mobile enum safety** → Task 7 (regen + guard test + documented released-build caveat). **§9 seeded-data caveat** → out of scope (documented). **§10 testing** → each task's tests + Task 10. **§11 rollout order** → task order.
- **Types consistent:** `CategoryBudget{GroupId,CategoryId,Amount}`, `BudgetData{Month,Income,Spent,Net,Untracked,Categories}`, `BudgetCategory{CategoryId,Name,Target,Spent,Over}`, `UpsertCategoryBudgetCommand{CategoryId,Amount}` used identically across backend, swagger, and desktop tasks.
- **Amount wire type:** decimal → string on the wire (swagger `string`, desktop sends `String(amount)`), consistent with the existing `quantity`/`unitPrice` convention. The generated desktop model will type `amount` as `string`; the component converts.
