package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdhelper",
	Short: "An AI-powered CLI helper for generating commands",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	initStyles()

	rootCmd.AddCommand(generateCmd)
}

func initStyles() {
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgHiCyan, color.Bold).SprintFunc())

	usageTemplate := rootCmd.UsageTemplate()

	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Global Flags:`, `{{StyleHeading "Global Flags:"}}`,
	).Replace(usageTemplate)

	re := regexp.MustCompile(`(?m)^Flags:\s*$`)

	usageTemplate = re.ReplaceAllLiteralString(usageTemplate, `{{StyleHeading "Flags:"}}`)

	rootCmd.SetUsageTemplate(usageTemplate)
}
