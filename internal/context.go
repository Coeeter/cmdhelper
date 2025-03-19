package internal

import (
	"fmt"
	"os"
	"runtime"
)

type context struct {
	Os               string
	Shell            string
	CurrentDirectory string
	Children         []string
}

func (c *context) toPrompt() (string, error) {
	const prompt = `Respond with commands which are runnable according to the user request following the json schema in context block.

FOLLOW THE SCHEMA IN CONTEXT BLOCK AND DO NOT ADD ANYTHING ELSE RETURN ONLY JSON.`

	const schema = `{
	"type": "object",
	"properties": {
		"commands": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"command": {
						"type": "string"
					},
					"isInteractive: {
						"type": "boolean"
					},
					"reason": {
						"type": "string"
					}
				},
				"required": ["command", "reason"]
			}
		}
	},
	"required": ["commands"]
}`

	result := fmt.Sprintf(`---Context---
OS: '%s'
Shell: '%s'
Current Directory: '%s'
Children: '%v'
schema: '%s'
---Context---

%s`,
		c.Os, c.Shell, c.CurrentDirectory, c.Children, schema, prompt,
	)

	return result, nil
}

func createContext() (context, error) {
	wd, err := os.Getwd()
	if err != nil {
		return context{}, fmt.Errorf("failed to get current working directory: %w", err)
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		// Windows
		shell = os.Getenv("ComSpec")
	}

	if shell == "" {
		return context{}, fmt.Errorf("failed to get shell environment variable")
	}

	goos := runtime.GOOS

	children := []string{}
	dirEntries, err := os.ReadDir(wd)
	if err != nil {
		return context{}, fmt.Errorf("failed to walk directory: %w", err)
	}

	for _, entry := range dirEntries {
		children = append(children, entry.Name())
	}

	return context{
		Os:               goos,
		Shell:            shell,
		CurrentDirectory: wd,
		Children:         children,
	}, nil
}
