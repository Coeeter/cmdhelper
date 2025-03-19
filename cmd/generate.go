package cmd

import (
	"fmt"
	"os"

	"github.com/coeeter/cmdhelper/internal"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [prompt]",
	Short: "Generate a command based on a prompt",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		fmt.Println("Generating command for:", prompt)

		res, err := internal.GenerateCommand()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println(res)
	},
}
