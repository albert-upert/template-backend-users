package cmd

import (
	"os"

	"github.com/albert-upert/template-backend-users/version"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "project",
	Short: "Project is a simple project management tool",
}

func init() {
	// set default time to asia/jakarta
	_ = os.Setenv("TZ", "Asia/Jakarta")

	// added version information
	root.Version = version.Version
	root.SetVersionTemplate(version.String())

	// add available command
	root.AddCommand(serve)
}

// Execute will initiate all registered commands
func Execute() error {
	return root.Execute()
}
