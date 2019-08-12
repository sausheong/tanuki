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
	Short: "Shows Tanuki version number",
	Long:  `Shows Tanuki version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Tanuki v0.1")
	},
}
