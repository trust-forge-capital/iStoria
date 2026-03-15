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
			if liveConfig.Enabled && !liveConfig.NoClear {
				ClearScreen()
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

			colorReset := ""
			if noColor {
				colorReset = ""
			}

			fmt.Printf("%s=== Disk Information ===%s\n", "", colorReset)
			fmt.Printf("%-12s %10s %10s %10s %8s %s\n",
				"Mount", "Total", "Used", "Available", "Use%", "Filesystem")
			fmt.Println()

			for _, d := range diskInfo.Disks {
				_ = generateBar(d.UsedPercent)
				fmt.Printf("%-12s %10s %10s %10s %7.1f%% %s\n",
					d.Path,
					formatBytes(d.Total),
					formatBytes(d.Used),
					formatBytes(d.Available),
					d.UsedPercent,
					d.Filesystem)

				if d.InodesTotal > 0 {
					inodePercent := float64(d.InodesUsed) / float64(d.InodesTotal) * 100
					fmt.Printf("  Inodes: %s / %s (%.1f%%)\n",
						formatBytes(d.InodesUsed*1024),
						formatBytes(d.InodesTotal*1024),
						inodePercent)
				}
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

func generateBar(percent float64) string {
	const barLength = 20
	filled := int(percent / 100 * float64(barLength))
	bar := ""
	for i := 0; i < barLength; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

func init() {
	rootCmd.AddCommand(diskCmd)
}
