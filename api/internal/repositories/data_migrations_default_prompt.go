package repositories

import (
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"

	"gorm.io/gorm"
)

// Name of the one-time migration that rewrites an unmodified seeded "Default
// Prompt" row to request line items. Prompts.CreateDefaultPrompt only ever
// seeds this row once, on first boot, so an already-running install's stored
// prompt text is frozen at whatever CreateDefaultPrompt looked like the day
// it was seeded — a later change to that Go source text never reaches an
// existing database on its own.
const addLineItemsToDefaultPromptMigration = "add-line-items-to-default-prompt"

// preLineItemsDefaultPromptText is a frozen copy of the exact text
// CreateDefaultPrompt seeded before line items were added to it. Used only as
// a guard: the migration rewrites the "Default Prompt" row only when its
// stored text still matches this byte-for-byte, so an administrator who
// customized the seeded prompt is never overwritten.
const preLineItemsDefaultPromptText = `
Find the receipt's name, total cost, and date. Format the found data as:
{
	"name": store name,
	"amount": amount as a number,
	"date": date in ISO 18601 format in UTC with ALL time values set as 0,
	"categories": categories,
	"tags": tags
}
If a store name cannot be confidently found, use 'Default store name' as the default name.
Omit any value if not found with confidence. Assume the date is in the year @currentYear if not provided.
The amount must be a float or integer.
If the receipt represents a refund, return, or credit (money returned to the customer rather than charged), the amount and item amounts MUST be negative. Otherwise they are positive.

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
`

// lineItemsDefaultPromptText mirrors the current CreateDefaultPrompt text in
// services/prompts.go at the time this migration was written. It is
// deliberately a frozen copy rather than a shared reference to that
// function's string: a migration is a point-in-time rewrite, and a later
// change to the live default prompt should not silently change what an
// already-applied migration wrote.
const lineItemsDefaultPromptText = `
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
`

// addLineItemsToDefaultPrompt rewrites the seeded "Default Prompt" row's text
// to request line items, but only when it is still byte-for-byte identical to
// the pre-line-items seed text — an administrator's own edits are left
// untouched. A prompt not found (deleted, or a fresh install that seeds the
// new text directly) is a no-op, not an error.
func addLineItemsToDefaultPrompt(tx *gorm.DB) error {
	result := tx.Model(&models.Prompt{}).
		Where("name = ? AND prompt = ?", constants.DefaultPromptName, preLineItemsDefaultPromptText).
		Update("prompt", lineItemsDefaultPromptText)

	return result.Error
}
