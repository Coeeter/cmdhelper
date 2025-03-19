package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type Config struct {
	ApiKey string `json:"apiKey"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configFilePath := filepath.Join(home, ".cmdhelper")
	return configFilePath, nil
}

func LoadConfig() (Config, error) {
	var cfg Config

	configFilename, err := getConfigFilePath()
	if err != nil {
		return cfg, fmt.Errorf("failed to get config file path: %w", err)
	}

	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		return createConfig()
	}

	file, err := os.Open(configFilename)
	if err != nil {
		return cfg, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to decode config file: %w", err)
	}

	return cfg, nil
}

func createConfig() (Config, error) {
	var cfg Config

	configFilename, err := getConfigFilePath()
	if err != nil {
		return cfg, fmt.Errorf("failed to get config file path: %w", err)
	}

	fmt.Println(
		color.New(color.FgGreen).Sprint("Welcome to cmdhelper!"),
	)
	fmt.Println("To get started, you'll need to enter your API key.")
	fmt.Println()
	fmt.Println("You can get your API key by signing up at https://console.anthropic.com.")
	fmt.Println()

	for cfg.ApiKey == "" {
		fmt.Print("Enter your API key: ")
		_, err := fmt.Scanln(&cfg.ApiKey)
		if err != nil {
			fmt.Println(
				color.New(color.FgRed).Sprint("Error:"),
				"Failed to read input",
			)
			fmt.Println()
		}
	}

	file, err := os.Create(configFilename)
	if err != nil {
		return cfg, fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(cfg); err != nil {
		return cfg, fmt.Errorf("failed to encode config file: %w", err)
	}

	fmt.Println()
	fmt.Println(
		color.New(color.FgGreen).Sprint("Success:"),
		"Config file created",
	)
	fmt.Println()

	return cfg, nil
}

func DeleteConfig() error {
	configFilename, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}

	if err := os.Remove(configFilename); err != nil {
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	return nil
}
