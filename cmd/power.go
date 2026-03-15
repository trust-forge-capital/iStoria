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

		colorReset := ""
		if noColor {
			colorReset = ""
		}

		fmt.Printf("%s=== Power Information ===%s\n", "", colorReset)

		if !powerData.HasBattery {
			fmt.Println("No battery detected (desktop or AC power)")
			return
		}

		fmt.Printf("Battery: %d%%\n", powerData.Percent)

		status := "Discharging"
		if powerData.Charging {
			status = "Charging"
		} else if powerData.PowerPlugged {
			status = "Plugged In"
		}
		fmt.Printf("Status: %s\n", status)

		if powerData.TimeRemaining > 0 {
			hrs := powerData.TimeRemaining / 60
			mins := powerData.TimeRemaining % 60
			if powerData.Charging {
				fmt.Printf("Time to Full: %d:%02d\n", hrs, mins)
			} else {
				fmt.Printf("Time Remaining: %d:%02d\n", hrs, mins)
			}
		}

		if powerData.Watts > 0 {
			fmt.Printf("Power Draw: %.1f W\n", powerData.Watts)
		}

		if powerData.CycleCount > 0 {
			fmt.Printf("Cycle Count: %d\n", powerData.CycleCount)
		}

		if powerData.Health != "" {
			fmt.Printf("Health: %s\n", powerData.Health)
		}
	},
}

func init() {
	rootCmd.AddCommand(powerCmd)
}
