package cobracmd

import (
	"fmt"

	"github.com/iamsalnikov/mymigrate"
	"github.com/spf13/cobra"
)

// HistoryCmd is a cobra command that prints list of applied migrations
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "shows list of applied migrations",
	RunE:  HistoryRunE,
}

// HistoryRunE is a cobra run function for HistoryCmd command
func HistoryRunE(cmd *cobra.Command, args []string) error {
	list, err := mymigrate.History()
	if err != nil {
		return err
	}

	if len(list) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "History is empty")

		return nil
	}

	for _, mig := range list {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), mig)
	}

	return nil
}
