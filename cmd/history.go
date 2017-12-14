package cmd

import (
	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View deployment history",
	RunE:  doHistory,
}

func doHistory(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	RootCmd.AddCommand(historyCmd)
}
