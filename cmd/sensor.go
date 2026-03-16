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

		_ = noColor // reserved for future color support

		if !sensorData.HasSensors && len(sensorData.Temperatures) == 0 && len(sensorData.Fans) == 0 {
			fmt.Println("No sensor data. Install istats (macOS) or lm-sensors (Linux).")
			return
		}

		if len(sensorData.Temperatures) > 0 {
			fmt.Print("Temp: ")
			for i, t := range sensorData.Temperatures {
				if i > 0 {
					fmt.Print(" ")
				}
				indicator := ""
				if t.Critical {
					indicator = "🔥"
				}
				fmt.Printf("%s=%.1f%s%s", t.Name, t.Value, t.Unit, indicator)
				if i >= 5 {
					if len(sensorData.Temperatures) > 6 {
						fmt.Print(" ...")
					}
					break
				}
			}
			fmt.Println()
		}

		if len(sensorData.Fans) > 0 {
			fmt.Print("Fan: ")
			for i, f := range sensorData.Fans {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Printf("%s=%.0f%s", f.Name, f.Value, f.Unit)
				if i >= 5 {
					if len(sensorData.Fans) > 6 {
						fmt.Print(" ...")
					}
					break
				}
			}
			fmt.Println()
		}

		if len(sensorData.Voltages) > 0 {
			fmt.Print("Volt: ")
			for i, v := range sensorData.Voltages {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Printf("%s=%.2f%s", v.Name, v.Value, v.Unit)
				if i >= 5 {
					if len(sensorData.Voltages) > 6 {
						fmt.Print(" ...")
					}
					break
				}
			}
			fmt.Println()
		}

		if len(sensorData.Power) > 0 {
			fmt.Print("Power: ")
			for i, p := range sensorData.Power {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Printf("%s=%.1f%s", p.Name, p.Value, p.Unit)
				if i >= 5 {
					if len(sensorData.Power) > 6 {
						fmt.Print(" ...")
					}
					break
				}
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(sensorCmd)
}
