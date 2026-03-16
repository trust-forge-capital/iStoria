package collect

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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

// LinuxCollector implements Collector for Linux
type LinuxCollector struct{}

// NewLinuxCollector creates a Linux collector
func NewLinuxCollector() *LinuxCollector {
	return &LinuxCollector{}
}

// Platform returns platform information
func (c *LinuxCollector) Platform() (*PlatformInfo, error) {
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
func (c *LinuxCollector) CPU() (*CPUInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var model string
	var cores, threads int
	var freq, freqMin, freqMax uint64

	if len(cpuInfo) > 0 {
		model = cpuInfo[0].ModelName
		// Use Cores from cpu.Info
		cores = int(cpuInfo[0].Cores)
		freq = uint64(cpuInfo[0].Mhz * 1_000_000)

		// Get logical CPU count from cpu.Counts
		logicalCPUs, _ := cpu.Counts(true)
		if logicalCPUs > 0 {
			threads = logicalCPUs
		}

		// If cores/threads are 0 or 1, try to get from /proc/cpuinfo
		if cores <= 1 || threads <= 1 {
			procCores, procThreads := c.getCPUCoresFromProc()
			if procCores > cores {
				cores = procCores
			}
			if procThreads > threads {
				threads = procThreads
			}
		}
	}

	// If still invalid, try sysfs
	if cores <= 1 {
		if cpus, err := c.getCPUCountFromSysfs(); err == nil && cpus > 0 {
			cores = cpus
			threads = cpus * 2 // Assume hyperthreading
		}
	}

	if freq == 0 {
		freq, freqMin, freqMax = c.getCPUFreqFromProc()
	}

	return &CPUInfo{
		Model:        model,
		Cores:        cores,
		Threads:      threads,
		Frequency:    freq,
		FrequencyMin: freqMin,
		FrequencyMax: freqMax,
	}, nil
}

// getCPUCoresFromProc reads CPU core info from /proc/cpuinfo
func (c *LinuxCollector) getCPUCoresFromProc() (cores int, threads int) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return 1, 1
	}

	lines := strings.Split(string(data), "\n")

	var physicalIDs []int
	var cpuCores []int
	var siblings []int

	for _, line := range lines {
		if strings.HasPrefix(line, "processor") {
			threads++ // Each processor is a logical thread
		}
		if strings.HasPrefix(line, "physical id") {
			id := strings.TrimSpace(strings.Split(line, ":")[1])
			pid, _ := strconv.Atoi(id)
			physicalIDs = append(physicalIDs, pid)
		}
		if strings.HasPrefix(line, "cpu cores") {
			val := strings.TrimSpace(strings.Split(line, ":")[1])
			c, _ := strconv.Atoi(val)
			cpuCores = append(cpuCores, c)
		}
		if strings.HasPrefix(line, "siblings") {
			val := strings.TrimSpace(strings.Split(line, ":")[1])
			s, _ := strconv.Atoi(val)
			siblings = append(siblings, s)
		}
	}

	// Calculate unique physical CPUs
	uniquePhysical := make(map[int]bool)
	for _, pid := range physicalIDs {
		uniquePhysical[pid] = true
	}

	if len(uniquePhysical) > 0 && len(cpuCores) > 0 {
		cores = len(uniquePhysical) * cpuCores[0]
	}

	if threads <= 0 {
		threads = cores * 2
	}

	return cores, threads
}

// getCPUCountFromSysfs gets CPU count from sysfs
func (c *LinuxCollector) getCPUCountFromSysfs() (int, error) {
	data, err := os.ReadFile("/sys/devices/system/cpu/online")
	if err != nil {
		return 0, err
	}

	// Parse "0-11" or "0,1,2,3" format
	parts := strings.Split(strings.TrimSpace(string(data)), ",")
	var count int
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			start, _ := strconv.Atoi(rangeParts[0])
			end, _ := strconv.Atoi(rangeParts[1])
			count += end - start + 1
		} else {
			count++
		}
	}

	return count, nil
}

func (c *LinuxCollector) getCPUFreqFromProc() (uint64, uint64, uint64) {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err == nil {
		var cur uint64
		fmt.Sscanf(string(data), "%d", &cur)
		cur *= 1000

		minData, _ := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_min_freq")
		var min uint64
		fmt.Sscanf(string(minData), "%d", &min)
		min *= 1000

		maxData, _ := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq")
		var max uint64
		fmt.Sscanf(string(maxData), "%d", &max)
		max *= 1000

		return cur, min, max
	}
	return 0, 0, 0
}

// CPUPercent returns per-CPU usage
func (c *LinuxCollector) CPUPercent() (*CPUPercent, error) {
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
func (c *LinuxCollector) Mem() (*MemInfo, error) {
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
func (c *LinuxCollector) Disk() (*DiskUsage, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []DiskInfo
	for _, p := range parts {
		if p.Mountpoint == "" || p.Mountpoint == "/dev" {
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
			InodesTotal: usage.Total,
			InodesUsed:  usage.Used,
			InodesAvail: usage.Free,
		})
	}

	return &DiskUsage{Disks: disks}, nil
}

// Net returns network interface information
func (c *LinuxCollector) Net() (*NetUsage, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	stats, _ := net.IOCounters(false)

	var ifs []NetInfo
	var totalRx, totalTx uint64

	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "lo") {
			continue
		}

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
			if strings.Contains(addrStr, ".") {
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
			RxBytes:      ifaceStats.BytesRecv,
			TxBytes:      ifaceStats.BytesSent,
			RxPackets:    ifaceStats.PacketsRecv,
			TxPackets:    ifaceStats.PacketsSent,
		})
	}

	return &NetUsage{
		Interfaces: ifs,
		TotalRx:    totalRx,
		TotalTx:    totalTx,
	}, nil
}

// Sensors returns sensor information
func (c *LinuxCollector) Sensors() (*SensorData, error) {
	data := &SensorData{
		Temperatures: []SensorInfo{},
		Fans:         []SensorInfo{},
		Voltages:     []SensorInfo{},
		Power:        []SensorInfo{},
		HasSensors:   false,
	}

	cmd := exec.Command("sensors")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		lines := strings.Split(out.String(), "\n")
		var currentZone string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasSuffix(line, ":") {
				currentZone = strings.TrimSuffix(line, ":")
				continue
			}
			if strings.Contains(line, "+") && strings.Contains(line, "°C") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					tempStr := strings.Trim(fields[1], "+°C")
					var temp float64
					fmt.Sscanf(tempStr, "%f", &temp)

					data.Temperatures = append(data.Temperatures, SensorInfo{
						Name:       currentZone + " " + fields[0],
						Value:      temp,
						Unit:       "°C",
						SensorType: "temperature",
					})
					data.HasSensors = true
				}
			}
			if strings.Contains(line, "RPM") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					var rpm float64
					fmt.Sscanf(fields[1], "%f", &rpm)
					data.Fans = append(data.Fans, SensorInfo{
						Name:       fields[0],
						Value:      rpm,
						Unit:       "RPM",
						SensorType: "fan",
					})
					data.HasSensors = true
				}
			}
		}
	}

	if !data.HasSensors {
		c.readHWMON(data)
	}

	return data, nil
}

func (c *LinuxCollector) readHWMON(data *SensorData) {
	paths := []string{
		"/sys/class/hwmon/hwmon0",
		"/sys/class/hwmon/hwmon1",
		"/sys/class/hwmon/hwmon2",
	}

	for _, path := range paths {
		nameData, err := os.ReadFile(path + "/name")
		if err != nil {
			continue
		}

		hwmonName := strings.TrimSpace(string(nameData))

		for i := 0; i < 10; i++ {
			tempInput := fmt.Sprintf("%s/temp%d_input", path, i+1)
			tempData, err := os.ReadFile(tempInput)
			if err == nil {
				var temp int64
				fmt.Sscanf(string(tempData), "%d", &temp)
				tempC := float64(temp) / 1000.0
				if tempC > 0 && tempC < 150 {
					data.Temperatures = append(data.Temperatures, SensorInfo{
						Name:       fmt.Sprintf("%s temp%d", hwmonName, i+1),
						Value:      tempC,
						Unit:       "°C",
						SensorType: "temperature",
					})
					data.HasSensors = true
				}
			}
		}

		for i := 0; i < 10; i++ {
			fanInput := fmt.Sprintf("%s/fan%d_input", path, i+1)
			fanData, err := os.ReadFile(fanInput)
			if err == nil {
				var rpm int64
				fmt.Sscanf(string(fanData), "%d", &rpm)
				if rpm > 0 {
					data.Fans = append(data.Fans, SensorInfo{
						Name:       fmt.Sprintf("%s fan%d", hwmonName, i+1),
						Value:      float64(rpm),
						Unit:       "RPM",
						SensorType: "fan",
					})
					data.HasSensors = true
				}
			}
		}
	}
}

// Power returns battery/power information
func (c *LinuxCollector) Power() (*PowerInfo, error) {
	p := &PowerInfo{HasBattery: false}

	batPaths := []string{
		"/sys/class/power_supply/BAT0",
		"/sys/class/power_supply/BAT1",
	}

	for _, path := range batPaths {
		presentFile := path + "/present"
		if _, err := os.ReadFile(presentFile); err == nil {
			p.HasBattery = true

			capData, _ := os.ReadFile(path + "/capacity")
			fmt.Sscanf(string(capData), "%d", &p.Percent)

			statusData, _ := os.ReadFile(path + "/status")
			status := strings.TrimSpace(string(statusData))
			p.Charging = status == "Charging"
			p.PowerPlugged = status == "Full" || p.Charging

			voltData, _ := os.ReadFile(path + "/voltage_now")
			var volt int64
			fmt.Sscanf(string(voltData), "%d", &volt)
			p.Volts = float64(volt) / 1_000_000.0

			curData, _ := os.ReadFile(path + "/current_now")
			var cur int64
			fmt.Sscanf(string(curData), "%d", &cur)
			p.Amps = float64(-cur) / 1_000_000.0
			p.Watts = p.Volts * p.Amps

			break
		}
	}

	return p, nil
}

var _ Collector = (*LinuxCollector)(nil)
