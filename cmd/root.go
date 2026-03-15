package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "istoria",
	Short: "iStoria is a hardware monitoring CLI for SSH-first workflows",
	Long:  "iStoria helps you inspect local machine health from the terminal with human-friendly and JSON output.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "output in JSON format")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")
	rootCmd.PersistentFlags().Bool("quiet", false, "reduce non-essential output")
	rootCmd.PersistentFlags().String("config", "", "path to config file")

	// Live mode flags
	rootCmd.PersistentFlags().Bool("live", false, "enable live refresh mode")
	rootCmd.PersistentFlags().Int("interval", 1000, "refresh interval in milliseconds (default: 1000, min: 500)")
	rootCmd.PersistentFlags().Bool("no-clear", false, "don't clear screen in live mode (append output)")
}

// LiveConfig holds live mode configuration
type LiveConfig struct {
	Enabled   bool
	Interval  time.Duration
	NoClear   bool
}

// GetLiveConfig extracts live mode configuration from command
func GetLiveConfig(cmd *cobra.Command) *LiveConfig {
	live, _ := cmd.Flags().GetBool("live")
	interval, _ := cmd.Flags().GetInt("interval")
	noClear, _ := cmd.Flags().GetBool("no-clear")

	// Validate interval
	if interval < 500 {
		interval = 500 // minimum 500ms
	}

	return &LiveConfig{
		Enabled:   live,
		Interval:  time.Duration(interval) * time.Millisecond,
		NoClear:   noClear,
	}
}

// SetupLiveMode sets up signal handling for graceful exit in live mode
func SetupLiveMode() (chan os.Signal, func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	cleanup := func() {
		signal.Stop(sigChan)
		close(sigChan)
	}

	return sigChan, cleanup
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}
