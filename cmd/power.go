package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
)

var powerCmd = &cobra.Command{
	Use:   "power",
	Short: "Show power/battery information",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")

		c := collect.NewCollector()

		powerData, err := c.Power()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting power: %v\n", err)
			os.Exit(1)
		}

		if jsonFlag {
			j, _ := json.MarshalIndent(powerData, "", "  ")
			fmt.Println(string(j))
			return
		}

		_ = noColor // reserved for future color support

		if !powerData.HasBattery {
			fmt.Println("Power: no battery")
			return
		}

		status := "Discharging"
		if powerData.Charging {
			status = "Charging"
		} else if powerData.PowerPlugged {
			status = "Plugged"
		}

		line := fmt.Sprintf("Bat: %d%% | %s", powerData.Percent, status)
		if powerData.TimeRemaining > 0 {
			hrs := powerData.TimeRemaining / 60
			mins := powerData.TimeRemaining % 60
			if powerData.Charging {
				line += fmt.Sprintf(" | Full:%d:%02d", hrs, mins)
			} else {
				line += fmt.Sprintf(" | Left:%d:%02d", hrs, mins)
			}
		}
		if powerData.Watts > 0 {
			line += fmt.Sprintf(" | %.1fW", powerData.Watts)
		}
		if powerData.CycleCount > 0 {
			line += fmt.Sprintf(" | Cycles:%d", powerData.CycleCount)
		}
		if powerData.Health != "" {
			line += fmt.Sprintf(" | Health:%s", powerData.Health)
		}
		fmt.Println(line)
	},
}

func init() {
	rootCmd.AddCommand(powerCmd)
}
