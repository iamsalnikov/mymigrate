package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {

	test := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.Flag("email").Value)
			fmt.Println(args)
			return nil
		},
	}

	root := &cobra.Command{
		Use: "root",
	}

	test.Flags().String("path", "", "specify email")
	test.Flags().String("name", "", "specify email")

	root.AddCommand(test)
	root.Execute()

}
