package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of mirage",
		Long:  "Print the version number of mirage",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stdout, "mirage v0.1.4")
		},
	}
}
