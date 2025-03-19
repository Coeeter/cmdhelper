package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/briandowns/spinner"
	"github.com/coeeter/cmdhelper/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [prompt]",
	Short: "Generate a command based on a prompt",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Generating command for: %s", color.New(color.FgGreen).Sprint(prompt))
		s.Start()

		res, err := internal.GenerateCommand(prompt)
		s.Stop()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		color.New(color.FgGreen).Println("Generated commands:")

		var mainFlow []internal.ClaudeCommand
		var alternatives []internal.ClaudeCommand

		for _, cmd := range res.Commands {
			if cmd.CommandIndex >= 0 {
				mainFlow = append(mainFlow, cmd)
			} else {
				alternatives = append(alternatives, cmd)
			}
		}

		sort.Slice(mainFlow, func(i, j int) bool {
			return mainFlow[i].CommandIndex < mainFlow[j].CommandIndex
		})

		for i, cmd := range mainFlow {
			fmt.Printf("  %d. %s\n", i+1, color.New(color.FgCyan).Sprint(cmd.Command))
			fmt.Printf("     %s\n", color.New(color.FgYellow).Sprint(cmd.Reason))
		}

		if len(alternatives) > 0 {
			color.New(color.FgMagenta).Println("\nAlternative commands:")
			for i, cmd := range alternatives {
				count := i + 1 + len(mainFlow)
				fmt.Printf("  %d. %s\n", count, color.New(color.FgCyan).Sprint(cmd.Command))
				fmt.Printf("     %s\n", color.New(color.FgYellow).Sprint(cmd.Reason))
			}
		}
	},
}
