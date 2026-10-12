package structs

import "github.com/sashabaranov/go-openai/jsonschema"

type AiClientMessage struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type AiChatCompletionOptions struct {
	// Messages to send to the AI model
	Messages []AiClientMessage `json:"messages"`

	// Determines whether to decrypt the key
	DecryptKey bool `json:"decryptKey"`

	// ResponseSchema overrides the structured-output schema sent to the provider
	// when EnforceJsonResponseFormat is on. Nil means the receipt extraction schema.
	ResponseSchema *jsonschema.Definition `json:"-"`
}
