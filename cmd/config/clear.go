package config

import (
	"github.com/coeeter/cmdhelper/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var clearConfigCommand = &cobra.Command{
	Use:   "clear",
	Short: "Clear the current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.DeleteConfig()
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}

		cmd.Println(
			color.New(color.FgGreen).Sprint("Configuration cleared"),
		)
	},
}
