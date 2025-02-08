package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print application version",
	Long:  "print application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
