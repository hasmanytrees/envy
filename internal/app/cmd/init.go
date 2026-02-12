package cmd

import (
	"envy/internal/app/shell"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init SHELL",
	Short: "Initialize an envy session",
	Long:  `Outputs sh commands to initialize an envy session and register sh hooks.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shellType := args[0]

		if !slices.Contains(shell.ValidShellTypes, shellType) {
			return errors.New(fmt.Sprintf("%s is not a supported shell type; valid values are [%s]", shellType, strings.Join(shell.ValidShellTypes, ", ")))
		}

		sessionKey := ulid.Make().String()

		shell := shell.NewShell(shellType, sessionKey)

		return shell.Init(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
