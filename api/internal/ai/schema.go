package ai

import (
	"github.com/google/generative-ai-go/genai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// ReceiptExtractionSchema is the canonical JSON Schema an AI provider's structured-output
// feature is given to shape a receipt extraction response. It is the single source of
// truth for the response's shape (name/amount/date/categories/tags/receiptItems, and each
// item's name/nameZh/quantity/unitPrice/amount/status) — the prompt only needs to carry
// semantic instructions (refund sign, which categories to pick, etc.), not the JSON shape
// itself, since the provider enforces this shape directly when EnforceJsonResponseFormat
// is on. Shared by all three AI clients (Ollama, OpenAI, Gemini) so they cannot drift.
func ReceiptExtractionSchema() *jsonschema.Definition {
	idOnly := jsonschema.Definition{
		Type:        jsonschema.Object,
		Description: "An existing category/tag, referenced by id only",
		Properties: map[string]jsonschema.Definition{
			"id": {
				Type:        jsonschema.Integer,
				Description: "The id of an existing category/tag chosen from the candidate list",
			},
		},
		Required: []string{"id"},
	}

	item := jsonschema.Definition{
		Type:        jsonschema.Object,
		Description: "A single line item purchased on the receipt",
		Properties: map[string]jsonschema.Definition{
			"name": {
				Type:        jsonschema.String,
				Description: "The item's name in English",
			},
			"nameZh": {
				Type:        jsonschema.String,
				Description: "The item's name in Chinese, copied verbatim from the receipt text if it shows Chinese for this item. Do not translate the English name into Chinese. Empty string if the receipt shows no Chinese for this item.",
			},
			"quantity": {
				Type:        jsonschema.Number,
				Description: "The quantity purchased. If truly unreadable, use 1.",
			},
			"unitPrice": {
				Type:        jsonschema.Number,
				Description: "The price for a single unit. If truly unreadable, use the same value as amount.",
			},
			"amount": {
				Type:        jsonschema.Number,
				Description: "The item's total price, read directly from the receipt whenever printed. If no total is printed for this item but quantity and unitPrice are both known, compute quantity multiplied by unitPrice. If none of the three can be determined, use 0.",
			},
			"status": {
				Type:        jsonschema.String,
				Enum:        []string{"OPEN"},
				Description: "Always OPEN for an extracted item",
			},
		},
		// Every field is required, not just name/status: a smaller model (confirmed against a real
		// request) will happily compute amount = quantity * unitPrice internally and still leave
		// quantity/unitPrice/nameZh out of the JSON, because an optional key is a key it doesn't
		// have to surface even when it already has the value. Required only forces presence, not a
		// non-empty/non-zero value, so "0"/""/"1" above are still honest fallbacks when a field is
		// genuinely unreadable — this is the same class of bug fixed on the top-level
		// receiptItems/categories/tags fields, just one level down.
		Required: []string{"name", "nameZh", "quantity", "unitPrice", "amount", "status"},
	}

	return &jsonschema.Definition{
		Type:        jsonschema.Object,
		Description: "Structured data extracted from a receipt",
		Properties: map[string]jsonschema.Definition{
			"name": {
				Type:        jsonschema.String,
				Description: "The store/vendor name. Use 'Default store name' if it cannot be confidently found.",
			},
			"amount": {
				Type:        jsonschema.Number,
				Description: "The receipt total. Negative for a refund, return, or credit; otherwise positive.",
			},
			"date": {
				Type:        jsonschema.String,
				Description: "The receipt date, in ISO 8601 format in UTC with all time values set to 0",
			},
			"categories": {
				Type:        jsonschema.Array,
				Items:       &idOnly,
				Description: "Up to 2 category ids chosen from the candidate list. Empty array if none fit.",
			},
			"tags": {
				Type:        jsonschema.Array,
				Items:       &idOnly,
				Description: "Tag ids chosen from the candidate list. Empty array if none fit.",
			},
			"receiptItems": {
				Type:        jsonschema.Array,
				Items:       &item,
				Description: "Line items purchased on the receipt. Empty array if none are readable.",
			},
		},
		// categories/tags/receiptItems are listed as required too, even though the receipt can
		// legitimately have none of any of them: in JSON Schema, "required" only forces a key to
		// be PRESENT, not non-empty — an empty array still satisfies it. Without this, Ollama's
		// grammar-constrained decoding allows a model to omit the key entirely (confirmed against
		// a real request: the same model reliably emitted "receiptItems": [] when it was required,
		// and omitted the key altogether when it was merely an optional property), so the model
		// would sometimes skip past populating items even when the prompt asked it to.
		Required: []string{"name", "amount", "date", "categories", "tags", "receiptItems"},
	}
}

// toGenaiSchema converts the canonical schema into Gemini's own Schema struct — the
// generative-ai-go SDK does not accept a raw JSON Schema document, only this Go type. The
// two shapes are otherwise a direct field-for-field mapping.
func toGenaiSchema(definition *jsonschema.Definition) *genai.Schema {
	if definition == nil {
		return nil
	}

	schema := &genai.Schema{
		Type:        toGenaiType(definition.Type),
		Description: definition.Description,
		Nullable:    definition.Nullable,
		Enum:        definition.Enum,
		Required:    definition.Required,
	}

	if definition.Items != nil {
		schema.Items = toGenaiSchema(definition.Items)
	}

	if definition.Properties != nil {
		schema.Properties = make(map[string]*genai.Schema, len(definition.Properties))
		for name, property := range definition.Properties {
			property := property
			schema.Properties[name] = toGenaiSchema(&property)
		}
	}

	return schema
}

func toGenaiType(dataType jsonschema.DataType) genai.Type {
	switch dataType {
	case jsonschema.Object:
		return genai.TypeObject
	case jsonschema.String:
		return genai.TypeString
	case jsonschema.Number:
		return genai.TypeNumber
	case jsonschema.Integer:
		return genai.TypeInteger
	case jsonschema.Array:
		return genai.TypeArray
	case jsonschema.Boolean:
		return genai.TypeBoolean
	default:
		return genai.TypeUnspecified
	}
}
