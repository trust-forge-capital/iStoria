package collect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// DarwinCollector implements Collector for macOS
type DarwinCollector struct{}

// NewDarwinCollector creates a Darwin collector
func NewDarwinCollector() *DarwinCollector {
	return &DarwinCollector{}
}

// Platform returns platform information
func (c *DarwinCollector) Platform() (*PlatformInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	kernelArch := runtime.GOARCH
	appleSilicon := runtime.GOARCH == "arm64" && runtime.GOOS == "darwin"

	bootTime := uint64(hostInfo.BootTime)
	uptime := uint64(time.Now().Unix()) - bootTime

	return &PlatformInfo{
		OS:           runtime.GOOS,
		Arch:         kernelArch,
		Hostname:     hostInfo.Hostname,
		Platform:     hostInfo.Platform,
		Kernel:       hostInfo.KernelVersion,
		Uptime:       uptime,
		BootTime:     bootTime,
		Procs:        int(hostInfo.Procs),
		AppleSilicon: appleSilicon,
	}, nil
}

// CPU returns CPU information
func (c *DarwinCollector) CPU() (*CPUInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var model string
	var cores, threads int
	if len(cpuInfo) > 0 {
		model = cpuInfo[0].ModelName
		cores = int(cpuInfo[0].Cores)
		threads = int(cpuInfo[0].Cores) * 2 // approximate
	}

	freq, _ := c.getCPUMFreq()
	appleSilicon := runtime.GOARCH == "arm64"
	perfCores, effCores := c.getAppleSiliconCores()

	return &CPUInfo{
		Model:            model,
		Cores:            cores,
		Threads:          threads,
		Frequency:        freq,
		FrequencyMax:     freq,
		AppleSilicon:     appleSilicon,
		PerformanceCores: perfCores,
		EfficiencyCores:  effCores,
	}, nil
}

func (c *DarwinCollector) getCPUMFreq() (uint64, uint64) {
	cmd := exec.Command("sysctl", "-n", "hw.cpufrequency")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		var freq uint64
		fmt.Sscanf(out.String(), "%d", &freq)
		return freq, freq
	}
	return 0, 0
}

func (c *DarwinCollector) getAppleSiliconCores() (int, int) {
	cmd := exec.Command("sysctl", "-n", "hw.perflevel0.physicalcpu")
	var perfOut bytes.Buffer
	cmd.Stdout = &perfOut
	perfErr := cmd.Run()

	cmd = exec.Command("sysctl", "-n", "hw.perflevel1.physicalcpu")
	var effOut bytes.Buffer
	cmd.Stdout = &effOut
	effErr := cmd.Run()

	perfCores, effCores := 0, 0
	if perfErr == nil {
		fmt.Sscanf(perfOut.String(), "%d", &perfCores)
	}
	if effErr == nil {
		fmt.Sscanf(effOut.String(), "%d", &effCores)
	}

	if perfCores == 0 && effCores == 0 {
		info, _ := cpu.Info()
		if len(info) > 0 {
			perfCores = int(info[0].Cores)
		}
	}

	return perfCores, effCores
}

// CPUPercent returns per-CPU usage
func (c *DarwinCollector) CPUPercent() (*CPUPercent, error) {
	perCPU, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	total, err := cpu.Percent(time.Second, true)
	if err != nil {
		return nil, err
	}

	return &CPUPercent{
		PerCPU: perCPU,
		Total:  total[0],
	}, nil
}

// Mem returns memory information
func (c *DarwinCollector) Mem() (*MemInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	s, _ := mem.SwapMemory()

	return &MemInfo{
		Total:       v.Total,
		Available:   v.Available,
		Used:        v.Used,
		Free:        v.Free,
		UsedPercent: v.UsedPercent,
		SwapTotal:   s.Total,
		SwapUsed:    s.Used,
		SwapFree:    s.Free,
		Wired:       v.Wired,
	}, nil
}

// Disk returns disk usage information
func (c *DarwinCollector) Disk() (*DiskUsage, error) {
	parts, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	var disks []DiskInfo
	for _, p := range parts {
		if strings.Contains(p.Mountpoint, "/System/Volumes") {
			continue
		}

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		disks = append(disks, DiskInfo{
			Path:        p.Mountpoint,
			Filesystem:  p.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Available:   usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	if len(disks) == 0 {
		usage, _ := disk.Usage("/")
		disks = append(disks, DiskInfo{
			Path:        "/",
			Filesystem:  "apfs",
			Total:       usage.Total,
			Used:        usage.Used,
			Available:   usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	return &DiskUsage{Disks: disks}, nil
}

// Net returns network interface information
func (c *DarwinCollector) Net() (*NetUsage, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ifs []NetInfo
	var totalRx, totalTx uint64

	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "lo0") {
			continue
		}

		stats, _ := net.IOCounters(false)
		var ifaceStats []net.IOCountersStat
		for _, s := range stats {
			if s.Name == iface.Name {
				ifaceStats = append(ifaceStats, s)
				break
			}
		}

		if len(ifaceStats) > 0 {
			totalRx += ifaceStats[0].BytesRecv
			totalTx += ifaceStats[0].BytesSent
		}

		var ip4, ip6 string
		for _, addr := range iface.Addrs {
			addrStr := addr.Addr
			if strings.Contains(addrStr, ".") && !strings.Contains(addrStr, ":") {
				ip4 = addrStr
			} else if strings.Contains(addrStr, ":") {
				ip6 = addrStr
			}
		}

		ifs = append(ifs, NetInfo{
			Name:         iface.Name,
			MTU:          iface.MTU,
			HardwareAddr: iface.HardwareAddr,
			Flags:        strings.Join(iface.Flags, ","),
			IP4:          ip4,
			IP6:          ip6,
		})
	}

	return &NetUsage{
		Interfaces: ifs,
		TotalRx:    totalRx,
		TotalTx:    totalTx,
	}, nil
}

// Sensors returns sensor data for macOS using ioreg
func (c *DarwinCollector) Sensors() (*SensorData, error) {
	data := &SensorData{
		Temperatures: []SensorInfo{},
		Fans:         []SensorInfo{},
		Voltages:     []SensorInfo{},
		Power:        []SensorInfo{},
		HasSensors:   false,
	}

	// Try to get CPU temperature from ioreg
	cmd := exec.Command("ioreg", "-c", "AppleSMC", "-r", "-d", "1", "-a")
	output, err := cmd.Output()
	if err == nil {
		// Parse ioreg JSON output for temperature sensors
		var result map[string]interface{}
		if json.Unmarshal(output, &result) == nil {
			if children, ok := result["Children"].([]interface{}); ok {
				for _, child := range children {
					if c, ok := child.(map[string]interface{}); ok {
						if c["key"] != nil {
							key := c["key"].(string)
							if strings.Contains(key, "TC") || strings.Contains(key, "Tp") {
								if val, ok := c["CurrentValue"]; ok {
									if fval, ok := val.(float64); ok {
										tempC := fval / 256.0
										sensorName := strings.Replace(key, "0x", "", -1)
										data.Temperatures = append(data.Temperatures, SensorInfo{
											Name:       sensorName,
											Value:      tempC,
											Unit:       "°C",
											Critical:   tempC > 95,
											SensorType: "temperature",
										})
										data.HasSensors = true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Try to get fan speed from ioreg
	cmd = exec.Command("ioreg", "-c", "AppleFan", "-r", "-d", "1", "-a")
	output, err = cmd.Output()
	if err == nil {
		var result []interface{}
		if json.Unmarshal(output, &result) == nil {
			for _, item := range result {
				if m, ok := item.(map[string]interface{}); ok {
					if fanNum, ok := m["FanNumber"].(float64); ok {
						if current, ok := m["CurrentValue"].(float64); ok {
							data.Fans = append(data.Fans, SensorInfo{
								Name:       fmt.Sprintf("Fan %.0f", fanNum),
								Value:      current,
								Unit:       "RPM",
								SensorType: "fan",
							})
							data.HasSensors = true
						}
					}
				}
			}
		}
	}

	return data, nil
}

// Power returns power/battery information for macOS using ioreg
func (c *DarwinCollector) Power() (*PowerInfo, error) {
	info := &PowerInfo{
		HasBattery:    false,
		Charging:      false,
		Percent:       0,
		TimeRemaining: 0,
		PowerPlugged:  true,
		Amps:          0,
		Volts:         0,
		Watts:         0,
		CycleCount:    0,
		Health:        "Good",
	}

	// Get battery info from ioreg
	cmd := exec.Command("ioreg", "-c", "AppleSmartBattery", "-r", "-a")
	output, err := cmd.Output()
	if err == nil {
		var result []interface{}
		if json.Unmarshal(output, &result) == nil {
			for _, item := range result {
				if m, ok := item.(map[string]interface{}); ok {
					// Check if battery exists
					if isPresent, ok := m["IsPresent"].(float64); ok && isPresent == 1 {
						info.HasBattery = true

						// Battery percentage
						if percent, ok := m["CurrentCapacity"].(float64); ok {
							if max, ok := m["MaxCapacity"].(float64); ok && max > 0 {
								info.Percent = int(percent / max * 100)
							}
						}

						// Charging status
						if charging, ok := m["IsCharging"].(float64); ok && charging == 1 {
							info.Charging = true
							info.PowerPlugged = true
						}

						// Power source
						if powerSource, ok := m["PowerSourceState"].(string); ok {
							if powerSource == "AC Power" {
								info.PowerPlugged = true
							} else {
								info.PowerPlugged = false
							}
						}

						// Time remaining (in minutes)
						if timeToEmpty, ok := m["TimeToEmpty"].(float64); ok && timeToEmpty > 0 {
							info.TimeRemaining = int(timeToEmpty)
						} else if timeToFull, ok := m["TimeToFullCharge"].(float64); ok && timeToFull > 0 {
							info.TimeRemaining = int(timeToFull)
						}

						// Cycle count
						if cycleCount, ok := m["CycleCount"].(float64); ok {
							info.CycleCount = int(cycleCount)
						}

						// Battery health
						if health, ok := m["BatteryHealth"].(string); ok {
							info.Health = health
						} else {
							// Estimate health based on cycle count
							if info.CycleCount > 1000 {
								info.Health = "Fair"
							} else if info.CycleCount > 500 {
								info.Health = "Good"
							} else {
								info.Health = "Excellent"
							}
						}

						// Voltage
						if voltage, ok := m["Voltage"].(float64); ok {
							info.Volts = voltage / 1000.0 // mV to V
						}

						// Amperage
						if amperage, ok := m["Amperage"].(float64); ok {
							info.Amps = amperage / 1000.0 // mA to A
							info.Watts = info.Volts * info.Amps
						}

						// Design capacity
						if designCap, ok := m["DesignCapacity"].(float64); ok {
							if maxCap, ok := m["MaxCapacity"].(float64); ok && designCap > 0 {
								healthPct := maxCap / designCap * 100
								if healthPct < 50 {
									info.Health = "Replace Soon"
								} else if healthPct < 80 {
									info.Health = "Fair"
								}
							}
						}
					}
				}
			}
		}
	}

	return info, nil
}
