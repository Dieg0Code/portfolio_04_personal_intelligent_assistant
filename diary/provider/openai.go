package provider

import (
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var openaiAPIKey = os.Getenv("OPENAI_API_KEY")

func NewOperAiClient() *openai.Client {
	config := openai.DefaultConfig(openaiAPIKey)
	config.BaseURL = "https://models.inference.ai.azure.com"
	client := openai.NewClientWithConfig(config)

	return client
}
