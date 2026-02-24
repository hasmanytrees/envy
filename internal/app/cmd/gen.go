package cmd

import (
	"context"
	"envy/internal/app/shared"
	"envy/internal/app/shell"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate load and unload sh scripts",
	Long:  `Generate load and unload sh scripts for managing envy-specific environment variables.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return genPreRun(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return genRun(cmd)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}

func genPreRun(cmd *cobra.Command) error {
	shellType := os.Getenv("ENVY_SHELL")
	sessionKey := os.Getenv("ENVY_SESSION_KEY")

	if len(shellType) == 0 || len(sessionKey) == 0 {
		cmd.SilenceErrors = true
		return fmt.Errorf("could not load sh for session; run 'envy init SHELL' first")
	}

	sh := shell.NewShell(shellType, sessionKey)

	// add the shell to the command context so we can use it during Run
	ctx := cmd.Context()
	ctx = context.WithValue(ctx, "shell", sh)
	cmd.SetContext(ctx)

	return nil
}

func genRun(cmd *cobra.Command) error {
	// retrieve the shell from the command context
	sh := cmd.Context().Value("shell").(shell.Shell)

	// capture current env
	oldEnv := shared.NewEnv(os.Environ())

	// find load paths (sh specific)
	loadPaths := sh.FindLoadPaths()

	// gen load script (sh specific)
	loadLines, loadFilepath := sh.GenLoadFile(loadPaths)
	err := writeLines(loadLines, loadFilepath)
	if err != nil {
		return err
	}

	// launch subshell to execute load paths and export env (sh specific)
	subshell := sh.GetSubshellCmd()
	output, err := subshell.CombinedOutput()
	if err != nil {
		return err
	}

	newEnv := shared.NewEnv(strings.Split(string(output), "\n"))

	// compare envs
	changes := oldEnv.Diff(newEnv)

	// gen unload script (sh specific)
	undoLines, undoFilepath := sh.GenUndoFile(changes)
	err = writeLines(undoLines, undoFilepath)
	if err != nil {
		return err
	}

	return nil
}

func writeLines(lines []string, name string) error {
	content := strings.Join(lines, "\n")

	return os.WriteFile(name, []byte(content), 0644)
}
