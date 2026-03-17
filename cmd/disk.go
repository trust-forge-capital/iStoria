package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Show disk information",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		liveConfig := GetLiveConfig(cmd)
		c := collect.NewCollector()

		runDisk := func() {
			if liveConfig.Enabled {
				if !liveConfig.NoClear {
					for i := 0; i < 8; i++ {
						fmt.Print("\r\033[K")
						if i < 7 {
							fmt.Print("\033[A")
						}
					}
					fmt.Print("\r")
				} else {
					fmt.Println()
				}
			}

			diskInfo, err := c.Disk()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting disk: %v\n", err)
				os.Exit(1)
			}

			if jsonFlag {
				j, _ := json.MarshalIndent(diskInfo, "", "  ")
				fmt.Println(string(j))
				return
			}

			_ = noColor // reserved for future color support

			fmt.Printf("%-10s %8s %8s %8s %5s %s\n", "Disk", "Total", "Used", "Avail", "Use%", "Type")
			for _, d := range diskInfo.Disks {
				fmt.Printf("%-10s %8s %8s %8s %4.0f%% %s\n",
					d.Path,
					formatBytes(d.Total),
					formatBytes(d.Used),
					formatBytes(d.Available),
					d.UsedPercent,
					d.Filesystem)
			}
		}

		if liveConfig.Enabled {
			sigChan, cleanup := SetupLiveMode()
			defer cleanup()
			for {
				runDisk()
				select {
				case <-sigChan:
					fmt.Println("\nInterrupted, exiting...")
					return
				case <-time.After(liveConfig.Interval):
				}
			}
		} else {
			runDisk()
		}
	},
}

func init() {
	rootCmd.AddCommand(diskCmd)
}
