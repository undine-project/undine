package cmd

import (
	"github.com/spf13/cobra"
)

var (
	WatchMode bool
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for changes in the diagram.md file",
	Run: func(cmd *cobra.Command, args []string) {
		WatchMode = true
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
}
