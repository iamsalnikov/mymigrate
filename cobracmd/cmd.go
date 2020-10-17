package cobracmd

import "github.com/spf13/cobra"

// MigrateCmd is a cobra command to work with migrations
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "work with migrations",
}

func init() {
	MigrateCmd.AddCommand(CreateCmd, HistoryCmd)
}
