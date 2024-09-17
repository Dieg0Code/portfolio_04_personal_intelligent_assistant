package provider

import (
	openai "github.com/sashabaranov/go-openai"
)

func NewOperAiClient(openaiAPIKey string) *openai.Client {
	config := openai.DefaultConfig(openaiAPIKey)
	config.BaseURL = "https://models.inference.ai.azure.com"
	client := openai.NewClientWithConfig(config)

	return client
}
