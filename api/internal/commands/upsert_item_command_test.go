package commands

import (
	"github.com/shopspring/decimal"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestUpsertItemCommand_Validate_ValidInputs(t *testing.T) {
	receiptAmount := decimal.NewFromFloat(100.00)

	tests := map[string]struct {
		command  UpsertItemCommand
		isCreate bool
	}{
		"valid create": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(50.00),
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate: true,
		},
		"valid update": {
			command: UpsertItemCommand{
				Amount:    decimal.NewFromFloat(50.00),
				Name:      "Test Item",
				ReceiptId: 1,
				Status:    models.ITEM_OPEN,
			},
			isCreate: false,
		},
		"amount equals receipt amount": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(100.00),
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate: true,
		},
		"zero amount": {
			command: UpsertItemCommand{
				Amount: decimal.Zero,
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate: true,
		},
		"negative amount within receipt magnitude": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(-50.00),
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate: true,
		},
		"negative amount equal to receipt magnitude": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(-100.00),
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate(receiptAmount, test.isCreate)

			if len(vErr.Errors) > 0 {
				utils.PrintTestError(t, len(vErr.Errors), 0)
			}
		})
	}
}

func TestUpsertItemCommand_Validate_InvalidInputs(t *testing.T) {
	receiptAmount := decimal.NewFromFloat(100.00)

	tests := map[string]struct {
		command       UpsertItemCommand
		isCreate      bool
		expectedError string
	}{
		"amount exceeds receipt amount": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(150.00),
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate:      true,
			expectedError: "amount",
		},
		"negative amount magnitude exceeds receipt": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(-150.00),
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate:      true,
			expectedError: "amount",
		},
		"positive item on negative receipt exceeds magnitude": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(150.00),
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate:      true,
			expectedError: "amount",
		},
		"missing name": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(50.00),
				Status: models.ITEM_OPEN,
			},
			isCreate:      true,
			expectedError: "name",
		},
		"missing receipt id on update": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(50.00),
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			},
			isCreate:      false,
			expectedError: "receiptId",
		},
		"missing status": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(50.00),
				Name:   "Test Item",
			},
			isCreate:      true,
			expectedError: "status",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			vErr := test.command.Validate(receiptAmount, test.isCreate)

			if len(vErr.Errors) == 0 {
				utils.PrintTestError(t, len(vErr.Errors), "greater than 0")
			}

			if _, exists := vErr.Errors[test.expectedError]; !exists {
				utils.PrintTestError(t, "error should exist for field", test.expectedError)
			}
		})
	}
}

func TestUpsertItemCommand_Validate_MultipleErrors(t *testing.T) {
	command := UpsertItemCommand{}

	vErr := command.Validate(decimal.NewFromFloat(100.00), false)

	if len(vErr.Errors) < 3 {
		utils.PrintTestError(t, len(vErr.Errors), "at least 3")
	}

	if _, exists := vErr.Errors["name"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "name")
	}

	if _, exists := vErr.Errors["status"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "status")
	}

	if _, exists := vErr.Errors["receiptId"]; !exists {
		utils.PrintTestError(t, "error should exist for field", "receiptId")
	}
}

func TestUpsertItemCommand_ResolveAmount(t *testing.T) {
	qty := decimal.NewFromInt(3)
	unitPrice := decimal.NewFromFloat(2.50)
	amount := decimal.NewFromFloat(7.50)

	tests := map[string]struct {
		command        UpsertItemCommand
		expectedAmount decimal.Decimal
		expectUnitSet  bool
	}{
		"amount provided is left untouched even with quantity and unit price": {
			command: UpsertItemCommand{
				Amount:    decimal.NewFromFloat(9.99),
				Quantity:  &qty,
				UnitPrice: &unitPrice,
			},
			expectedAmount: decimal.NewFromFloat(9.99),
		},
		"amount computed from quantity and unit price when amount is zero": {
			command: UpsertItemCommand{
				Amount:    decimal.Zero,
				Quantity:  &qty,
				UnitPrice: &unitPrice,
			},
			expectedAmount: amount,
		},
		"unit price derived from amount and quantity when unit price missing": {
			command: UpsertItemCommand{
				Amount:   amount,
				Quantity: &qty,
			},
			expectedAmount: amount,
			expectUnitSet:  true,
		},
		"no quantity or unit price leaves amount untouched": {
			command: UpsertItemCommand{
				Amount: decimal.NewFromFloat(5.00),
			},
			expectedAmount: decimal.NewFromFloat(5.00),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			test.command.ResolveAmount()

			if !test.command.Amount.Equal(test.expectedAmount) {
				utils.PrintTestError(t, test.command.Amount.String(), test.expectedAmount.String())
			}

			if test.expectUnitSet && test.command.UnitPrice == nil {
				utils.PrintTestError(t, "unit price should have been derived", "nil")
			}
		})
	}
}

func TestUpsertItemCommand_ResolveAmount_RecursesIntoLinkedItems(t *testing.T) {
	qty := decimal.NewFromInt(2)
	unitPrice := decimal.NewFromFloat(4.00)

	command := UpsertItemCommand{
		Amount: decimal.NewFromFloat(8.00),
		Name:   "Parent",
		LinkedItems: []UpsertItemCommand{
			{
				Amount:    decimal.Zero,
				Quantity:  &qty,
				UnitPrice: &unitPrice,
				Name:      "Child",
			},
		},
	}

	command.ResolveAmount()

	if !command.LinkedItems[0].Amount.Equal(decimal.NewFromFloat(8.00)) {
		utils.PrintTestError(t, command.LinkedItems[0].Amount.String(), "8")
	}
}

func TestUpsertItemCommand_Validate_NegativeReceiptMagnitudeChecks(t *testing.T) {
	negativeReceiptAmount := decimal.NewFromFloat(-100.00)

	tests := map[string]struct {
		amount       decimal.Decimal
		expectPasses bool
	}{
		"negative item within negative receipt magnitude": {
			amount:       decimal.NewFromFloat(-50.00),
			expectPasses: true,
		},
		"negative item equal to negative receipt magnitude": {
			amount:       decimal.NewFromFloat(-100.00),
			expectPasses: true,
		},
		"negative item exceeds negative receipt magnitude": {
			amount:       decimal.NewFromFloat(-150.00),
			expectPasses: false,
		},
		"positive item exceeds negative receipt magnitude": {
			amount:       decimal.NewFromFloat(150.00),
			expectPasses: false,
		},
		"positive item within negative receipt magnitude": {
			amount:       decimal.NewFromFloat(50.00),
			expectPasses: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			cmd := UpsertItemCommand{
				Amount: test.amount,
				Name:   "Test Item",
				Status: models.ITEM_OPEN,
			}

			vErr := cmd.Validate(negativeReceiptAmount, true)
			_, amountErrExists := vErr.Errors["amount"]

			if test.expectPasses && amountErrExists {
				utils.PrintTestError(t, "amount should pass magnitude check", test.amount.String())
			}
			if !test.expectPasses && !amountErrExists {
				utils.PrintTestError(t, "amount should fail magnitude check", test.amount.String())
			}
		})
	}
}
