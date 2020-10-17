package cobracmd

import (
	"fmt"

	"github.com/iamsalnikov/mymigrate"
	"github.com/spf13/cobra"
)

// NewListCmd is a cobra command that prints list of applied migrations
var NewListCmd = &cobra.Command{
	Use:   "new-list",
	Short: "shows list of new migrations",
	RunE:  NewListRunE,
}

// NewListRunE is a cobra run function for NewListCmd command
func NewListRunE(cmd *cobra.Command, args []string) error {
	list, err := mymigrate.NewNames()
	if err != nil {
		return err
	}

	if len(list) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "There are no new migrations")

		return nil
	}

	for _, mig := range list {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), mig)
	}

	return nil
}
