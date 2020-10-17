package cobracmd

import (
	"fmt"

	"github.com/iamsalnikov/mymigrate"
	"github.com/spf13/cobra"
)

// ApplyCmd is a cobra command that applies new migrations
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply new migrations",
	RunE:  ApplyRunE,
}

// ApplyRunE is a cobra run function for ApplyCmd command
func ApplyRunE(cmd *cobra.Command, args []string) error {
	list, err := mymigrate.Apply()
	if err != nil {
		return err
	}

	if len(list) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "There are no new migrations")

		return nil
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "List of applied migrations:")
	for _, mig := range list {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), mig)
	}

	return nil
}
