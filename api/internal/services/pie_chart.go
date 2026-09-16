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

type PieChartService struct {
	BaseService
}

func NewPieChartService(tx *gorm.DB) PieChartService {
	service := PieChartService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

func (service PieChartService) GetPieChartData(
	userId uint,
	groupId string,
	command commands.PieChartDataCommand,
) (structs.PieChartData, error) {
	receiptRepository := repositories.NewReceiptRepository(service.TX)
	permissionService := NewPermissionService(service.TX)

	uintGroupId, err := utils.StringToUint(groupId)
	if err != nil {
		return structs.PieChartData{}, err
	}

	pagedRequest := commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          -1,
			PageSize:      -1,
			OrderBy:       "date",
			SortDirection: commands.DESCENDING,
		},
		Filter: command.Filter,
	}

	// Restrict any category/tag filter to what the caller may see (anti-probing).
	err = permissionService.IntersectReceiptFilterWithGrants(userId, uintGroupId, &pagedRequest.Filter)
	if err != nil {
		return structs.PieChartData{}, err
	}

	receipts, _, err := receiptRepository.GetPagedReceiptsByGroupId(
		userId,
		groupId,
		pagedRequest,
		[]string{"Categories", "Tags"},
		permissionService.PaidByListResolver(userId),
	)
	if err != nil {
		return structs.PieChartData{}, err
	}

	// Replace categories/tags the caller cannot see with a (Restricted) marker so
	// their spend aggregates into its own slice rather than collapsing into
	// Uncategorized/Untagged (which would hide it among genuinely uncategorized
	// receipts). This matches the reporting engine's treatment of hidden values.
	err = permissionService.SubstituteRestrictedCategoriesTags(userId, receipts)
	if err != nil {
		return structs.PieChartData{}, err
	}

	// The pie chart is a spending chart, so drop receipts in income-flagged
	// categories (mirrors BudgetService.GetBudgetData). Income categories default
	// to IsIncome=false, so installs that never flag income are unaffected.
	receipts = excludeIncomeReceipts(receipts)

	// Narrow to the widget's selected date range (all-time when both are empty).
	receipts = filterReceiptsByDateRange(receipts, command.StartDate, command.EndDate)

	pieChartData := structs.PieChartData{
		Data: []structs.PieChartDataPoint{},
	}

	switch command.ChartGrouping {
	case models.CHART_GROUPING_CATEGORIES:
		pieChartData.Data = service.groupByCategories(receipts)
	case models.CHART_GROUPING_TAGS:
		pieChartData.Data = service.groupByTags(receipts)
	case models.CHART_GROUPING_PAIDBY:
		pieChartData.Data, err = service.groupByPaidBy(receipts)
		if err != nil {
			return structs.PieChartData{}, err
		}
	}

	return pieChartData, nil
}

// filterReceiptsByDateRange keeps receipts whose Date falls in [start, end).
// Each bound is an optional RFC3339 timestamp; an empty or unparseable bound is
// treated as open, so both empty means no filtering (all-time).
func filterReceiptsByDateRange(receipts []models.Receipt, startDate string, endDate string) []models.Receipt {
	var start, end time.Time
	hasStart, hasEnd := false, false
	if parsed, err := time.Parse(time.RFC3339, startDate); err == nil {
		start, hasStart = parsed, true
	}
	if parsed, err := time.Parse(time.RFC3339, endDate); err == nil {
		end, hasEnd = parsed, true
	}
	if !hasStart && !hasEnd {
		return receipts
	}

	result := make([]models.Receipt, 0, len(receipts))
	for _, receipt := range receipts {
		if hasStart && receipt.Date.Before(start) {
			continue
		}
		if hasEnd && !receipt.Date.Before(end) {
			continue
		}
		result = append(result, receipt)
	}
	return result
}

// excludeIncomeReceipts drops any receipt that carries at least one income-flagged
// category, so a spending pie chart never counts income (which is modelled as a
// positive amount in an IsIncome category). Same rule as BudgetService.
func excludeIncomeReceipts(receipts []models.Receipt) []models.Receipt {
	result := make([]models.Receipt, 0, len(receipts))
	for _, receipt := range receipts {
		isIncome := false
		for _, category := range receipt.Categories {
			if category.IsIncome {
				isIncome = true
				break
			}
		}
		if !isIncome {
			result = append(result, receipt)
		}
	}
	return result
}

func (service PieChartService) groupByCategories(receipts []models.Receipt) []structs.PieChartDataPoint {
	categoryAmounts := make(map[string]decimal.Decimal)
	categoryColors := make(map[string]string)

	for _, receipt := range receipts {
		if len(receipt.Categories) == 0 {
			if _, exists := categoryAmounts["Uncategorized"]; !exists {
				categoryAmounts["Uncategorized"] = decimal.NewFromInt(0)
			}
			categoryAmounts["Uncategorized"] = categoryAmounts["Uncategorized"].Add(receipt.Amount)
		} else {
			for _, category := range receipt.Categories {
				if _, exists := categoryAmounts[category.Name]; !exists {
					categoryAmounts[category.Name] = decimal.NewFromInt(0)
				}
				categoryAmounts[category.Name] = categoryAmounts[category.Name].Add(receipt.Amount)
				categoryColors[category.Name] = category.Color
			}
		}
	}

	return service.convertToDataPointsWithColors(categoryAmounts, categoryColors)
}

func (service PieChartService) groupByTags(receipts []models.Receipt) []structs.PieChartDataPoint {
	tagAmounts := make(map[string]decimal.Decimal)

	for _, receipt := range receipts {
		if len(receipt.Tags) == 0 {
			if _, exists := tagAmounts["Untagged"]; !exists {
				tagAmounts["Untagged"] = decimal.NewFromInt(0)
			}
			tagAmounts["Untagged"] = tagAmounts["Untagged"].Add(receipt.Amount)
		} else {
			for _, tag := range receipt.Tags {
				if _, exists := tagAmounts[tag.Name]; !exists {
					tagAmounts[tag.Name] = decimal.NewFromInt(0)
				}
				tagAmounts[tag.Name] = tagAmounts[tag.Name].Add(receipt.Amount)
			}
		}
	}

	return service.convertToDataPoints(tagAmounts)
}

func (service PieChartService) groupByPaidBy(receipts []models.Receipt) ([]structs.PieChartDataPoint, error) {
	userAmounts := make(map[uint]decimal.Decimal)
	userNames := make(map[uint]string)
	userRepository := repositories.NewUserRepository(service.TX)

	for _, receipt := range receipts {
		if _, exists := userAmounts[receipt.PaidByUserID]; !exists {
			userAmounts[receipt.PaidByUserID] = decimal.NewFromInt(0)
		}
		userAmounts[receipt.PaidByUserID] = userAmounts[receipt.PaidByUserID].Add(receipt.Amount)

		if _, exists := userNames[receipt.PaidByUserID]; !exists {
			user, err := userRepository.GetUserById(receipt.PaidByUserID)
			if err != nil {
				userNames[receipt.PaidByUserID] = "Unknown User"
			} else {
				if len(user.DisplayName) > 0 {
					userNames[receipt.PaidByUserID] = user.DisplayName
				} else {
					userNames[receipt.PaidByUserID] = user.Username
				}
			}
		}
	}

	dataPoints := make([]structs.PieChartDataPoint, 0, len(userAmounts))
	for userId, amount := range userAmounts {
		floatVal, _ := amount.Float64()
		dataPoints = append(dataPoints, structs.PieChartDataPoint{
			Label: userNames[userId],
			Value: floatVal,
		})
	}

	return dataPoints, nil
}

func (service PieChartService) convertToDataPoints(amounts map[string]decimal.Decimal) []structs.PieChartDataPoint {
	dataPoints := make([]structs.PieChartDataPoint, 0, len(amounts))
	for name, amount := range amounts {
		floatVal, _ := amount.Float64()
		dataPoints = append(dataPoints, structs.PieChartDataPoint{
			Label: name,
			Value: floatVal,
		})
	}
	return dataPoints
}

// convertToDataPointsWithColors is convertToDataPoints plus a per-label color
// (used by the category grouping so each slice keeps its category's stored
// color). A label with no color (e.g. Uncategorized) carries an empty string,
// which the client renders from its fallback palette.
func (service PieChartService) convertToDataPointsWithColors(amounts map[string]decimal.Decimal, colors map[string]string) []structs.PieChartDataPoint {
	dataPoints := make([]structs.PieChartDataPoint, 0, len(amounts))
	for name, amount := range amounts {
		floatVal, _ := amount.Float64()
		dataPoints = append(dataPoints, structs.PieChartDataPoint{
			Label: name,
			Value: floatVal,
			Color: colors[name],
		})
	}
	return dataPoints
}
