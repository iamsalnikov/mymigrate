package cobracmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/iamsalnikov/mymigrate"
	"github.com/spf13/cobra"
)

// DownCmd is a cobra command that applies new migrations
var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "down number of migrations",
	RunE:  DownRunE,
}

// DownRunE is a cobra run function for DownCmd command
func DownRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("please pass count of migrations to down")
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	list, err := mymigrate.Down(number)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "There is nothing to down")

		return nil
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "List of downed migrations:")
	for _, mig := range list {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), mig)
	}

	return nil
}
