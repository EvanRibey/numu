package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of numu",
	Long:  `All software has versions. This is numu's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("numu component creator v0.1 -- HEAD")
	},
}
