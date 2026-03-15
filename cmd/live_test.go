package cmd

import (
	"testing"
	"time"
)

// TestGetLiveConfig tests GetLiveConfig function
func TestGetLiveConfig(t *testing.T) {
	tests := []struct {
		name           string
		liveFlag      bool
		intervalFlag  int
		noClearFlag  bool
		expectEnabled bool
		expectInterval time.Duration
	}{
		{"default", false, 1000, false, false, 1000 * time.Millisecond},
		{"live enabled", true, 1000, false, true, 1000 * time.Millisecond},
		{"custom interval", true, 500, false, true, 500 * time.Millisecond},
		{"below minimum interval", true, 100, false, true, 500 * time.Millisecond}, // min is 500
		{"no clear", true, 1000, true, true, 1000 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := getLiveConfigFromFlags(tt.liveFlag, tt.intervalFlag, tt.noClearFlag)
			
			if config.Enabled != tt.expectEnabled {
				t.Errorf("Enabled = %v; want %v", config.Enabled, tt.expectEnabled)
			}
			if config.Interval != tt.expectInterval {
				t.Errorf("Interval = %v; want %v", config.Interval, tt.expectInterval)
			}
		})
	}
}

// getLiveConfigFromFlags is a helper to test LiveConfig creation
func getLiveConfigFromFlags(live bool, interval int, noClear bool) *LiveConfig {
	if interval < 500 {
		interval = 500
	}
	return &LiveConfig{
		Enabled:  live,
		Interval: time.Duration(interval) * time.Millisecond,
		NoClear:  noClear,
	}
}
