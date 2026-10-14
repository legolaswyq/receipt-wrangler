package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"sort"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai/jsonschema"
	"gorm.io/gorm"
)

var (
	// ErrBankStatementForbidden is a 403: the caller selected a category or payer
	// outside what they may use in the group.
	ErrBankStatementForbidden = errors.New("you do not have access to one or more selected categories or users")
	// ErrBankStatementInvalidSelection is a 400: the reviewed rows are inconsistent
	// with the stored categories or the target group.
	ErrBankStatementInvalidSelection = errors.New("invalid bank statement import")
)

type BankStatementImportService struct {
	BaseService
}

func NewBankStatementImportService(tx *gorm.DB) BankStatementImportService {
	return BankStatementImportService{
		BaseService: BaseService{
			DB: repositories.GetDB(),
			TX: tx,
		},
	}
}

// Preview parses the uploaded statement PDFs and returns their rows for review
// as one list: transfers between the user's own accounts unselected, rows that
// look already imported flagged, and a suggested category where one can be
// worked out. It creates nothing.
func (service BankStatementImportService) Preview(userId uint, command commands.BankStatementPreviewCommand) (structs.BankStatementPreview, error) {
	texts := make([]string, len(command.Files))
	for i, file := range command.Files {
		text, err := service.extractText(file)
		if err != nil {
			return structs.BankStatementPreview{}, err
		}
		texts[i] = text
	}

	return service.previewFromTexts(userId, command.GroupId, texts)
}

func (service BankStatementImportService) previewFromTexts(userId uint, groupId uint, texts []string) (structs.BankStatementPreview, error) {
	preview := structs.BankStatementPreview{
		Accounts:     []structs.BankStatementAccount{},
		Transactions: []structs.BankStatementTransaction{},
	}

	seen := map[structs.BankStatementAccount]bool{}
	for _, text := range texts {
		statement, err := ParseAnzStatementText(text)
		if err != nil {
			return structs.BankStatementPreview{}, err
		}

		// The same statement picked twice (a re-downloaded copy) is read once,
		// or every one of its rows would be offered for import twice.
		if seen[statement.BankStatementAccount] {
			continue
		}
		seen[statement.BankStatementAccount] = true

		preview.Accounts = append(preview.Accounts, statement.BankStatementAccount)
		preview.Transactions = append(preview.Transactions, statement.Transactions...)
	}

	markOwnAccountTransfers(preview.Transactions)

	// Newest first across accounts. Dates are YYYY-MM-DD, and the sort is stable
	// so rows of one account keep their statement order within a day.
	sort.SliceStable(preview.Transactions, func(i, j int) bool {
		return preview.Transactions[i].Date > preview.Transactions[j].Date
	})

	err := service.markDuplicates(groupId, preview.Transactions)
	if err != nil {
		return structs.BankStatementPreview{}, err
	}

	categories, err := service.visibleCategories(userId, groupId)
	if err != nil {
		return structs.BankStatementPreview{}, err
	}

	service.suggestCategories(groupId, categories, preview.Transactions)

	return preview, nil
}

func (service BankStatementImportService) extractText(file []byte) (string, error) {
	tempFile, err := os.CreateTemp("", "bank-statement-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(file)
	closeErr := tempFile.Close()
	if err != nil {
		return "", err
	}
	if closeErr != nil {
		return "", closeErr
	}

	return ExtractPdfText(tempFile.Name())
}

// markOwnAccountTransfers finds money moved between the uploaded accounts that
// the bank did not label a transfer (an automatic payment into a savings
// account, say). The two sides of one movement are a withdrawal from one account
// and a deposit into a different account, on the same day, for the same amount
// and naming the same counterparty: the bank prints the account holder on both.
// The counterparty is what keeps an unrelated purchase of the same amount on the
// same day from being taken for the other side. Both rows are flagged and
// unselected, since importing them would count the same money as spending and as
// income. Pairs are matched one for one, and never within a single account.
func markOwnAccountTransfers(transactions []structs.BankStatementTransaction) {
	sideKey := func(transaction structs.BankStatementTransaction, isIncome bool) string {
		return fmt.Sprintf("%s|%s|%s|%t", transaction.Date, transaction.Amount.StringFixed(2),
			strings.ToLower(strings.TrimSpace(transaction.Details)), isIncome)
	}

	unmatched := map[string][]int{}
	for i := range transactions {
		transaction := transactions[i]
		if transaction.Transfer {
			continue
		}

		oppositeKey := sideKey(transaction, !transaction.IsIncome)
		candidates := unmatched[oppositeKey]
		position := -1
		for p, candidate := range candidates {
			if transactions[candidate].AccountNumber != transaction.AccountNumber {
				position = p
				break
			}
		}

		if position < 0 {
			ownKey := sideKey(transaction, transaction.IsIncome)
			unmatched[ownKey] = append(unmatched[ownKey], i)
			continue
		}

		for _, index := range []int{i, candidates[position]} {
			transactions[index].Transfer = true
			transactions[index].Selected = false
		}
		unmatched[oppositeKey] = append(candidates[:position:position], candidates[position+1:]...)
	}
}

func statementDuplicateKey(date time.Time, name string, amount string) string {
	return date.UTC().Format(structs.BankStatementDateFormat) + "|" + strings.ToLower(strings.TrimSpace(name)) + "|" + amount
}

// markDuplicates flags rows matching an existing receipt in the group by date,
// name and amount, and unselects them. Matching consumes existing receipts one
// for one, so two genuinely identical purchases on the same day are only both
// flagged when both are already present.
func (service BankStatementImportService) markDuplicates(groupId uint, transactions []structs.BankStatementTransaction) error {
	if len(transactions) == 0 {
		return nil
	}

	earliest, latest := time.Time{}, time.Time{}
	dates := make([]time.Time, len(transactions))
	for i, transaction := range transactions {
		date, err := time.Parse(structs.BankStatementDateFormat, transaction.Date)
		if err != nil {
			return err
		}
		dates[i] = date

		if earliest.IsZero() || date.Before(earliest) {
			earliest = date
		}
		if date.After(latest) {
			latest = date
		}
	}

	receiptRepository := repositories.NewReceiptRepository(service.TX)
	existing, err := receiptRepository.GetReceiptsByGroupIdInDateRange(groupId, earliest, latest.AddDate(0, 0, 1))
	if err != nil {
		return err
	}

	remaining := make(map[string]int, len(existing))
	for _, receipt := range existing {
		remaining[statementDuplicateKey(receipt.Date, receipt.Name, receipt.Amount.StringFixed(2))]++
	}

	for i := range transactions {
		key := statementDuplicateKey(dates[i], transactions[i].Name, transactions[i].Amount.StringFixed(2))
		if remaining[key] > 0 {
			remaining[key]--
			transactions[i].Duplicate = true
			transactions[i].Selected = false
		}
	}

	return nil
}

func (service BankStatementImportService) visibleCategories(userId uint, groupId uint) ([]models.Category, error) {
	categoryRepository := repositories.NewCategoryRepository(service.TX)
	categories, err := categoryRepository.GetAllCategories("id, name, description, is_income")
	if err != nil {
		return nil, err
	}

	permissionService := NewPermissionService(service.TX)
	return permissionService.GetVisibleCategoriesForUser(userId, groupId, categories)
}

type statementMerchant struct {
	Index    int    `json:"index"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsIncome bool   `json:"isIncome"`
}

type statementCategorySuggestions struct {
	Suggestions []struct {
		Index      int  `json:"index"`
		CategoryId uint `json:"categoryId"`
	} `json:"suggestions"`
}

func statementMerchantKey(transaction structs.BankStatementTransaction) string {
	return fmt.Sprintf("%t|%s", transaction.IsIncome, strings.ToLower(transaction.Name))
}

// suggestCategories fills SuggestedCategoryId. Suggestions are a convenience for
// the review step, so every failure here (AI not configured, unreachable, or
// answering nonsense) is logged and leaves the affected rows unsuggested rather
// than failing the preview.
func (service BankStatementImportService) suggestCategories(groupId uint, categories []models.Category, transactions []structs.BankStatementTransaction) {
	if len(categories) == 0 {
		return
	}

	categoryById := make(map[uint]models.Category, len(categories))
	incomeCategoryIds := []uint{}
	for _, category := range categories {
		categoryById[category.ID] = category
		if category.IsIncome {
			incomeCategoryIds = append(incomeCategoryIds, category.ID)
		}
	}

	merchants := []statementMerchant{}
	merchantIndexByKey := map[string]int{}
	for _, transaction := range transactions {
		key := statementMerchantKey(transaction)
		if _, seen := merchantIndexByKey[key]; transaction.Duplicate || seen {
			continue
		}

		merchantIndexByKey[key] = len(merchants)
		merchants = append(merchants, statementMerchant{
			Index:    len(merchants),
			Name:     transaction.Name,
			Type:     transaction.Type,
			IsIncome: transaction.IsIncome,
		})
	}

	suggested, err := service.requestCategorySuggestions(groupId, categories, merchants)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "bank statement category suggestion failed: "+err.Error())
	}

	for i := range transactions {
		transaction := &transactions[i]
		merchantIndex, ok := merchantIndexByKey[statementMerchantKey(*transaction)]
		if !ok {
			continue
		}

		// A suggestion is only kept when it names a category the caller can use
		// and of the right kind: an expense category on a deposit would silently
		// count that income as spending.
		if categoryId, found := suggested[merchantIndex]; found {
			if category, known := categoryById[categoryId]; known && category.IsIncome == transaction.IsIncome {
				id := categoryId
				transaction.SuggestedCategoryId = &id
				continue
			}
		}

		// With a single income category there is nothing to choose between.
		if transaction.IsIncome && len(incomeCategoryIds) == 1 {
			id := incomeCategoryIds[0]
			transaction.SuggestedCategoryId = &id
		}
	}
}

func bankStatementSuggestionSchema() *jsonschema.Definition {
	return &jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"suggestions": {
				Type:        jsonschema.Array,
				Description: "One entry per merchant",
				Items: &jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"index": {
							Type:        jsonschema.Integer,
							Description: "The merchant's index from the input list",
						},
						"categoryId": {
							Type:        jsonschema.Integer,
							Description: "The id of the best matching candidate category, or 0 when none fits",
						},
					},
					Required: []string{"index", "categoryId"},
				},
			},
		},
		Required: []string{"suggestions"},
	}
}

func (service BankStatementImportService) requestCategorySuggestions(
	groupId uint,
	categories []models.Category,
	merchants []statementMerchant,
) (map[int]uint, error) {
	suggested := map[int]uint{}
	if len(merchants) == 0 {
		return suggested, nil
	}

	processingService, err := NewSystemReceiptProcessingService(service.TX, utils.UintToString(groupId))
	if err != nil {
		return suggested, err
	}
	if processingService.ReceiptProcessingSettings.ID == 0 {
		return suggested, nil
	}

	type candidate struct {
		Id          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IsIncome    bool   `json:"isIncome"`
	}
	candidates := make([]candidate, len(categories))
	for i, category := range categories {
		candidates[i] = candidate{Id: category.ID, Name: category.Name, Description: category.Description, IsIncome: category.IsIncome}
	}

	categoriesJson, err := json.Marshal(candidates)
	if err != nil {
		return suggested, err
	}
	merchantsJson, err := json.Marshal(merchants)
	if err != nil {
		return suggested, err
	}

	prompt := fmt.Sprintf(`You categorise bank statement transactions. Merchant names are often truncated.
For each merchant choose the single best category id from the candidate categories, or 0 when none fits.
A merchant with isIncome true may only use a category with isIncome true; every other merchant may only use a category with isIncome false.
Respond with JSON only, in the shape {"suggestions":[{"index":0,"categoryId":1}]}, with one entry per merchant.

Candidate categories: %s

Merchants: %s`, categoriesJson, merchantsJson)

	aiService := AiService{ReceiptProcessingSettings: processingService.ReceiptProcessingSettings}
	response, _, err := aiService.CreateChatCompletion(structs.AiChatCompletionOptions{
		Messages:       []structs.AiClientMessage{{Role: "user", Content: prompt}},
		DecryptKey:     true,
		ResponseSchema: bankStatementSuggestionSchema(),
	})
	if err != nil {
		return suggested, err
	}

	parsed := statementCategorySuggestions{}
	err = json.Unmarshal([]byte(processingService.cleanResponse(response)), &parsed)
	if err != nil {
		return suggested, err
	}

	for _, suggestion := range parsed.Suggestions {
		if suggestion.CategoryId > 0 {
			suggested[suggestion.Index] = suggestion.CategoryId
		}
	}

	return suggested, nil
}

// Import creates one resolved receipt per reviewed statement row, atomically.
func (service BankStatementImportService) Import(userId uint, command commands.BankStatementImportCommand) (structs.BankStatementImportResult, error) {
	result := structs.BankStatementImportResult{}

	groupRepository := repositories.NewGroupRepository(service.TX)
	group, err := groupRepository.GetGroupById(utils.UintToString(command.GroupId), false, false, false)
	if err != nil {
		return result, err
	}
	if group.IsAllGroup {
		return result, fmt.Errorf("%w: receipts cannot be imported into the All group", ErrBankStatementInvalidSelection)
	}

	categoryById, err := service.resolveImportCategories(userId, command)
	if err != nil {
		return result, err
	}

	now := time.Now().UTC()
	receipts := make([]models.Receipt, 0, len(command.Transactions))
	for _, transaction := range command.Transactions {
		date, err := transaction.ParsedDate()
		if err != nil {
			return result, fmt.Errorf("%w: %s", ErrBankStatementInvalidSelection, err.Error())
		}

		receipt := models.Receipt{
			Name:         strings.TrimSpace(transaction.Name),
			Amount:       transaction.Amount,
			Date:         date,
			ResolvedDate: &now,
			PaidByUserID: command.PaidByUserId,
			Status:       models.RESOLVED,
			GroupId:      command.GroupId,
			Categories:   []models.Category{},
		}
		receipt.CreatedBy = &userId

		if transaction.CategoryId != nil && *transaction.CategoryId != 0 {
			receipt.Categories = []models.Category{categoryById[*transaction.CategoryId]}
		}

		receipts = append(receipts, receipt)
	}

	err = service.GetDB().Transaction(func(tx *gorm.DB) error {
		// Categories already exist; only the join rows are written.
		err := tx.Omit("Categories.*").CreateInBatches(&receipts, 100).Error
		if err != nil {
			return err
		}

		notificationRepository := repositories.NewNotificationRepository(tx)
		body := fmt.Sprintf("%d transactions were imported from a bank statement into the group %s.",
			len(receipts), repositories.BuildParamaterisedString("groupId", command.GroupId, "name", "string"))
		return notificationRepository.SendNotificationToGroup(command.GroupId, "Bank Statement Imported", body,
			models.NOTIFICATION_TYPE_NORMAL, []interface{}{userId})
	})
	if err != nil {
		return result, err
	}

	result.CreatedCount = len(receipts)
	return result, nil
}

// resolveImportCategories authorizes the import's category and payer selection
// and returns the referenced categories by id. A deposit must use an
// income-flagged category and a withdrawal must not, because the income flag is
// the only thing that tells the two apart once they are receipts.
func (service BankStatementImportService) resolveImportCategories(userId uint, command commands.BankStatementImportCommand) (map[uint]models.Category, error) {
	categoryIds := []uint{}
	seen := map[uint]struct{}{}
	for _, transaction := range command.Transactions {
		if transaction.CategoryId == nil || *transaction.CategoryId == 0 {
			continue
		}
		if _, ok := seen[*transaction.CategoryId]; !ok {
			seen[*transaction.CategoryId] = struct{}{}
			categoryIds = append(categoryIds, *transaction.CategoryId)
		}
	}

	permissionService := NewPermissionService(service.TX)

	allowed, err := permissionService.ValidateCategoryTagSelection(userId, command.GroupId, categoryIds, nil)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrBankStatementForbidden
	}

	allowed, err = permissionService.ValidateReceiptUserSelection(userId, command.GroupId, command.PaidByUserId, nil)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrBankStatementForbidden
	}

	categoryById := map[uint]models.Category{}
	if len(categoryIds) > 0 {
		categoryRepository := repositories.NewCategoryRepository(service.TX)
		categories, err := categoryRepository.GetByIds(categoryIds)
		if err != nil {
			return nil, err
		}
		for _, category := range categories {
			categoryById[category.ID] = category
		}
	}

	for _, transaction := range command.Transactions {
		if transaction.CategoryId == nil || *transaction.CategoryId == 0 {
			continue
		}

		category, ok := categoryById[*transaction.CategoryId]
		if !ok {
			return nil, fmt.Errorf("%w: category %d does not exist", ErrBankStatementInvalidSelection, *transaction.CategoryId)
		}
		if category.IsIncome != transaction.IsIncome {
			return nil, fmt.Errorf("%w: %q needs %s category", ErrBankStatementInvalidSelection, transaction.Name,
				map[bool]string{true: "an income", false: "an expense"}[transaction.IsIncome])
		}
	}

	return categoryById, nil
}
