package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "List all available mvnenv commands",
	Long: `Display a list of all available mvnenv commands.

This command lists all mvnenv commands available for use, one per line.
Useful for command discovery and scripting.`,
	Example: `  mvnenv commands`,
	RunE:    runCommands,
}

func init() {
	rootCmd.AddCommand(commandsCmd)
}

func runCommands(cmd *cobra.Command, args []string) error {
	// Get all commands from root
	commands := rootCmd.Commands()

	for _, c := range commands {
		if !c.Hidden {
			fmt.Println(c.Name())
		}
	}

	return nil
}
