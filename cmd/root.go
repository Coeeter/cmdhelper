package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "cmdhelper",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Hello World!")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
