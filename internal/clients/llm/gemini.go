package llm

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/genai"
)

type GeminiClient struct{}

func NewGeminiClient() *GeminiClient {
	return &GeminiClient{}
}

func (c *GeminiClient) GenerateFeedback(ctx context.Context, prompt string, apiKey string) (string, error) {
	timeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client, err := genai.NewClient(timeCtx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return "", fmt.Errorf("Error creating GenAI client: %w", err)
	}

	result, err := client.Models.GenerateContent(
		timeCtx,
		"gemini-2.5-flash-lite",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("Error generating content: %w", err)
	}

	return result.Text(), nil
}
