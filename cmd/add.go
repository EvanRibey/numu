package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new component to the project",
	Long:  "Creates a new component in provided feature dir in the project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many arguments called. You can only specificy one command at a time.")
		} else {
			fmt.Println("Add command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
