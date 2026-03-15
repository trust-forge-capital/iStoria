package cmd

import (
	"testing"
)

// TestFormatBytes tests the formatBytes helper function
func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
		{1610612736, "1.5 GB"},
	}

	for _, tt := range tests {
		result := formatBytes(tt.input)
		if result != tt.expected {
			t.Errorf("formatBytes(%d) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

// TestFormatHz tests the formatHz helper function
func TestFormatHz(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0 Hz"},
		{1000, "1 KHz"},
		{1000000, "1 MHz"},
		{1000000000, "1.00 GHz"},
		{2500000000, "2.50 GHz"},
	}

	for _, tt := range tests {
		result := formatHz(tt.input)
		if result != tt.expected {
			t.Errorf("formatHz(%d) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

// TestFormatUptime tests the formatUptime helper function
func TestFormatUptime(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0m"},
		{60, "1m"},
		{3600, "1h 0m"},
		{3660, "1h 1m"},
		{86400, "1d 0h 0m"},
		{90060, "1d 1h 1m"},
	}

	for _, tt := range tests {
		result := formatUptime(tt.input)
		if result != tt.expected {
			t.Errorf("formatUptime(%d) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}
