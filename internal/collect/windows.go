package collect

import (
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// WindowsCollector implements Collector for Windows
type WindowsCollector struct{}

// NewWindowsCollector creates a Windows collector
func NewWindowsCollector() *WindowsCollector {
	return &WindowsCollector{}
}

// Platform returns platform information
func (c *WindowsCollector) Platform() (*PlatformInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	bootTime := uint64(hostInfo.BootTime)
	uptime := uint64(time.Now().Unix()) - bootTime

	return &PlatformInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: hostInfo.Hostname,
		Platform: hostInfo.Platform,
		Kernel:   hostInfo.KernelVersion,
		Uptime:   uptime,
		BootTime: bootTime,
		Procs:    int(hostInfo.Procs),
	}, nil
}

// CPU returns CPU information
func (c *WindowsCollector) CPU() (*CPUInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var model string
	var cores, threads int
	if len(cpuInfo) > 0 {
		model = cpuInfo[0].ModelName
		cores = int(cpuInfo[0].Cores)
		threads = int(cpuInfo[0].Cores)
	}

	return &CPUInfo{
		Model:   model,
		Cores:   cores,
		Threads: threads,
	}, nil
}

// CPUPercent returns per-CPU usage
func (c *WindowsCollector) CPUPercent() (*CPUPercent, error) {
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
func (c *WindowsCollector) Mem() (*MemInfo, error) {
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
	}, nil
}

// Disk returns disk usage information
func (c *WindowsCollector) Disk() (*DiskUsage, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []DiskInfo
	for _, p := range parts {
		if p.Mountpoint == "" {
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

	return &DiskUsage{Disks: disks}, nil
}

// Net returns network interface information
func (c *WindowsCollector) Net() (*NetUsage, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	stats, _ := net.IOCounters(false)

	var ifs []NetInfo
	var totalRx, totalTx uint64

	for _, iface := range ifaces {
		var ifaceStats net.IOCountersStat
		for _, s := range stats {
			if s.Name == iface.Name {
				ifaceStats = s
				break
			}
		}

		totalRx += ifaceStats.BytesRecv
		totalTx += ifaceStats.BytesSent

		var ip4, ip6 string
		for _, addr := range iface.Addrs {
			addrStr := addr.Addr
			if len(addrStr) > 0 {
				if contains(addrStr, ".") {
					ip4 = addrStr
				} else if contains(addrStr, ":") {
					ip6 = addrStr
				}
			}
		}

		ifs = append(ifs, NetInfo{
			Name:         iface.Name,
			MTU:          iface.MTU,
			HardwareAddr: iface.HardwareAddr,
			IP4:          ip4,
			IP6:          ip6,
			RxBytes:      ifaceStats.BytesRecv,
			TxBytes:      ifaceStats.BytesSent,
		})
	}

	return &NetUsage{
		Interfaces: ifs,
		TotalRx:    totalRx,
		TotalTx:    totalTx,
	}, nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

// Sensors returns sensor information for Windows
func (c *WindowsCollector) Sensors() (*SensorData, error) {
	data := &SensorData{
		Temperatures: []SensorInfo{},
		Fans:         []SensorInfo{},
		Voltages:     []SensorInfo{},
		Power:        []SensorInfo{},
		HasSensors:   false,
	}

	// Windows doesn't have standardized sensor APIs without additional libraries
	// Common approaches:
	// 1. OpenHardwareMonitor WMI
	// 2. LibreHardwareMonitor
	// 3. SpeedFan
	
	// For now, we'll try to read from registry or common locations
	// Note: Full sensor support requires external libraries like go-hardware-monitor
	
	return data, nil
}

// Power returns battery/power information for Windows
func (c *WindowsCollector) Power() (*PowerInfo, error) {
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
		Health:        "Unknown",
	}

	// Try to get battery status using systeminfo or WMI
	// Note: Full battery support requires golang.org/x/sys/windows
	
	// Check memory as a fallback for power status
	_, err := mem.VirtualMemory()
	if err == nil {
		// If we have memory info, system is running
		info.PowerPlugged = true
	}
	
	// Try to use systeminfo command as fallback
	cmd := exec.Command("systeminfo")
	output, err := cmd.Output()
	if err == nil {
		outputStr := string(output)
		
		// Check for battery status
		if strings.Contains(outputStr, "Battery") {
			info.HasBattery = true
			
			// Parse battery percentage from systeminfo output
			lines := strings.Split(outputStr, "\n")
			for _, line := range lines {
				if strings.Contains(line, "Battery") && strings.Contains(line, "%") {
					// Try to extract percentage
					re := regexp.MustCompile(`(\d+)%`)
					matches := re.FindStringSubmatch(line)
					if len(matches) > 1 {
						info.Percent, _ = strconv.Atoi(matches[1])
					}
				}
			}
		}
	}
	
	// If no battery detected, it's likely a desktop
	if !info.HasBattery {
		info.HasBattery = false
		info.Percent = 0
		info.PowerPlugged = true
		info.Health = "N/A"
	}

	return info, nil
}

var _ Collector = (*WindowsCollector)(nil)
