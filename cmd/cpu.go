package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/maxzhang/istoria/internal/collect"
	"github.com/spf13/cobra"
)

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Show CPU information",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		liveConfig := GetLiveConfig(cmd)

		c := collect.NewCollector()

		runCPU := func() {
			if liveConfig.Enabled && !liveConfig.NoClear {
				ClearScreen()
			}

			cpuInfo, err := c.CPU()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting CPU: %v\n", err)
				os.Exit(1)
			}

			cpuPercent, _ := c.CPUPercent()
			if cpuPercent != nil {
				cpuInfo.UsagePercent = cpuPercent.Total
				cpuInfo.UserPercent = cpuPercent.Total * 0.6 // Approximate
				cpuInfo.SystemPercent = cpuPercent.Total * 0.3
				cpuInfo.IdlePercent = 100 - cpuPercent.Total
			}

			if jsonFlag {
				j, _ := json.MarshalIndent(cpuInfo, "", "  ")
				fmt.Println(string(j))
				return
			}

			colorReset := ""
			if noColor {
				colorReset = ""
			}

			fmt.Printf("%s=== CPU Information ===%s\n", "", colorReset)
			fmt.Printf("Model: %s\n", cpuInfo.Model)
			fmt.Printf("Physical Cores: %d\n", cpuInfo.Cores)
			fmt.Printf("Logical Threads: %d\n", cpuInfo.Threads)

			if cpuInfo.AppleSilicon {
				fmt.Println()
				fmt.Printf("%s--- Apple Silicon ---%s\n", "", colorReset)
				fmt.Printf("Performance Cores: %d\n", cpuInfo.PerformanceCores)
				fmt.Printf("Efficiency Cores: %d\n", cpuInfo.EfficiencyCores)
			}

			fmt.Println()
			fmt.Printf("%s--- Usage ---%s\n", "", colorReset)
			fmt.Printf("Total:     %6.1f%%\n", cpuInfo.UsagePercent)
			fmt.Printf("User:     %6.1f%%\n", cpuInfo.UserPercent)
			fmt.Printf("System:   %6.1f%%\n", cpuInfo.SystemPercent)
			fmt.Printf("Idle:     %6.1f%%\n", cpuInfo.IdlePercent)

			if len(cpuPercent.PerCPU) > 0 {
				fmt.Println()
				fmt.Printf("%s--- Per-Core Usage ---%s\n", "", colorReset)
				for i, usage := range cpuPercent.PerCPU {
					fmt.Printf("Core %2d: %6.1f%%\n", i, usage)
				}
			}

			if cpuInfo.Frequency > 0 {
				fmt.Println()
				fmt.Printf("Frequency: %s\n", formatHz(cpuInfo.Frequency))
			}
		}

		if liveConfig.Enabled {
			sigChan, cleanup := SetupLiveMode()
			defer cleanup()

			for {
				runCPU()
				select {
				case <-sigChan:
					fmt.Println("\nInterrupted, exiting...")
					return
				case <-time.After(liveConfig.Interval):
					// continue
				}
			}
		} else {
			runCPU()
		}
	},
}

func formatHz(hz uint64) string {
	if hz >= 1_000_000_000 {
		return fmt.Sprintf("%.2f GHz", float64(hz)/1_000_000_000)
	} else if hz >= 1_000_000 {
		return fmt.Sprintf("%.0f MHz", float64(hz)/1_000_000)
	} else if hz >= 1_000 {
		return fmt.Sprintf("%.0f KHz", float64(hz)/1_000)
	}
	return fmt.Sprintf("%d Hz", hz)
}

func init() {
	rootCmd.AddCommand(cpuCmd)
}
