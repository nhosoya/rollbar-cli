package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nhosoya/rollbar-cli/internal/client"
	"github.com/spf13/cobra"
)

var occurrencesLimit int

var occurrencesCmd = &cobra.Command{
	Use:     "occurrences <item_id>",
	Aliases: []string{"occ"},
	Short:   "List occurrences for an item",
	Long:    `List occurrences (instances) for a specific Rollbar item.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		occurrences, err := c.GetOccurrences(args[0], occurrencesLimit)
		if err != nil {
			return err
		}

		// Output simplified list
		type OccurrenceSummary struct {
			ID        int64  `json:"id"`
			ItemID    string `json:"item_id"`
			Timestamp string `json:"timestamp"`
		}

		summaries := make([]OccurrenceSummary, len(occurrences))
		for i, occ := range occurrences {
			summaries[i] = OccurrenceSummary{
				ID:        occ.ID,
				ItemID:    occ.ItemID,
				Timestamp: occ.Timestamp,
			}
		}

		output, err := json.MarshalIndent(summaries, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(occurrencesCmd)

	occurrencesCmd.Flags().IntVarP(&occurrencesLimit, "limit", "n", 10, "Number of occurrences")
}
