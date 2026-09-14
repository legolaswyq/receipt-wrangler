package commands

import (
	"github.com/shopspring/decimal"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
)

type UpsertItemCommand struct {
	Amount          decimal.Decimal         `json:"amount"`
	Quantity        *decimal.Decimal        `json:"quantity"`
	UnitPrice       *decimal.Decimal        `json:"unitPrice"`
	NameZh          string                  `json:"nameZh"`
	ChargedToUserId *uint                   `json:"chargedToUserId"`
	IsTaxed         bool                    `json:"isTaxed"`
	Name            string                  `json:"name"`
	ReceiptId       uint                    `json:"receiptId"`
	Status          models.ItemStatus       `json:"status"`
	Categories      []UpsertCategoryCommand `json:"categories"`
	Tags            []UpsertTagCommand      `json:"tags"`
	LinkedItems     []UpsertItemCommand     `json:"linkedItems"`
}

// ResolveAmount fills in whichever of amount/quantity/unitPrice the source (AI extraction or a
// manual entry) left out, preferring values the source actually provided. It recurses into linked
// items so the same rule applies uniformly.
func (item *UpsertItemCommand) ResolveAmount() {
	if item.Quantity != nil && item.UnitPrice != nil && item.Amount.IsZero() {
		computed := item.Quantity.Mul(*item.UnitPrice)
		item.Amount = computed
	} else if item.Quantity != nil && !item.Quantity.IsZero() && item.UnitPrice == nil && !item.Amount.IsZero() {
		computed := item.Amount.Div(*item.Quantity)
		item.UnitPrice = &computed
	}

	for i := range item.LinkedItems {
		item.LinkedItems[i].ResolveAmount()
	}
}

func (item *UpsertItemCommand) Validate(receiptAmount decimal.Decimal, isCreate bool) structs.ValidatorError {
	errors := make(map[string]string)
	vErr := structs.ValidatorError{}

	if item.Amount.Abs().GreaterThan(receiptAmount.Abs()) {
		errors["amount"] = "Amount cannot be greater than receipt amount"
	}

	if len(item.Name) == 0 {
		errors["name"] = "Name is required"
	}

	if !isCreate {
		if item.ReceiptId == 0 {
			errors["receiptId"] = "Receipt Id is required"
		}
	}

	if len(item.Status) == 0 {
		errors["status"] = "Status is required"
	}

	vErr.Errors = errors
	return vErr
}
