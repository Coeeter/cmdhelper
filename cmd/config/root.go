package config

import "github.com/spf13/cobra"

var ConfigRootCommand = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

func init() {
	ConfigRootCommand.AddCommand(showConfigCmd)
	ConfigRootCommand.AddCommand(clearConfigCommand)
}
