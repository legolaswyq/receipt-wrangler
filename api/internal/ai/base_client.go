package ai

import (
	config "receipt-wrangler/api/internal/env"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/sashabaranov/go-openai/jsonschema"
)

type BaseClient struct {
	Options                   structs.AiChatCompletionOptions
	ReceiptProcessingSettings models.ReceiptProcessingSettings
}

func (baseClient BaseClient) getKey(decryptKey bool) (string, error) {
	if decryptKey && len(baseClient.ReceiptProcessingSettings.Key) > 0 {
		return baseClient.decryptKey()
	}

	return baseClient.ReceiptProcessingSettings.Key, nil
}

func (baseClient BaseClient) decryptKey() (string, error) {
	return utils.DecryptB64EncodedData(config.GetEncryptionKey(), baseClient.ReceiptProcessingSettings.Key)
}

// responseSchema is the structured-output schema to enforce: the caller's
// override when one is supplied, otherwise the receipt extraction schema.
func (baseClient BaseClient) responseSchema() *jsonschema.Definition {
	if baseClient.Options.ResponseSchema != nil {
		return baseClient.Options.ResponseSchema
	}

	return ReceiptExtractionSchema()
}
