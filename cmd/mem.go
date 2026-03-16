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

			_ = noColor // reserved for future color support

			fmt.Printf("Mem: %s | Used:%s(%.0f%%) | Avail:%s | Free:%s\n",
				formatBytes(memInfo.Total),
				formatBytes(memInfo.Used),
				memInfo.UsedPercent,
				formatBytes(memInfo.Available),
				formatBytes(memInfo.Free))

			if memInfo.SwapTotal > 0 {
				fmt.Printf("Swap: %s | Used:%s | Free:%s\n",
					formatBytes(memInfo.SwapTotal),
					formatBytes(memInfo.SwapUsed),
					formatBytes(memInfo.SwapFree))
			}

			if memInfo.Wired > 0 || memInfo.Compressed > 0 {
				extra := ""
				if memInfo.Wired > 0 {
					extra += fmt.Sprintf(" Wired:%s", formatBytes(memInfo.Wired))
				}
				if memInfo.Compressed > 0 {
					extra += fmt.Sprintf(" Comp:%s", formatBytes(memInfo.Compressed))
				}
				fmt.Printf("macOS:%s\n", extra)
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
