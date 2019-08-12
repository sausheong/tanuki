package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tanuki",
	Short: "Tanuki is a polyglot web framework",
	Long:  `Tanuki allows web developers to create web applications and services in multiple programming languages`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute executes the Tanuki command line application
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		danger("Cannot execute Tanuki root command", err)
		os.Exit(1)
	}
}
