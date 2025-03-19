package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	ctx "github.com/coeeter/cmdhelper/internal/context"
	"github.com/liushuangls/go-anthropic/v2"
)

type ClaudeCommand struct {
	CommandIndex int    `json:"commandIndex"`
	Command      string `json:"command"`
	Reason       string `json:"reason"`
}

type ClaudeResponse struct {
	Commands []ClaudeCommand `json:"commands"`
}

var jsonRegex = regexp.MustCompile(`\{[\s\S]*\}`)

func GenerateCommand(prompt string) (ClaudeResponse, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return ClaudeResponse{}, fmt.Errorf("failed to load config: %w", err)
	}

	ctx, err := ctx.CreateContext()
	if err != nil {
		return ClaudeResponse{}, fmt.Errorf("failed to create context: %w", err)
	}

	systemPrompt, err := ctx.ToPrompt()
	if err != nil {
		return ClaudeResponse{}, fmt.Errorf("failed to generate prompt: %w", err)
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
		return ClaudeResponse{}, fmt.Errorf("failed to generate command: %w", err)
	}

	responseText := res.Content[0].GetText()

	if responseText == "" {
		return ClaudeResponse{}, fmt.Errorf("response is empty")
	}

	matches := jsonRegex.FindString(responseText)
	if len(matches) == 0 {
		return ClaudeResponse{}, fmt.Errorf("failed to find json in response")
	}

	var response ClaudeResponse

	err = json.Unmarshal([]byte(matches), &response)
	if err != nil {
		return ClaudeResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}
