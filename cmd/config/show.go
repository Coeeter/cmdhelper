package config

import (
	"github.com/coeeter/cmdhelper/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.LoadConfig()
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}

		cmd.Println(
			color.New(color.FgGreen).Sprint("Current configuration:"),
			"\nAPI Key:", config.ApiKey,
		)
	},
}
