package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
	"os"
)

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Show network interface information",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")

		// Check for live mode
		liveConfig := GetLiveConfig(cmd)
		if liveConfig.Enabled {
			if jsonFlag {
				RunNetLiveJSON(liveConfig)
			} else {
				RunNetLive(liveConfig)
			}
			return
		}

		c := collect.NewCollector()

		netData, err := c.Net()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting network: %v\n", err)
			os.Exit(1)
		}

		if jsonFlag {
			j, _ := json.MarshalIndent(netData, "", "  ")
			fmt.Println(string(j))
			return
		}

		_ = noColor // reserved for future color support
		noColorFlag := noColor

		// Filter to show only interfaces with activity or valid IPs
		var active []collect.NetInfo
		for _, ni := range netData.Interfaces {
			if ni.Name == "lo0" || ni.Name == "lo" {
				continue
			}
			if ni.RxBytes > 0 || ni.TxBytes > 0 || ni.IP4 != "" {
				active = append(active, ni)
			}
		}

		if len(active) > 0 {
			fmt.Printf("%-6s %-15s %12s %12s\n", 
				Colorize("Iface", ColorBold, noColorFlag),
				Colorize("IP", ColorGray, noColorFlag),
				Colorize("RX", ColorGray, noColorFlag),
				Colorize("TX", ColorGray, noColorFlag))
			for _, ni := range active {
				fmt.Printf("%-6s %-15s %12s %12s\n",
					Colorize(ni.Name, ColorCyan, noColorFlag),
					ni.IP4,
					Colorize(formatBytes(ni.RxBytes), ColorGreen, noColorFlag),
					Colorize(formatBytes(ni.TxBytes), ColorYellow, noColorFlag))
			}
		}
		fmt.Printf("%s %s | %s\n",
			Colorize("Total:", ColorBold, noColorFlag),
			Colorize("RX:"+formatBytes(netData.TotalRx), ColorGreen, noColorFlag),
			Colorize("TX:"+formatBytes(netData.TotalTx), ColorYellow, noColorFlag))
	},
}

func init() {
	rootCmd.AddCommand(netCmd)
}
