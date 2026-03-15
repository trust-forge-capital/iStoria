package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
)

var memCmd = &cobra.Command{
	Use:   "mem",
	Short: "Show memory information",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		liveConfig := GetLiveConfig(cmd)
		c := collect.NewCollector()

		runMem := func() {
			if liveConfig.Enabled && !liveConfig.NoClear {
				ClearScreen()
			}

			memInfo, err := c.Mem()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting memory: %v\n", err)
				os.Exit(1)
			}

			if jsonFlag {
				j, _ := json.MarshalIndent(memInfo, "", "  ")
				fmt.Println(string(j))
				return
			}

			colorReset := ""
			if noColor {
				colorReset = ""
			}

			fmt.Printf("%s=== Memory Information ===%s\n", "", colorReset)
			fmt.Printf("Total:     %s\n", formatBytes(memInfo.Total))
			fmt.Printf("Used:      %s (%s)\n", formatBytes(memInfo.Used), fmt.Sprintf("%.1f%%", memInfo.UsedPercent))
			fmt.Printf("Available: %s\n", formatBytes(memInfo.Available))
			fmt.Printf("Free:      %s\n", formatBytes(memInfo.Free))

			if memInfo.SwapTotal > 0 {
				fmt.Println()
				fmt.Printf("%s--- Swap ---%s\n", "", colorReset)
				fmt.Printf("Total: %s\n", formatBytes(memInfo.SwapTotal))
				fmt.Printf("Used:  %s\n", formatBytes(memInfo.SwapUsed))
				fmt.Printf("Free:  %s\n", formatBytes(memInfo.SwapFree))
			}

			// macOS specific
			if memInfo.Wired > 0 || memInfo.Compressed > 0 {
				fmt.Println()
				fmt.Printf("%s--- macOS Specific ---%s\n", "", colorReset)
				if memInfo.Wired > 0 {
					fmt.Printf("Wired:     %s\n", formatBytes(memInfo.Wired))
				}
				if memInfo.Compressed > 0 {
					fmt.Printf("Compressed: %s\n", formatBytes(memInfo.Compressed))
				}
			}
		}

		if liveConfig.Enabled {
			sigChan, cleanup := SetupLiveMode()
			defer cleanup()
			for {
				runMem()
				select {
				case <-sigChan:
					fmt.Println("\nInterrupted, exiting...")
					return
				case <-time.After(liveConfig.Interval):
				}
			}
		} else {
			runMem()
		}
	},
}

func init() {
	rootCmd.AddCommand(memCmd)
}
