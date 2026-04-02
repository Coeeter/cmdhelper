# cmdhelper

An AI-powered CLI tool that generates shell commands from natural language prompts using Claude AI.

## Overview

`cmdhelper` takes a plain-English description of what you want to do and uses the Anthropic Claude API to generate the appropriate shell command(s). It is context-aware — it reads your current OS, shell, working directory, directory contents, and git status to suggest commands that are relevant to your environment.

## Features

- **Natural language to command** — describe a task in plain English and get runnable shell commands back
- **Context-aware** — passes your OS, shell, current directory, file listing, and git information to Claude for better suggestions
- **Multi-step commands** — for tasks that require multiple commands, they are returned in the correct execution order
- **Alternative commands** — where applicable, alternative approaches are also suggested
- **Interactive execution** — select, optionally edit, and execute a command directly from the CLI
- **Configurable** — stores your Anthropic API key in `~/.cmdhelper`

## Requirements

- Go 1.24.1 or later
- An [Anthropic API key](https://console.anthropic.com)

## Installation

```bash
go install github.com/coeeter/cmdhelper@latest
```

Or build from source:

```bash
git clone https://github.com/coeeter/cmdhelper.git
cd cmdhelper
go build -o cmdhelper .
```

## Configuration

On first run, `cmdhelper` will prompt you to enter your Anthropic API key. The key is saved to `~/.cmdhelper`.

You can manage your configuration at any time using the `config` subcommand:

```bash
# Show the current configuration
cmdhelper config show

# Clear the saved configuration
cmdhelper config clear
```

## Usage

```bash
cmdhelper generate "<prompt>"
```

### Examples

```bash
# Find all Go files modified in the last 7 days
cmdhelper generate "find all go files modified in the last week"

# Stage and commit all changes
cmdhelper generate "stage all changes and commit with message 'fix: update handler'"

# List running Docker containers
cmdhelper generate "show me all running docker containers"
```

### Interactive workflow

1. `cmdhelper` sends your prompt along with environment context to Claude.
2. The generated commands are displayed, grouped into a **main flow** (ordered steps) and **alternative commands**.
3. Enter the number of the command you want to run.
4. Optionally edit the command before executing — press **Enter** to run it as-is.
5. After execution, choose to run another command, re-run, or exit.

```
Generated commands:
  1. find . -name "*.go" -mtime -7
     Lists all Go files modified within the last 7 days

Enter the number of the command you want to execute (or q to exit): 1

You selected: find . -name "*.go" -mtime -7
Edit the command (default: press enter to execute as is): 
Executing: find . -name "*.go" -mtime -7
...
Command executed successfully.

Execute another command? (y/n or number of command): n
```

## Project Structure

```
cmdhelper/
├── main.go                      # Entry point
├── cmd/
│   ├── root.go                  # Root cobra command
│   ├── generate.go              # generate subcommand
│   └── config/
│       ├── root.go              # config subcommand group
│       ├── show.go              # config show
│       └── clear.go             # config clear
└── internal/
    ├── claude.go                # Anthropic API integration
    ├── config.go                # Config load/save/delete
    └── context/
        ├── context.go           # Environment context builder
        ├── schema.json          # JSON schema for Claude responses
        └── systemPrompt.txt     # System prompt sent to Claude
```

## Dependencies

| Package | Purpose |
|---|---|
| [spf13/cobra](https://github.com/spf13/cobra) | CLI framework |
| [liushuangls/go-anthropic](https://github.com/liushuangls/go-anthropic) | Anthropic Claude API client |
| [fatih/color](https://github.com/fatih/color) | Coloured terminal output |
| [briandowns/spinner](https://github.com/briandowns/spinner) | Loading spinner |

## License

This project is open source. See the repository for license details.
