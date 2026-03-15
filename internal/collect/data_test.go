package collect

import (
	"testing"
)

// TestCollectorCreation tests NewCollector function
func TestCollectorCreation(t *testing.T) {
	collector := NewCollector()
	if collector == nil {
		t.Error("NewCollector() returned nil")
	}
}

// TestCPUInfoFields tests CPUInfo structure
func TestCPUInfoFields(t *testing.T) {
	info := CPUInfo{
		Model:            "Apple M4",
		Cores:            10,
		Threads:          20,
		UsagePercent:    50.0,
		UserPercent:      30.0,
		SystemPercent:    20.0,
		IdlePercent:      50.0,
		AppleSilicon:    true,
		PerformanceCores: 4,
		EfficiencyCores:  6,
	}

	if info.Model != "Apple M4" {
		t.Errorf("Model = %s; want Apple M4", info.Model)
	}
	if info.Cores != 10 {
		t.Errorf("Cores = %d; want 10", info.Cores)
	}
	if !info.AppleSilicon {
		t.Error("AppleSilicon should be true")
	}
}

// TestMemInfoFields tests MemInfo structure
func TestMemInfoFields(t *testing.T) {
	info := MemInfo{
		Total:         16000000000,
		Used:          8000000000,
		Free:          4000000000,
		Available:     8000000000,
		UsedPercent:   50.0,
		SwapTotal:     2000000000,
		SwapUsed:      1000000000,
		SwapFree:      1000000000,
		Wired:         2000000000,
		Compressed:    2000000000,
	}

	if info.Total != 16000000000 {
		t.Errorf("Total = %d; want 16000000000", info.Total)
	}
	if info.UsedPercent != 50.0 {
		t.Errorf("UsedPercent = %f; want 50.0", info.UsedPercent)
	}
}

// TestDiskInfoFields tests DiskInfo structure
func TestDiskInfoFields(t *testing.T) {
	info := DiskInfo{
		Path:         "/",
		Filesystem:   "apfs",
		Total:        500000000000,
		Used:         250000000000,
		Available:    250000000000,
		UsedPercent:  50.0,
		InodesTotal:  1000000,
		InodesUsed:   500000,
		InodesAvail:  500000,
	}

	if info.Path != "/" {
		t.Errorf("Path = %s; want /", info.Path)
	}
	if info.UsedPercent != 50.0 {
		t.Errorf("UsedPercent = %f; want 50.0", info.UsedPercent)
	}
}

// TestPlatformInfoFields tests PlatformInfo structure
func TestPlatformInfoFields(t *testing.T) {
	info := PlatformInfo{
		OS:            "darwin",
		Arch:          "arm64",
		Hostname:      "test-mac",
		Kernel:        "23.0.0",
		Uptime:        3600,
		BootTime:      1700000000,
		Procs:         100,
		AppleSilicon:  true,
	}

	if info.OS != "darwin" {
		t.Errorf("OS = %s; want darwin", info.OS)
	}
	if !info.AppleSilicon {
		t.Error("AppleSilicon should be true")
	}
}

// TestNetInfoFields tests NetInfo structure
func TestNetInfoFields(t *testing.T) {
	info := NetInfo{
		Name:        "en0",
		MTU:         1500,
		HardwareAddr: "aa:bb:cc:dd:ee:ff",
		Flags:       "up,running",
		IP4:         "192.168.1.100",
		IP6:         "fd00::1",
		RxBytes:     1000000,
		TxBytes:     2000000,
	}

	if info.Name != "en0" {
		t.Errorf("Name = %s; want en0", info.Name)
	}
	if info.RxBytes != 1000000 {
		t.Errorf("RxBytes = %d; want 1000000", info.RxBytes)
	}
}

// TestPowerInfoFields tests PowerInfo structure
func TestPowerInfoFields(t *testing.T) {
	info := PowerInfo{
		HasBattery:    true,
		Charging:      true,
		Percent:       80,
		TimeRemaining: 300,
		PowerPlugged:  true,
		Amps:          1.5,
		Volts:         12.0,
		Watts:         18.0,
		CycleCount:    500,
		Health:        "Good",
	}

	if !info.HasBattery {
		t.Error("HasBattery should be true")
	}
	if info.Percent != 80 {
		t.Errorf("Percent = %d; want 80", info.Percent)
	}
}

// TestSensorDataFields tests SensorData structure
func TestSensorDataFields(t *testing.T) {
	info := SensorData{
		Temperatures: []SensorInfo{
			{Name: "CPU", Value: 65.0, Unit: "°C", Critical: false},
		},
		Fans:       []SensorInfo{},
		Voltages:   []SensorInfo{},
		Power:      []SensorInfo{},
		HasSensors: true,
	}

	if len(info.Temperatures) != 1 {
		t.Errorf("Temperatures length = %d; want 1", len(info.Temperatures))
	}
	if !info.HasSensors {
		t.Error("HasSensors should be true")
	}
}
