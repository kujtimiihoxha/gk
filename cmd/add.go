package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Use to add additional transports to a service",
}

func init() {
	RootCmd.AddCommand(addCmd)
}
