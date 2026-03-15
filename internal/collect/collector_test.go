package collect

import (
	"encoding/json"
	"os"
	"testing"
)

func init() {
	// Fix macOS PATH for gopsutil in tests
	os.Setenv("PATH", "/usr/sbin:/usr/bin:"+os.Getenv("PATH"))
}

// TestCollectorInitialization tests that NewCollector returns a valid collector
func TestCollectorInitialization(t *testing.T) {
	c := NewCollector()
	if c == nil {
		t.Fatal("NewCollector returned nil")
	}

	// Test that Platform() returns valid data
	platform, err := c.Platform()
	if err != nil {
		t.Skipf("Platform() returned error (may need ioreg): %v", err)
		return
	}
	if platform.OS == "" {
		t.Error("Platform.OS should not be empty")
	}
	if platform.Hostname == "" {
		t.Error("Platform.Hostname should not be empty")
	}
}

// TestMemInfoJSON tests that MemInfo marshals to expected JSON keys
func TestMemInfoJSON(t *testing.T) {
	mem := MemInfo{
		Total:       16000000000,
		Available:   8000000000,
		Used:        8000000000,
		Free:        4000000000,
		UsedPercent: 50.0,
		SwapTotal:   2000000000,
		SwapUsed:    500000000,
		SwapFree:    1500000000,
	}

	// Test JSON marshaling
	data, err := json.Marshal(mem)
	if err != nil {
		t.Fatalf("JSON Marshal failed: %v", err)
	}

	// Verify JSON contains expected snake_case keys
	var result map[string]interface{}
	json.Unmarshal(data, &result)

	requiredKeys := []string{"total_bytes", "available_bytes", "used_percent", "swap_total_bytes"}
	for _, key := range requiredKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("JSON should contain key: %s", key)
		}
	}
}

// TestCPUInfoJSON tests that CPUInfo marshals to expected JSON keys
func TestCPUInfoJSON(t *testing.T) {
	cpu := CPUInfo{
		Model:            "Apple M4",
		Cores:            10,
		Threads:          20,
		UsagePercent:     25.5,
		AppleSilicon:     true,
		PerformanceCores: 4,
		EfficiencyCores:  6,
	}

	data, err := json.Marshal(cpu)
	if err != nil {
		t.Fatalf("JSON Marshal failed: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	requiredKeys := []string{"usage_percent", "apple_silicon", "performance_cores"}
	for _, key := range requiredKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("JSON should contain key: %s", key)
		}
	}
}

// TestDiskInfoJSON tests that DiskInfo marshals to expected JSON keys
func TestDiskInfoJSON(t *testing.T) {
	disk := DiskInfo{
		Path:        "/",
		Filesystem:  "apfs",
		Total:       500000000000,
		UsedPercent: 50.0,
	}

	data, err := json.Marshal(disk)
	if err != nil {
		t.Fatalf("JSON Marshal failed: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	if _, ok := result["total_bytes"]; !ok {
		t.Error("JSON should contain total_bytes")
	}
	if _, ok := result["used_percent"]; !ok {
		t.Error("JSON should contain used_percent")
	}
}

// TestPowerInfoJSON tests that PowerInfo marshals to expected JSON keys
func TestPowerInfoJSON(t *testing.T) {
	power := PowerInfo{
		HasBattery:   true,
		Percent:      75,
		PowerPlugged: false,
		CycleCount:   342,
	}

	data, err := json.Marshal(power)
	if err != nil {
		t.Fatalf("JSON Marshal failed: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	requiredKeys := []string{"has_battery", "battery_percent", "cycle_count"}
	for _, key := range requiredKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("JSON should contain key: %s", key)
		}
	}
}

// TestSensorDataJSON tests that SensorData marshals to expected JSON keys
func TestSensorDataJSON(t *testing.T) {
	sensors := SensorData{
		Temperatures: []SensorInfo{
			{Name: "CPU", Value: 52.5, Unit: "°C", SensorType: "temperature"},
		},
		Fans: []SensorInfo{
			{Name: "Fan 1", Value: 1200, Unit: "RPM", SensorType: "fan"},
		},
		HasSensors: true,
	}

	data, err := json.Marshal(sensors)
	if err != nil {
		t.Fatalf("JSON Marshal failed: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	requiredKeys := []string{"temperatures", "fans", "has_sensors"}
	for _, key := range requiredKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("JSON should contain key: %s", key)
		}
	}
}

// TestPlatformInfoJSON tests that PlatformInfo marshals to expected JSON keys
func TestPlatformInfoJSON(t *testing.T) {
	platform := PlatformInfo{
		OS:           "darwin",
		Arch:         "arm64",
		Hostname:     "test-mac",
		Uptime:       3600,
		AppleSilicon: true,
	}

	data, err := json.Marshal(platform)
	if err != nil {
		t.Fatalf("JSON Marshal failed: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	requiredKeys := []string{"boot_time", "apple_silicon"}
	for _, key := range requiredKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("JSON should contain key: %s", key)
		}
	}
}
