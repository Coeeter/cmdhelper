package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
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

		var selectedCmdIndex int = -1

		for {
			if selectedCmdIndex == -1 {
				selectedCmdIndex, err = getSelectedCommand(res.Commands)
				if err != nil {
					fmt.Println("Invalid command number")
					continue
				}
			}

			var selectedCmd internal.ClaudeCommand
			if selectedCmdIndex <= len(mainFlow) {
				selectedCmd = mainFlow[selectedCmdIndex-1]
			} else {
				selectedCmd = alternatives[selectedCmdIndex-len(mainFlow)-1]
			}

			fmt.Printf("\nYou selected: %s\n", color.New(color.FgCyan).Sprint(selectedCmd.Command))
			fmt.Print("Edit the command (default: press enter to execute as is): ")

			reader := bufio.NewReader(os.Stdin)
			editedCmd, _ := reader.ReadString('\n')
			editedCmd = strings.TrimSpace(editedCmd)

			if editedCmd == "" {
				editedCmd = selectedCmd.Command
			}

			fmt.Printf("Executing: %s\n", color.New(color.FgCyan).Sprint(editedCmd))

			err = executeCommand(editedCmd)
			selectedCmdIndex = -1
			if err != nil {
				fmt.Println("Error executing command:", err)
			} else {
				color.New(color.FgGreen).Println("Command executed successfully.")
			}

			fmt.Print("\nExecute another command? (y/n or number of command): ")
			var again string
			fmt.Scanln(&again)
			if again == "n" {
				os.Exit(0)
			}

			if again == "y" {
				continue
			}

			var nextIndex int
			_, err = fmt.Sscanf(again, "%d", &nextIndex)
			if err != nil || nextIndex < 1 || nextIndex > len(res.Commands) {
				fmt.Println("Invalid command number")
				continue
			}

			selectedCmdIndex = nextIndex
		}
	},
}

func getSelectedCommand(commands []internal.ClaudeCommand) (int, error) {
	fmt.Print("\nEnter the number of the command you want to execute (or q to exit): ")

	var input string
	fmt.Scanln(&input)

	if input == "q" {
		os.Exit(0)
	}

	var selectedIndex int
	_, err := fmt.Sscanf(input, "%d", &selectedIndex)
	if err != nil || selectedIndex < 1 || selectedIndex > len(commands) {
		return -1, fmt.Errorf("invalid command number")
	}
	return selectedIndex, nil
}

func executeCommand(cmd string) error {
	var shell string
	if runtime.GOOS == "windows" {
		shell = "cmd"
	} else {
		shell = os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
	}

	var cmdExec *exec.Cmd
	if shell == "cmd" {
		cmdExec = exec.Command(shell, "/C", cmd)
	} else {
		cmdExec = exec.Command(shell, "-c", cmd)
	}

	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Stdin = os.Stdin

	return cmdExec.Run()
}
