package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show machine summary",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")

		c := collect.NewCollector()

		// Collect all data
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

		stat := collect.StatInfo{
			Platform: *platform,
			CPU:      *cpuInfo,
			Mem:      *memInfo,
		}

		// Get root disk
		if len(diskInfo.Disks) > 0 {
			stat.Disk = diskInfo.Disks[0]
		}

		// Format uptime
		stat.Uptime = formatUptime(platform.Uptime)

		// Get load average on Linux
		if platform.OS == "linux" {
			stat.LoadAvg = getLoadAvg()
		}

		if jsonFlag {
			j, _ := json.MarshalIndent(stat, "", "  ")
			fmt.Println(string(j))
			return
		}

		// Human readable output (compact)
		_ = noColor // reserved for future color support

		fmt.Printf("Host: %s | OS: %s/%s | Up: %s | Kernel: %s\n", platform.Hostname, platform.OS, platform.Arch, stat.Uptime, platform.Kernel)

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

		if stat.LoadAvg != "" {
			fmt.Printf("Load: %s\n", stat.LoadAvg)
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
