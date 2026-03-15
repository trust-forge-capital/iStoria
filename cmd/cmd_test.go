package cmd

import (
	"testing"
)

// TestCPUCmdExists tests that cpu command is registered
func TestCPUCmdExists(t *testing.T) {
	// Test that cpu command can be created without panic
	cmd := cpuCmd
	if cmd == nil {
		t.Error("cpuCmd should not be nil")
	}
	if cmd.Use != "cpu" {
		t.Errorf("cpuCmd.Use = %s; want cpu", cmd.Use)
	}
}

// TestMemCmdExists tests that mem command is registered
func TestMemCmdExists(t *testing.T) {
	cmd := memCmd
	if cmd == nil {
		t.Error("memCmd should not be nil")
	}
	if cmd.Use != "mem" {
		t.Errorf("memCmd.Use = %s; want mem", cmd.Use)
	}
}

// TestDiskCmdExists tests that disk command is registered
func TestDiskCmdExists(t *testing.T) {
	cmd := diskCmd
	if cmd == nil {
		t.Error("diskCmd should not be nil")
	}
	if cmd.Use != "disk" {
		t.Errorf("diskCmd.Use = %s; want disk", cmd.Use)
	}
}

// TestNetCmdExists tests that net command is registered
func TestNetCmdExists(t *testing.T) {
	cmd := netCmd
	if cmd == nil {
		t.Error("netCmd should not be nil")
	}
	if cmd.Use != "net" {
		t.Errorf("netCmd.Use = %s; want net", cmd.Use)
	}
}

// TestSensorCmdExists tests that sensor command is registered
func TestSensorCmdExists(t *testing.T) {
	cmd := sensorCmd
	if cmd == nil {
		t.Error("sensorCmd should not be nil")
	}
	if cmd.Use != "sensor" {
		t.Errorf("sensorCmd.Use = %s; want sensor", cmd.Use)
	}
}

// TestPowerCmdExists tests that power command is registered
func TestPowerCmdExists(t *testing.T) {
	cmd := powerCmd
	if cmd == nil {
		t.Error("powerCmd should not be nil")
	}
	if cmd.Use != "power" {
		t.Errorf("powerCmd.Use = %s; want power", cmd.Use)
	}
}
