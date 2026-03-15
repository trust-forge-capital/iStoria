package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
)

var sensorCmd = &cobra.Command{
	Use:   "sensor",
	Short: "Show sensor information (temperature, fan, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")

		c := collect.NewCollector()

		sensorData, err := c.Sensors()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting sensors: %v\n", err)
			os.Exit(1)
		}

		if jsonFlag {
			j, _ := json.MarshalIndent(sensorData, "", "  ")
			fmt.Println(string(j))
			return
		}

		colorReset := ""
		if noColor {
			colorReset = ""
		}

		if !sensorData.HasSensors && len(sensorData.Temperatures) == 0 && len(sensorData.Fans) == 0 {
			fmt.Println("No sensor data available on this system.")
			fmt.Println("Note: On macOS, install 'istats' for sensor support.")
			fmt.Println("On Linux, ensure 'lm-sensors' is installed.")
			return
		}

		if len(sensorData.Temperatures) > 0 {
			fmt.Printf("%s--- Temperatures ---%s\n", "", colorReset)
			for _, t := range sensorData.Temperatures {
				indicator := ""
				if t.Critical {
					indicator = " 🔥"
				}
				fmt.Printf("%-25s %7.1f %s%s\n", t.Name, t.Value, t.Unit, indicator)
			}
			fmt.Println()
		}

		if len(sensorData.Fans) > 0 {
			fmt.Printf("%s--- Fans ---%s\n", "", colorReset)
			for _, f := range sensorData.Fans {
				fmt.Printf("%-25s %7.0f %s\n", f.Name, f.Value, f.Unit)
			}
			fmt.Println()
		}

		if len(sensorData.Voltages) > 0 {
			fmt.Printf("%s--- Voltages ---%s\n", "", colorReset)
			for _, v := range sensorData.Voltages {
				fmt.Printf("%-25s %7.2f %s\n", v.Name, v.Value, v.Unit)
			}
			fmt.Println()
		}

		if len(sensorData.Power) > 0 {
			fmt.Printf("%s--- Power ---%s\n", "", colorReset)
			for _, p := range sensorData.Power {
				fmt.Printf("%-25s %7.1f %s\n", p.Name, p.Value, p.Unit)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(sensorCmd)
}
