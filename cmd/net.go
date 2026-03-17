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
			fmt.Printf("%-8s %-15s %10s %10s\n", "Iface", "IP", "RX", "TX")
			for _, ni := range active {
				fmt.Printf("%-8s %-15s %10s %10s\n",
					ni.Name,
					ni.IP4,
					formatBytes(ni.RxBytes),
					formatBytes(ni.TxBytes))
			}
		}
		fmt.Printf("Total: RX:%s | TX:%s\n", formatBytes(netData.TotalRx), formatBytes(netData.TotalTx))
	},
}

func init() {
	rootCmd.AddCommand(netCmd)
}
