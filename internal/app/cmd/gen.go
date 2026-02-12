package cmd

import (
	"envy/internal/app/shell"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	shellType  = os.Getenv("ENVY_SHELL")
	sessionKey = os.Getenv("ENVY_SESSION_KEY")
	sh         shell.Shell
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate load and unload sh scripts",
	Long:  `Generate load and unload sh scripts for managing envy-specific environment variables.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(shellType) == 0 || len(sessionKey) == 0 {
			cmd.SilenceUsage = true
			return fmt.Errorf("could not load sh for session; run 'envy init SHELL' first")
		}

		sh = shell.NewShell(shellType, sessionKey)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// capture current env
		oldEnv := shell.NewEnv(os.Environ())

		// find load paths (sh specific)
		loadPaths, err := sh.FindLoadPaths()
		if err != nil {
			return err
		}

		// gen load script (sh specific)
		loadLines, loadFilepath := sh.GenLoadFile(loadPaths)
		err = writeLines(loadLines, loadFilepath)
		if err != nil {
			return err
		}

		// launch subshell to execute load paths and export env (sh specific)
		subshell := sh.GetSubshellCmd()
		output, err := subshell.CombinedOutput()
		if err != nil {
			return err
		}

		newEnv := shell.NewEnv(strings.Split(string(output), "\n"))

		// compare envs
		changes := oldEnv.Diff(newEnv)

		// gen unload script (sh specific)
		undoLines, undoFilepath := sh.GenUndoFile(changes)
		err = writeLines(undoLines, undoFilepath)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}

func writeLines(lines []string, name string) error {
	content := strings.Join(lines, "\n")

	return os.WriteFile(name, []byte(content), 0644)
}
