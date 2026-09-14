package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
)

type PromptService struct {
	BaseService
}

func NewPromptService(tx *gorm.DB) PromptService {
	service := PromptService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

func (service PromptService) CreateDefaultPrompt() (models.Prompt, error) {
	db := service.GetDB()
	var defaultPromptCount int64
	db.Model(models.Prompt{}).Where("name = ?", constants.DefaultPromptName).Count(&defaultPromptCount)

	defaultPrompt := fmt.Sprintf(`
Find the receipt's name, total cost, date, and line items. Format the found data as:
{
	"name": store name,
	"amount": amount as a number,
	"date": date in ISO 18601 format in UTC with ALL time values set as 0,
	"categories": categories,
	"tags": tags,
	"receiptItems": line items
}
If a store name cannot be confidently found, use 'Default store name' as the default name.
Omit any value if not found with confidence. Assume the date is in the year @currentYear if not provided.
The amount must be a float or integer.
If the receipt represents a refund, return, or credit (money returned to the customer rather than charged), the amount and item amounts MUST be negative. Otherwise they are positive.

For every line item purchased on the receipt, add an entry to "receiptItems" formatted as:
{
	"name": the item's name in English,
	"nameZh": the item's name in Chinese, ONLY if the receipt text is in or includes Chinese, otherwise omit this field,
	"quantity": the quantity purchased, as a number,
	"unitPrice": the price for a single unit, as a number,
	"amount": the item's total price (quantity multiplied by unitPrice) as a number,
	"status": "OPEN"
}
Provide "quantity", "unitPrice", and "amount" whenever they can be confidently read from the receipt. If only two of the three can be found, omit the third rather than guessing. If none of the three can be found for a line item, omit that line item entirely rather than inventing values.
If the receipt has no readable line items, return an empty array for "receiptItems".

Please do NOT add any additional information, only valid JSON.
Please return the json in plaintext ONLY, do not ever return it in a code block or any other format.

Choose up to 2 categories from the given list based on the receipt's items and store name. If no categories fit, please return an empty array for the field and do not select any categories. When selecting categories, select only the id, like:
{
	Id: category id
}

Emphasize the relationship between the category and the receipt, and use the description of the category to fine tune the results. Do not return categories that have an empty name or do not exist.
If there are no categories to chose from, then please make categories an empty array.
Likewise, if there are not tags to choose from, then make tags an empty array.

Categories to chose from: @categories

Follow the same process as described for categories for tags.

Tags to chose from: @tags

Receipt text: @ocrText
`)

	if defaultPromptCount > 0 {
		return models.Prompt{}, errors.New("default prompt already exists")
	}

	promptRepository := repositories.NewPromptRepository(service.TX)
	command := commands.UpsertPromptCommand{
		Name:        constants.DefaultPromptName,
		Description: "Default prompt used for previous versions of Receipt Wrangler.",
		Prompt:      defaultPrompt,
	}

	return promptRepository.CreatePrompt(command)
}
