package internal

import "fmt"

func GenerateCommand() (string, error) {
	context, err := createContext()
	if err != nil {
		return "", fmt.Errorf("failed to create context: %w", err)
	}

	prompt, err := context.toPrompt()
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt: %w", err)
	}

	return prompt, nil
}
