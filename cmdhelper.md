# Project Name: cmdhelper

## Overview
`cmdhelper` is an AI-powered Command Line Interface (CLI) tool built in Go. It acts as an intelligent terminal assistant that translates natural language prompts into executable shell commands using the Anthropic Claude API. Unlike a standard chatbot, it is deeply integrated into the user's local environment.

## The Core Problem
Developers frequently forget syntax for complex CLI tools (like `git`, `docker`, `find`, `tar`). Leaving the terminal to search the web or ask an AI breaks their state of flow. Furthermore, standard AI answers often fail because the AI doesn't know the user's operating system, shell, or local file structure. 

## Key Features

1. **Natural Language to Command**: Instantly translates English prompts into accurate shell commands. Example: `cmdhelper generate "find all go files modified in the last week"`.
2. **Deep Context-Awareness**: Before querying the AI, `cmdhelper` auto-detects:
   - Operating system (Linux, macOS, or Windows)
   - Currently used shell (zsh, bash, fish, PowerShell, etc.)
   - Working directory
   - Directory contents (files and folders)
   - Git branch, status, and both staged/unstaged files if in a git repo
3. **Multi-step Command Support**: For complex tasks, returns a sequence of commands (with proper order and index), guiding multi-command workflows.
4. **Alternative Commands**: Suggests alternatives when multiple approaches exist.
5. **Interactive Execution Mode**: Read explanations for each command, choose them in a numbered list, edit before execution, execute directly in the CLI, and repeat as needed—all with safety checks.
6. **Strict Structured Output via JSON Schema**: Embeds a JSON schema in the prompt and strictly enforces JSON output from the AI for safe and reliable parsing (eliminating parsing errors and hallucinated explanations).
7. **Customizable & Secure Configuration**:
   - Stores the Anthropic API key securely in `~/.cmdhelper`
   - Offers CLI commands to view or erase config
8. **Project Extensibility**: Modular design with clear separation (command parser, configuration manager, context collector, AI integration), making it easy to add features or support new shells.
9. **Rich Terminal Experience**:
   - Syntax highlighted outputs
   - Spinner during network calls
   - Informative error handling for API or parsing failures
   - Descriptive usage and help options
10. **Easy Installation**:
    - Can be installed with `go install` or built from source
    - Minimal dependencies (respects Go best practices)

## Technical Stack
- **Language:** Go (Golang)—for reliability, portability, and system integration
- **Core dependencies:**
  - `spf13/cobra` (CLI framework)
  - `liushuangls/go-anthropic` (Claude API client)
  - `fatih/color` (syntax-highlighted terminal output)
  - `briandowns/spinner` (loading spinners)

## Learning Journey & Surprising Outcomes
- **Drastic power-up of LLM accuracy** when feeding OS/shell/context/git info—AI hallucination dropped dramatically when providing real system context.
- **Challenge of enforcing pure JSON:** LLMs frequently returned unstructured text, so strict prompts and schema validation were critical to get machine-usable results every time.
- **System design modularity:** Separating context-building, AI communication, and command parsing allowed for quick iteration and makes future extensions (e.g., adding support for new APIs or shells) much simpler.

## Example Use Cases
- Quickly get the right flags for complex `find`, `tar`, or `rsync` commands
- Stage and commit changes with a descriptive message as a single prompt
- List or manipulate Docker containers without memorizing CLI options
- Teach new developers shell best-practices by showing explanations side by side with commands

## Summary
`cmdhelper` keeps developers "in the zone" by coupling local system intelligence with large language models for on-demand, personalized terminal command generation—and immediate safe execution.

---

## Terminal Walkthrough (What You Actually See)

### First Run — API Key Setup
On the very first run, before any command is generated, `cmdhelper` detects there is no config file and prompts for an API key:

```
Welcome to cmdhelper!
To get started, you'll need to enter your API key.

You can get your API key by signing up at https://console.anthropic.com.

Enter your API key: sk-ant-...

Success: Config file created
```

The key is saved to `~/.cmdhelper` and never asked for again.

---

### Scenario 1 — Simple one-off command
**User types:**
```
$ cmdhelper generate "show me all running docker containers"
```

**Terminal shows (spinner while waiting for AI):**
```
⠙ Generating command for: show me all running docker containers
```

**After AI responds:**
```
Generated commands:
  1. docker ps
     Lists all currently running Docker containers

Enter the number of the command you want to execute (or q to exit): 1

You selected: docker ps
Edit the command (default: press enter to execute as is):
Executing: docker ps

CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
3f4a2b1c9d0e   nginx     ...       ...       Up 2 hrs  ...       webserver

Command executed successfully.

Execute another command? (y/n or number of command): n
```

---

### Scenario 2 — Multi-step command with alternatives
**User types:**
```
$ cmdhelper generate "stage all changes and commit with message 'fix: update handler'"
```

**Terminal shows:**
```
Generated commands:
  1. git add .
     Stages all changes in the current directory
  2. git commit -m 'fix: update handler'
     Commits the staged changes with the provided message

Alternative commands:
  3. git commit -am 'fix: update handler'
     Stages and commits all tracked file changes in a single step

Enter the number of the command you want to execute (or q to exit): 1

You selected: git add .
Edit the command (default: press enter to execute as is):
Executing: git add .

Command executed successfully.

Execute another command? (y/n or number of command): 2

You selected: git commit -m 'fix: update handler'
Edit the command (default: press enter to execute as is):
Executing: git commit -m 'fix: update handler'

[main 3a9f1c2] fix: update handler
 1 file changed, 4 insertions(+), 1 deletion(-)

Command executed successfully.

Execute another command? (y/n or number of command): n
```

---

### Scenario 3 — Editing a command before running
**User types:**
```
$ cmdhelper generate "find all go files modified in the last week"
```

**Terminal shows:**
```
Generated commands:
  1. find . -name "*.go" -mtime -7
     Lists all Go files modified within the last 7 days

Enter the number of the command you want to execute (or q to exit): 1

You selected: find . -name "*.go" -mtime -7
Edit the command (default: press enter to execute as is): find . -name "*.go" -mtime -14

Executing: find . -name "*.go" -mtime -14

./main.go
./cmd/generate.go
./internal/claude.go

Command executed successfully.

Execute another command? (y/n or number of command): n
```

---

### Config Management Commands

**View the stored API key:**
```
$ cmdhelper config show

Current configuration:
API Key: sk-ant-xxxxxxxxxxxxxxxxxxxxxxxx
```

**Clear the saved config (e.g. to switch API keys):**
```
$ cmdhelper config clear

Configuration cleared
```

---

### Help Output
```
$ cmdhelper --help

Usage:
  cmdhelper [command]

Available Commands:
  config      Manage configuration
  generate    Generate a command based on a prompt

Flags:
  -h, --help   help for cmdhelper
```
