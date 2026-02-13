package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export environment variables",
	Long:  `Export environment variables.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return exportRun(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func exportRun(writer io.Writer) error {
	for _, line := range os.Environ() {
		_, err := io.WriteString(writer, line+"\n")
		if err != nil {
			return err
		}
	}

	return nil
}
