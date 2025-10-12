package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/version"
)

var shellCmd = &cobra.Command{
	Use:   "shell <version>",
	Short: "Set the shell-specific Maven version",
	Long: `Set the Maven version for the current shell session.

This command displays instructions for setting the MVNENV_MAVEN_VERSION
environment variable. The shell version takes precedence over both local
and global versions.

After running this command, you need to set the environment variable manually:
  PowerShell: $env:MVNENV_MAVEN_VERSION = "3.9.4"
  cmd.exe: set MVNENV_MAVEN_VERSION=3.9.4`,
	Example: `  mvnenv shell 3.9.4
  mvnenv shell 3.6.3`,
	Args: cobra.ExactArgs(1),
	RunE: runShell,
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

func runShell(cmd *cobra.Command, args []string) error {
	ver := args[0]
	mvnenvRoot := getMvnenvRoot()

	// Verify version is installed
	resolver := version.NewVersionResolver(mvnenvRoot)
	if !resolver.IsVersionInstalled(ver) {
		return fmt.Errorf("version '%s' not installed", ver)
	}

	// Output instructions for setting environment variable
	fmt.Printf("%s\n", ver)
	fmt.Println()
	fmt.Println("To set this version in your current shell session:")
	fmt.Println("  PowerShell: $env:MVNENV_MAVEN_VERSION = \"" + ver + "\"")
	fmt.Println("  cmd.exe: set MVNENV_MAVEN_VERSION=" + ver)

	return nil
}
