package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/version"
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
