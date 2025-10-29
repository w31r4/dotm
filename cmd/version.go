package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of dotm",
	Long:  `All software has versions. This is dotm's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dotm version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
