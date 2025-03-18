package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdhelper",
	Short: "An AI-powered CLI helper for generating commands",
	Long:  `cmdhelper helps you generate and execute shell commands using AI assistance.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
