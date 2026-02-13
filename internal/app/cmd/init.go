package cmd

import (
	"envy/internal/app/shell"
	"errors"
	"fmt"
	"io"
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
		return initRun(args[0], cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initRun(shellType string, writer io.Writer) error {
	if !slices.Contains(shell.SupportedShellTypes, shellType) {
		return errors.New(fmt.Sprintf("%s is not a supported shell type; valid values are [%s]", shellType, strings.Join(shell.SupportedShellTypes, ", ")))
	}

	sessionKey := ulid.Make().String()

	sh := shell.NewShell(shellType, sessionKey)

	return sh.Init(writer)
}
