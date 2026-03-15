package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/maxzhang/istoria/internal/collect"
)

// LiveRenderer handles live mode output rendering
type LiveRenderer struct {
	collector collect.Collector
	config    *LiveConfig
	lastRx    uint64
	lastTx    uint64
	lastTime  time.Time
}

// NewLiveRenderer creates a new live renderer
func NewLiveRenderer(c collect.Collector, config *LiveConfig) *LiveRenderer {
	return &LiveRenderer{
		collector: c,
		config:    config,
		lastRx:    0,
		lastTx:    0,
		lastTime:  time.Now(),
	}
}

// RunNetLive runs the live network monitoring
func RunNetLive(config *LiveConfig) error {
	c := collect.NewCollector()
	renderer := NewLiveRenderer(c, config)

	sigChan, cleanup := SetupLiveMode()
	defer cleanup()

	for {
		select {
		case <-sigChan:
			// Handle Ctrl+C gracefully
			fmt.Println("\nExiting live mode...")
			return nil

		case <-time.After(config.Interval):
			// Collect and render
			if err := renderer.renderNet(); err != nil {
				return err
			}
		}
	}
}

func (r *LiveRenderer) renderNet() error {
	netData, err := r.collector.Net()
	if err != nil {
		return err
	}

	now := time.Now()
	elapsed := now.Sub(r.lastTime).Seconds()

	// Calculate rates
	var rxRate, txRate float64
	if r.lastRx > 0 && elapsed > 0 {
		rxRate = float64(netData.TotalRx-r.lastRx) / elapsed
		txRate = float64(netData.TotalTx-r.lastTx) / elapsed
	}

	r.lastRx = netData.TotalRx
	r.lastTx = netData.TotalTx
	r.lastTime = now

	// Clear screen if not no-clear mode
	if !r.config.NoClear {
		ClearScreen()
	}

	// Print header
	fmt.Println("=== iStoria Live Network Monitor ===")
	fmt.Printf("Press Ctrl+C to exit | Interval: %v | Mode: %s\n\n",
		r.config.Interval, map[bool]string{false: "refresh", true: "append"}[r.config.NoClear])

	// Print interface table
	fmt.Printf("%-10s %15s %15s %15s %15s\n",
		"Interface", "IPv4", "RX Rate", "TX Rate", "Total RX")
	fmt.Printf("%s\n", "--------------------------------------------------------------------------------")

	for _, ni := range netData.Interfaces {
		if ni.Name == "lo0" || ni.Name == "lo" {
			continue
		}
		fmt.Printf("%-10s %15s %15s %15s %15s\n",
			ni.Name,
			ni.IP4,
			formatBytesRate(uint64(rxRate)),
			formatBytesRate(uint64(txRate)),
			formatBytes(ni.RxBytes))
	}

	fmt.Println()
	fmt.Printf("Total RX: %s | Total TX: %s\n",
		formatBytes(netData.TotalRx),
		formatBytes(netData.TotalTx))
	fmt.Printf("RX Rate: %s/s | TX Rate: %s/s\n",
		formatBytesRate(uint64(rxRate)),
		formatBytesRate(uint64(txRate)))

	return nil
}

// RunNetLiveJSON runs live network monitoring in JSON mode
func RunNetLiveJSON(config *LiveConfig) error {
	c := collect.NewCollector()

	sigChan, cleanup := SetupLiveMode()
	defer cleanup()

	for {
		select {
		case <-sigChan:
			return nil

		case <-time.After(config.Interval):
			netData, err := c.Net()
			if err != nil {
				continue
			}

			// Calculate rates
			now := time.Now()
			output := map[string]interface{}{
				"timestamp":      now.Format(time.RFC3339),
				"command":        "net",
				"total_rx_bytes": netData.TotalRx,
				"total_tx_bytes": netData.TotalTx,
				"interfaces":     netData.Interfaces,
			}

			data, err := json.Marshal(output)
			if err == nil {
				fmt.Println(string(data))
			}
		}
	}
}

func formatBytesRate(b uint64) string {
	const unit = 1024.0
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := unit, 0.0
	for float64(b)/div >= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB/s", float64(b)/div, "KMGTPE"[int(exp)])
}
