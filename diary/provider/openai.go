package provider

import (
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func NewOperAiClient() *openai.Client {
	apiKey := os.Getenv("OPENAI_API_KEY")
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://models.inference.ai.azure.com"
	client := openai.NewClientWithConfig(config)

	return client
}
