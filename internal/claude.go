package internal

import (
	"context"
	"fmt"

	"github.com/liushuangls/go-anthropic/v2"
)

func GenerateCommand(prompt string) (string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	ctx, err := createContext()
	if err != nil {
		return "", fmt.Errorf("failed to create context: %w", err)
	}

	systemPrompt, err := ctx.toPrompt()
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt: %w", err)
	}

	client := anthropic.NewClient(cfg.ApiKey)
	res, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		System: systemPrompt,
		Model:  anthropic.ModelClaude3Dot7SonnetLatest,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(prompt),
		},
		MaxTokens: 1000,
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate command: %w", err)
	}

	return res.Content[0].GetText(), nil
}
