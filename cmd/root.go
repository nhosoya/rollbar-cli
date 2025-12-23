package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rollbar",
	Short: "A CLI tool to query Rollbar errors",
	Long:  `rollbar-cli is a lightweight CLI tool to query Rollbar errors.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
