package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
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

		colorReset := ""
		if noColor {
			colorReset = ""
		}

		fmt.Printf("%s=== Network Interfaces ===%s\n", "", colorReset)
		fmt.Printf("%-10s %-18s %-18s %12s %12s\n",
			"Interface", "IPv4", "IPv6", "RX", "TX")
		fmt.Println()

		for _, ni := range netData.Interfaces {
			// Skip loopback for cleaner output
			if ni.Name == "lo0" || ni.Name == "lo" {
				continue
			}
			fmt.Printf("%-10s %-18s %-18s %12s %12s\n",
				ni.Name,
				ni.IP4,
				truncateIP(ni.IP6),
				formatBytes(ni.RxBytes),
				formatBytes(ni.TxBytes))
		}

		fmt.Println()
		fmt.Printf("Total RX: %s | Total TX: %s\n",
			formatBytes(netData.TotalRx),
			formatBytes(netData.TotalTx))
	},
}

func truncateIP(ip string) string {
	// Handle empty or very long IPs
	if ip == "" {
		return "-"
	}
	if len(ip) > 18 {
		// Remove IPv6 prefix for display
		if strings.HasPrefix(ip, "fe80::") {
			ip = ip[:min(len(ip), 15)]
		} else if len(ip) > 15 {
			return ip[:12] + "..."
		}
	}
	return ip
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	rootCmd.AddCommand(netCmd)
}
