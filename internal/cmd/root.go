package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "mirage",
		Short: "Universal client with fingerprint spoofing",
		Long: `Mirage v0.1.0

Universal client with fingerprint spoofing
Examples:
	mirage http https://example.com -m get -fp firefox-linux`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	httpCmd := newHttpCommand()
	rootCmd.AddCommand(httpCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
