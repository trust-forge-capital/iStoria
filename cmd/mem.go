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
			if liveConfig.Enabled {
				if !liveConfig.NoClear {
					for i := 0; i < 3; i++ {
						fmt.Print("\r\033[K")
						if i < 2 {
							fmt.Print("\033[A")
						}
					}
					fmt.Print("\r")
				} else {
					fmt.Println()
				}
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
			noColorFlag := noColor

			// Memory line
			memUsedStr := fmt.Sprintf("%s/%s", formatBytes(memInfo.Used), formatBytes(memInfo.Total))
			memBar := FmtBar(memInfo.UsedPercent, 10, noColorFlag)
			fmt.Printf(" %-4s %-20s %s %s %s\n",
				Colorize("Mem", ColorBold, noColorFlag),
				Colorize(memUsedStr, ColorYellow, noColorFlag),
				memBar,
				FmtPercent(memInfo.UsedPercent, noColorFlag),
				Colorize(formatBytes(memInfo.Available)+" avail", ColorGray, noColorFlag))

			// Swap line
			if memInfo.SwapTotal > 0 {
				swapBar := FmtBar(float64(memInfo.SwapUsed)*100/float64(memInfo.SwapTotal), 10, noColorFlag)
				fmt.Printf(" %-4s %-20s %s %s\n",
					Colorize("Swap", ColorBold, noColorFlag),
					fmt.Sprintf("%s/%s", formatBytes(memInfo.SwapUsed), formatBytes(memInfo.SwapTotal)),
					swapBar,
					FmtPercent(float64(memInfo.SwapUsed)*100/float64(memInfo.SwapTotal), noColorFlag))
			}

			// macOS specific
			if memInfo.Wired > 0 || memInfo.Compressed > 0 {
				extra := ""
				if memInfo.Wired > 0 {
					extra += fmt.Sprintf("Wired:%s", formatBytes(memInfo.Wired))
				}
				if memInfo.Compressed > 0 {
					if extra != "" {
						extra += " "
					}
					extra += fmt.Sprintf("Comp:%s", formatBytes(memInfo.Compressed))
				}
				fmt.Printf(" %-4s %s\n",
					Colorize("OS", ColorBold, noColorFlag),
					Colorize(extra, ColorGray, noColorFlag))
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
