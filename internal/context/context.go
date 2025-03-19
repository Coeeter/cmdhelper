package context

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	_ "embed"
)

type claudeContext struct {
	Os               string
	Shell            string
	CurrentDirectory string
	Children         []string
	GitBranch        string
	GitStatus        string
	StagedFiles      []string
	UnstagedFiles    []string
}

//go:embed systemPrompt.txt
var systemPrompt string

//go:embed schema.json
var schema string

func (c *claudeContext) ToPrompt() (string, error) {
	result := fmt.Sprintf(`---Context---
OS: '%s'
Shell: '%s'
Current Directory: '%s'
Children: '%v'
Git Branch: '%s'
Git Status: '%s'
Staged Files: '%v'
Unstaged Files: '%v'
Schema: '%s'
---Context---

%s`,
		c.Os,
		c.Shell,
		c.CurrentDirectory,
		c.Children,
		c.GitBranch,
		c.GitStatus,
		c.StagedFiles,
		c.UnstagedFiles,
		schema,
		systemPrompt,
	)

	return result, nil
}

func CreateContext() (claudeContext, error) {
	wd, err := os.Getwd()
	if err != nil {
		return claudeContext{}, fmt.Errorf("failed to get current working directory: %w", err)
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		// Windows
		shell = os.Getenv("ComSpec")
	}

	if shell == "" {
		return claudeContext{}, fmt.Errorf("failed to get shell environment variable")
	}

	goos := runtime.GOOS

	children := []string{}
	dirEntries, err := os.ReadDir(wd)
	if err != nil {
		return claudeContext{}, fmt.Errorf("failed to walk directory: %w", err)
	}

	for _, entry := range dirEntries {
		children = append(children, entry.Name())
	}

	gitBranch, gitStatus, stagedFiles, unstagedFiles, err := getGitInfo(wd)
	if err != nil {
		return claudeContext{}, fmt.Errorf("failed to get git info: %w", err)
	}

	return claudeContext{
		Os:               goos,
		Shell:            shell,
		CurrentDirectory: wd,
		Children:         children,
		GitBranch:        gitBranch,
		GitStatus:        gitStatus,
		StagedFiles:      stagedFiles,
		UnstagedFiles:    unstagedFiles,
	}, nil
}

func getGitInfo(directory string) (string, string, []string, []string, error) {
	cmd := exec.Command("git", "-C", directory, "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	if err != nil {
		return "", "", nil, nil, nil // Not a git repository, return empty values
	}

	cmd = exec.Command("git", "-C", directory, "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, err := cmd.Output()
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to get git branch: %w", err)
	}

	cmd = exec.Command("git", "-C", directory, "status", "--short")
	statusOutput, err := cmd.Output()
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to get git status: %w", err)
	}

	cmd = exec.Command("git", "-C", directory, "diff", "--cached", "--name-only")
	stagedFilesOutput, err := cmd.Output()
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	cmd = exec.Command("git", "-C", directory, "diff", "--name-only")
	unstagedFilesOutput, err := cmd.Output()
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to get unstaged files: %w", err)
	}

	stagedFiles := parseGitFiles(stagedFilesOutput)
	unstagedFiles := parseGitFiles(unstagedFilesOutput)

	return string(branchOutput), string(statusOutput), stagedFiles, unstagedFiles, nil
}

func parseGitFiles(output []byte) []string {
	var files []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files
}
