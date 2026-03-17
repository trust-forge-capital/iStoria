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
			if liveConfig.Enabled {
				if !liveConfig.NoClear {
					for i := 0; i < 4; i++ {
						fmt.Print("\r\033[K")
						if i < 3 {
							fmt.Print("\033[A")
						}
					}
					fmt.Print("\r")
				} else {
					fmt.Println()
				}
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

			_ = noColor // reserved for future color support
			noColorFlag := noColor

			// CPU model line
			cpuModel := cpuInfo.Model
			if len(cpuModel) > 25 {
				cpuModel = cpuModel[:22] + "..."
			}
			cpuCores := fmt.Sprintf("%dC/%dT", cpuInfo.Cores, cpuInfo.Threads)
			var cpuDetails string
			if cpuInfo.AppleSilicon {
				cpuDetails = fmt.Sprintf("P:%d E:%d", cpuInfo.PerformanceCores, cpuInfo.EfficiencyCores)
			}
			if cpuInfo.Frequency > 0 {
				cpuDetails += " " + formatHz(cpuInfo.Frequency)
			}
			fmt.Printf(" %-3s %-25s %-12s %s\n",
				Colorize("CPU", ColorBold, noColorFlag),
				Colorize(cpuModel, ColorCyan, noColorFlag),
				cpuCores,
				cpuDetails)

			// Usage line
			fmt.Printf(" %-3s %s\n",
				Colorize("Usage", ColorBold, noColorFlag),
				fmt.Sprintf("Total:%s User:%s Sys:%s Idle:%s",
					FmtPercent(cpuInfo.UsagePercent, noColorFlag),
					FmtPercent(cpuInfo.UserPercent, noColorFlag),
					FmtPercent(cpuInfo.SystemPercent, noColorFlag),
					FmtPercent(cpuInfo.IdlePercent, noColorFlag)))

			// Per-core usage
			if cpuPercent != nil && len(cpuPercent.PerCPU) > 0 {
				fmt.Printf(" %-3s ", Colorize("Cores", ColorBold, noColorFlag))
				for i, usage := range cpuPercent.PerCPU {
					if i > 0 {
						fmt.Print(" ")
					}
					fmt.Printf("C%d:%s", i, FmtPercent(usage, noColorFlag))
					if i >= 7 {
						if len(cpuPercent.PerCPU) > 8 {
							fmt.Print(" ...")
						}
						break
					}
				}
				fmt.Println()
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
