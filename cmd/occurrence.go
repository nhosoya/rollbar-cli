package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nhosoya/rollbar-cli/internal/client"
	"github.com/spf13/cobra"
)

var occurrenceFull bool

var occurrenceCmd = &cobra.Command{
	Use:     "occurrence <occurrence_id>",
	Aliases: []string{"o"},
	Short:   "Show occurrence details",
	Long:    `Show details of a specific occurrence.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		if occurrenceFull {
			// Return raw API response
			raw, err := c.GetOccurrenceRaw(args[0])
			if err != nil {
				return err
			}

			// Pretty print the JSON
			var data interface{}
			if err := json.Unmarshal(raw, &data); err != nil {
				return err
			}
			output, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(output))
			return nil
		}

		occ, err := c.GetOccurrence(args[0])
		if err != nil {
			return err
		}

		// Format output with essential fields
		result := map[string]interface{}{
			"id":        occ.ID,
			"item_id":   occ.ItemID,
			"timestamp": occ.Timestamp,
			"data":      client.FormatOccurrenceData(occ.Data),
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(occurrenceCmd)

	occurrenceCmd.Flags().BoolVarP(&occurrenceFull, "full", "f", false, "Show full occurrence data (verbose)")
}
