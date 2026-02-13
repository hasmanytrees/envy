package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envy COMMAND",
	Short: "A manager for environment variables",
	Long:  `envy is a CLI tool for managing environment variables.`,
}

func Execute() error {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd.Execute()
}
