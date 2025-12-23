package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nhosoya/rollbar-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	itemsLimit  int
	itemsStatus string
	itemsLevel  string
	itemsEnv    string
)

var itemsCmd = &cobra.Command{
	Use:   "items",
	Short: "List recent error items",
	Long:  `List recent error items from Rollbar with optional filters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		items, err := c.GetItems(itemsLimit, itemsStatus, itemsLevel, itemsEnv)
		if err != nil {
			return err
		}

		output, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(itemsCmd)

	itemsCmd.Flags().IntVarP(&itemsLimit, "limit", "n", 10, "Number of items")
	itemsCmd.Flags().StringVarP(&itemsStatus, "status", "s", "active", "Filter by status: active, resolved, muted")
	itemsCmd.Flags().StringVarP(&itemsLevel, "level", "l", "", "Filter by level: error, warning, critical")
	itemsCmd.Flags().StringVarP(&itemsEnv, "env", "e", "", "Filter by environment")
}
