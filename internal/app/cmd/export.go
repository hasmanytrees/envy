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
	Run: func(cmd *cobra.Command, args []string) {
		exportRun(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func exportRun(writer io.Writer) {
	for _, line := range os.Environ() {
		io.WriteString(writer, line)
	}
}
