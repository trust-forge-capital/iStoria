package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show machine summary",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		liveConfig := GetLiveConfig(cmd)

		c := collect.NewCollector()

		runStat := func() {
			if liveConfig.Enabled {
				if !liveConfig.NoClear {
					// Use carriage return to go to line start (more compatible than ANSI)
					for i := 0; i < 5; i++ {
						fmt.Print("\r\033[K") // Clear line and return to start
						if i < 4 {
							fmt.Print("\033[A") // Move up one line
						}
					}
					fmt.Print("\r")
				} else {
					fmt.Println()
				}
			}

			platform, err := c.Platform()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting platform: %v\n", err)
				os.Exit(1)
			}

			cpuInfo, err := c.CPU()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting CPU: %v\n", err)
				os.Exit(1)
			}

			cpuPercent, _ := c.CPUPercent()
			if cpuPercent != nil {
				cpuInfo.UsagePercent = cpuPercent.Total
			}

			memInfo, err := c.Mem()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting memory: %v\n", err)
				os.Exit(1)
			}

			diskInfo, err := c.Disk()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting disk: %v\n", err)
				os.Exit(1)
			}

			statInfo := collect.StatInfo{
				Platform: *platform,
				CPU:      *cpuInfo,
				Mem:      *memInfo,
			}

			if len(diskInfo.Disks) > 0 {
				statInfo.Disk = diskInfo.Disks[0]
			}

			statInfo.Uptime = formatUptime(platform.Uptime)

			if platform.OS == "linux" {
				statInfo.LoadAvg = getLoadAvg()
			}

			if jsonFlag {
				j, _ := json.MarshalIndent(statInfo, "", "  ")
				fmt.Println(string(j))
				return
			}

			_ = noColor // reserved for future color support

			fmt.Printf("Host: %s | OS: %s/%s | Up: %s | Kernel: %s\n", platform.Hostname, platform.OS, platform.Arch, statInfo.Uptime, platform.Kernel)

			cpuLine := fmt.Sprintf("CPU: %s | %dC/%dT", cpuInfo.Model, cpuInfo.Cores, cpuInfo.Threads)
			if cpuInfo.AppleSilicon {
				cpuLine += fmt.Sprintf(" | P:%d E:%d", cpuInfo.PerformanceCores, cpuInfo.EfficiencyCores)
			}
			cpuLine += fmt.Sprintf(" | %.1f%%", cpuInfo.UsagePercent)
			fmt.Println(cpuLine)

			fmt.Printf("Mem: %s | Used:%s(%.0f%%) | Avail:%s\n",
				formatBytes(memInfo.Total),
				formatBytes(memInfo.Used),
				memInfo.UsedPercent,
				formatBytes(memInfo.Available))

			if len(diskInfo.Disks) > 0 {
				d := diskInfo.Disks[0]
				fmt.Printf("Disk: %s | %s | Used:%s(%.0f%%) | Avail:%s\n",
					d.Path,
					formatBytes(d.Total),
					formatBytes(d.Used),
					d.UsedPercent,
					formatBytes(d.Available))
			}

			if statInfo.LoadAvg != "" {
				fmt.Printf("Load: %s\n", statInfo.LoadAvg)
			}
		}

		if liveConfig.Enabled {
			sigChan, cleanup := SetupLiveMode()
			defer cleanup()

			for {
				runStat()
				select {
				case <-sigChan:
					fmt.Println("\nInterrupted, exiting...")
					return
				case <-time.After(liveConfig.Interval):
					// continue
				}
			}
		} else {
			runStat()
		}
	},
}

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	mins := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	} else {
		return fmt.Sprintf("%dm", mins)
	}
}

func getLoadAvg() string {
	// Simple load average for Linux - in real implementation use gopsutil/host
	return ""
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func init() {
	rootCmd.AddCommand(statCmd)
}
