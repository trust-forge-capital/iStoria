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
			noColorFlag := noColor

			// Host line with alignment
			hostStr := fmt.Sprintf("%s %s/%s | Up:%s | %s",
				Colorize("●", ColorGreen, noColorFlag),
				platform.OS, platform.Arch,
				statInfo.Uptime,
				platform.Kernel)
			fmt.Printf("%-20s %s\n", Colorize(platform.Hostname, ColorBold, noColorFlag), hostStr)

			// CPU line
			cpuModel := cpuInfo.Model
			if len(cpuModel) > 20 {
				cpuModel = cpuModel[:17] + "..."
			}
			cpuCores := fmt.Sprintf("%dC/%dT", cpuInfo.Cores, cpuInfo.Threads)
			var cpuDetails string
			if cpuInfo.AppleSilicon {
				cpuDetails = fmt.Sprintf("P:%d E:%d", cpuInfo.PerformanceCores, cpuInfo.EfficiencyCores)
			}
			cpuUsage := FmtPercent(cpuInfo.UsagePercent, noColorFlag)
			fmt.Printf(" CPU %-20s %-12s %-12s %s\n",
				Colorize(cpuModel, ColorCyan, noColorFlag),
				cpuCores,
				cpuDetails,
				cpuUsage)

			// Memory line with usage bar
			memUsedStr := fmt.Sprintf("%s/%s", formatBytes(memInfo.Used), formatBytes(memInfo.Total))
			memBar := FmtBar(memInfo.UsedPercent, 10, noColorFlag)
			fmt.Printf(" Mem %-20s %s %s %s\n",
				Colorize(memUsedStr, ColorYellow, noColorFlag),
				memBar,
				FmtPercent(memInfo.UsedPercent, noColorFlag),
				Colorize(formatBytes(memInfo.Available)+" avail", ColorGray, noColorFlag))

			// Disk line (if available)
			if len(diskInfo.Disks) > 0 {
				d := diskInfo.Disks[0]
				diskUsedStr := fmt.Sprintf("%s/%s", formatBytes(d.Used), formatBytes(d.Total))
				diskBar := FmtBar(d.UsedPercent, 10, noColorFlag)
				fmt.Printf(" Disk %-20s %s %s %s\n",
					Colorize(d.Path, ColorCyan, noColorFlag),
					diskUsedStr,
					diskBar,
					FmtPercent(d.UsedPercent, noColorFlag))
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
