package services

import (
	"time"

	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/shopspring/decimal"
)

// BudgetService computes the current-month budget widget figures: total
// income/spend, net, per-category budget vs. actual, and untracked spend
// (expense receipts that touch no budgeted category).
type BudgetService struct {
	BaseService
}

func NewBudgetService(tx *gorm.DB) BudgetService {
	return BudgetService{BaseService: BaseService{DB: repositories.GetDB(), TX: tx}}
}

// GetBudgetData computes the budget figures for the month containing now,
// scoped to groupId and the categories/tags/receipts userId may see.
//
// Computation rules:
//   - Month window is [first-of-month(now), first-of-next-month) in now's location.
//   - A receipt with >=1 category whose IsIncome is true is an income receipt:
//     its full amount adds to Income, and it is excluded from Spent, the
//     per-category bars, and Untracked.
//   - Every other receipt is an expense receipt. Spent sums their amounts once
//     each.
//   - Per-category Spent sums expense amounts across receipts containing that
//     budgeted category (a receipt with multiple budgeted categories fans out
//     and is counted under each, matching the pie chart's convention).
//   - Untracked sums expense amounts over receipts touching no budgeted category.
//   - Net = Income - Spent. Over = category Spent > Target.
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
			Page:          -1,
			PageSize:      -1,
			OrderBy:       "date",
			SortDirection: commands.DESCENDING,
		},
	}

	err = permissionService.IntersectReceiptFilterWithGrants(userId, uintGroupId, &pagedRequest.Filter)
	if err != nil {
		return structs.BudgetData{}, err
	}

	receipts, _, err := receiptRepository.GetPagedReceiptsByGroupId(
		userId,
		groupId,
		pagedRequest,
		[]string{"Categories"},
		permissionService.PaidByListResolver(userId),
	)
	if err != nil {
		return structs.BudgetData{}, err
	}

	err = permissionService.SubstituteRestrictedCategoriesTags(userId, receipts)
	if err != nil {
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
		for _, category := range receipt.Categories {
			if category.IsIncome {
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
		for _, category := range receipt.Categories {
			if _, ok := budgetByCategory[category.ID]; ok {
				touchesBudget = true
				categoryNames[category.ID] = category.Name
				perCategorySpent[category.ID] = perCategorySpent[category.ID].Add(receipt.Amount)
			}
		}
		if !touchesBudget {
			untracked = untracked.Add(receipt.Amount)
		}
	}

	data := structs.BudgetData{
		Month:      monthStart.Format("2006-01"),
		Income:     budgetToFloat(income),
		Spent:      budgetToFloat(spent),
		Net:        budgetToFloat(income.Sub(spent)),
		Untracked:  budgetToFloat(untracked),
		Categories: []structs.BudgetCategory{},
	}

	for _, budget := range budgets {
		catSpent := perCategorySpent[budget.CategoryId]
		name := categoryNames[budget.CategoryId]
		if name == "" {
			name = service.categoryName(budget.CategoryId)
		}
		data.Categories = append(data.Categories, structs.BudgetCategory{
			CategoryId: budget.CategoryId,
			Name:       name,
			Target:     budgetToFloat(budget.Amount),
			Spent:      budgetToFloat(catSpent),
			Over:       catSpent.GreaterThan(budget.Amount),
		})
	}

	return data, nil
}

// categoryName resolves the name of a budgeted category that had no matching
// receipt this month (so it never appeared in categoryNames above).
func (service BudgetService) categoryName(categoryId uint) string {
	categoryRepository := repositories.NewCategoryRepository(service.TX)
	category, err := categoryRepository.GetCategoryById(utils.UintToString(categoryId))
	if err != nil {
		return "Unknown"
	}
	return category.Name
}

func budgetToFloat(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}
