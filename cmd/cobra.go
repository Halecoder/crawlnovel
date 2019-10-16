package cmd

import (
	"crawlnovel/cmd/api"
	"crawlnovel/cmd/convert"
	"crawlnovel/cmd/download"
	"crawlnovel/cmd/search"
	"crawlnovel/cmd/task"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:               "download",
	Short:             "download API server",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `Start download API server`,
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
}

func init() {
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(download.StartCmd)
	rootCmd.AddCommand(convert.StartCmd)
	rootCmd.AddCommand(search.StartCmd)
	rootCmd.AddCommand(task.StartCmd)
}

//Execute : run commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
