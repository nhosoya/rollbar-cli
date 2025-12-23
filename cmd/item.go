package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nhosoya/rollbar-cli/internal/client"
	"github.com/spf13/cobra"
)

var itemCmd = &cobra.Command{
	Use:   "item <item_id>",
	Short: "Show item details",
	Long:  `Show details of a specific Rollbar item.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		item, err := c.GetItem(args[0])
		if err != nil {
			return err
		}

		output, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(itemCmd)
}
