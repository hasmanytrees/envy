package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export environment variables",
	Long:  `Export environment variables.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		exportRun(cmd.Println)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func exportRun(println func(i ...interface{})) {
	for _, line := range os.Environ() {
		println(line)
	}
}
