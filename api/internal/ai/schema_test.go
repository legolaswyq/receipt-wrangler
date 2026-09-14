package ai

import (
	"encoding/json"
	"receipt-wrangler/api/internal/utils"
	"testing"

	"github.com/google/generative-ai-go/genai"
)

func TestReceiptExtractionSchema_MarshalsExpectedShape(t *testing.T) {
	schema := ReceiptExtractionSchema()

	bytes, err := json.Marshal(schema)
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		utils.PrintTestError(t, err, nil)
	}

	if decoded["type"] != "object" {
		utils.PrintTestError(t, decoded["type"], "object")
	}

	properties, ok := decoded["properties"].(map[string]interface{})
	if !ok {
		utils.PrintTestError(t, decoded["properties"], "a properties object")
		return
	}

	for _, field := range []string{"name", "amount", "date", "categories", "tags", "receiptItems"} {
		if _, ok := properties[field]; !ok {
			utils.PrintTestError(t, "missing top-level field "+field, "present")
		}
	}

	requiredFields, ok := decoded["required"].([]interface{})
	if !ok {
		utils.PrintTestError(t, decoded["required"], "a required array")
		return
	}
	for _, field := range []string{"name", "amount", "date", "categories", "tags", "receiptItems"} {
		if !containsString(requiredFields, field) {
			utils.PrintTestError(t, "top-level field "+field+" missing from required", "present")
		}
	}

	receiptItems, ok := properties["receiptItems"].(map[string]interface{})
	if !ok {
		utils.PrintTestError(t, properties["receiptItems"], "an object")
		return
	}
	itemSchema, ok := receiptItems["items"].(map[string]interface{})
	if !ok {
		utils.PrintTestError(t, receiptItems["items"], "an item schema")
		return
	}
	itemProperties, ok := itemSchema["properties"].(map[string]interface{})
	if !ok {
		utils.PrintTestError(t, itemSchema["properties"], "an item properties object")
		return
	}
	for _, field := range []string{"name", "nameZh", "quantity", "unitPrice", "amount", "status"} {
		if _, ok := itemProperties[field]; !ok {
			utils.PrintTestError(t, "missing item field "+field, "present")
		}
	}

	itemRequiredFields, ok := itemSchema["required"].([]interface{})
	if !ok {
		utils.PrintTestError(t, itemSchema["required"], "a required array")
		return
	}
	for _, field := range []string{"name", "nameZh", "quantity", "unitPrice", "amount", "status"} {
		if !containsString(itemRequiredFields, field) {
			utils.PrintTestError(t, "item field "+field+" missing from required", "present")
		}
	}
}

func containsString(values []interface{}, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestToGenaiSchema_ConvertsTypesAndNestedShapes(t *testing.T) {
	converted := toGenaiSchema(ReceiptExtractionSchema())

	if converted.Type != genai.TypeObject {
		utils.PrintTestError(t, converted.Type, genai.TypeObject)
	}

	receiptItems, ok := converted.Properties["receiptItems"]
	if !ok {
		utils.PrintTestError(t, "receiptItems missing from converted schema", "present")
		return
	}
	if receiptItems.Type != genai.TypeArray {
		utils.PrintTestError(t, receiptItems.Type, genai.TypeArray)
	}
	if receiptItems.Items == nil {
		utils.PrintTestError(t, nil, "a non-nil Items schema")
		return
	}
	if receiptItems.Items.Type != genai.TypeObject {
		utils.PrintTestError(t, receiptItems.Items.Type, genai.TypeObject)
	}

	quantity, ok := receiptItems.Items.Properties["quantity"]
	if !ok {
		utils.PrintTestError(t, "quantity missing from converted item schema", "present")
		return
	}
	if quantity.Type != genai.TypeNumber {
		utils.PrintTestError(t, quantity.Type, genai.TypeNumber)
	}

	status, ok := receiptItems.Items.Properties["status"]
	if !ok {
		utils.PrintTestError(t, "status missing from converted item schema", "present")
		return
	}
	if len(status.Enum) != 1 || status.Enum[0] != "OPEN" {
		utils.PrintTestError(t, status.Enum, []string{"OPEN"})
	}
}

func TestToGenaiSchema_NilInputReturnsNil(t *testing.T) {
	if toGenaiSchema(nil) != nil {
		utils.PrintTestError(t, "non-nil result", nil)
	}
}
