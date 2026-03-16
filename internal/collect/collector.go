package collect

import "runtime"

// PlatformInfo represents basic platform information
type PlatformInfo struct {
	OS           string `json:"os"`            // "darwin", "linux", "windows"
	Arch         string `json:"arch"`          // "arm64", "amd64", etc.
	Hostname     string `json:"hostname"`      // machine hostname
	Platform     string `json:"platform"`      // platform-specific string
	Kernel       string `json:"kernel"`        // kernel version
	Uptime       uint64 `json:"uptime"`        // uptime in seconds
	BootTime     uint64 `json:"boot_time"`     // boot time Unix timestamp
	Procs        int    `json:"procs"`         // number of processors
	AppleSilicon bool   `json:"apple_silicon"` // true if Apple Silicon
}

// CPUInfo represents CPU information
type CPUInfo struct {
	Model            string  `json:"model"`             // CPU model name
	Cores            int     `json:"cores"`             // physical cores
	Threads          int     `json:"threads"`           // logical processors
	Frequency        uint64  `json:"frequency"`         // CPU frequency in Hz
	FrequencyMin     uint64  `json:"frequency_min"`     // min frequency
	FrequencyMax     uint64  `json:"frequency_max"`     // max frequency
	UsagePercent     float64 `json:"usage_percent"`     // total CPU usage
	UserPercent      float64 `json:"user_percent"`      // user CPU usage
	SystemPercent    float64 `json:"system_percent"`    // system CPU usage
	IdlePercent      float64 `json:"idle_percent"`      // idle CPU usage
	IowaitPercent    float64 `json:"iowait_percent"`    // iowait (Linux)
	NicePercent      float64 `json:"nice_percent"`      // nice CPU usage
	AppleSilicon     bool    `json:"apple_silicon"`     // Apple Silicon flag
	PerformanceCores int     `json:"performance_cores"` // P-cores (Apple Silicon)
	EfficiencyCores  int     `json:"efficiency_cores"`  // E-cores (Apple Silicon)
}

// CPUPercent returns per-CPU usage
type CPUPercent struct {
	PerCPU []float64 `json:"per_core_usage_percent"` // per-CPU usage percentages
	Total  float64   `json:"usage_percent"`          // total CPU usage
}

// MemInfo represents memory information
type MemInfo struct {
	Total       uint64  `json:"total_bytes"`      // total memory in bytes
	Available   uint64  `json:"available_bytes"`  // available memory
	Used        uint64  `json:"used_bytes"`       // used memory
	Free        uint64  `json:"free_bytes"`       // free memory
	UsedPercent float64 `json:"used_percent"`     // usage percentage
	SwapTotal   uint64  `json:"swap_total_bytes"` // swap total
	SwapUsed    uint64  `json:"swap_used_bytes"`  // swap used
	SwapFree    uint64  `json:"swap_free_bytes"`  // swap free
	Wired       uint64  `json:"wired_bytes"`      // wired memory (macOS)
	Compressed  uint64  `json:"compressed_bytes"` // compressed memory (macOS)
}

// DiskInfo represents disk information
type DiskInfo struct {
	Path        string  `json:"path"`         // mount point
	Filesystem  string  `json:"filesystem"`   // filesystem type
	Total       uint64  `json:"total_bytes"`  // total bytes
	Used        uint64  `json:"used_bytes"`   // used bytes
	Available   uint64  `json:"free_bytes"`   // available bytes
	UsedPercent float64 `json:"used_percent"` // usage percentage
	InodesTotal uint64  `json:"inodes_total"` // inodes total
	InodesUsed  uint64  `json:"inodes_used"`  // inodes used
	InodesAvail uint64  `json:"inodes_avail"` // inodes available
}

// DiskUsage returns multiple disks
type DiskUsage struct {
	Disks []DiskInfo `json:"volumes"` // list of disks
}

// NetInfo represents network interface information
type NetInfo struct {
	Name         string `json:"name"`        // interface name
	MTU          int    `json:"mtu"`         // MTU
	HardwareAddr string `json:"mac_address"` // MAC address
	Flags        string `json:"flags"`       // flags (up, down, etc)
	IP4          string `json:"ipv4"`        // IPv4 address
	IP6          string `json:"ipv6"`        // IPv6 address
	RxBytes      uint64 `json:"rx_bytes"`    // total bytes received
	TxBytes      uint64 `json:"tx_bytes"`    // total bytes transmitted
	RxPackets    uint64 `json:"rx_packets"`  // packets received
	TxPackets    uint64 `json:"tx_packets"`  // packets transmitted
	RxErrors     uint64 `json:"rx_errors"`   // receive errors
	TxErrors     uint64 `json:"tx_errors"`   // transmit errors
	RxDropped    uint64 `json:"rx_dropped"`  // dropped received
	TxDropped    uint64 `json:"tx_dropped"`  // dropped transmitted
	SpeedMbps    int    `json:"speed_mbps"`  // interface speed (Mbps)
}

// NetUsage returns network interface usage
type NetUsage struct {
	Interfaces []NetInfo `json:"interfaces"`     // list of interfaces
	TotalRx    uint64    `json:"total_rx_bytes"` // total received
	TotalTx    uint64    `json:"total_tx_bytes"` // total transmitted
}

// SensorInfo represents sensor information
type SensorInfo struct {
	Name       string  `json:"name"`        // sensor name
	Value      float64 `json:"value"`       // sensor value
	Unit       string  `json:"unit"`        // unit (°C, RPM, V, W, etc)
	Critical   bool    `json:"critical"`    // critical threshold exceeded
	Threshold  float64 `json:"threshold"`   // critical threshold
	SensorType string  `json:"sensor_type"` // "temperature", "fan", "voltage", "power"
	Location   string  `json:"location"`    // sensor location
}

// SensorData returns sensor readings
type SensorData struct {
	Temperatures []SensorInfo `json:"temperatures"` // temperature sensors
	Fans         []SensorInfo `json:"fans"`         // fan sensors
	Voltages     []SensorInfo `json:"voltages"`     // voltage sensors
	Power        []SensorInfo `json:"power"`        // power sensors
	HasSensors   bool         `json:"has_sensors"`  // whether sensors are available
}

// PowerInfo represents power/battery information
type PowerInfo struct {
	HasBattery    bool    `json:"has_battery"`     // has battery
	Charging      bool    `json:"charging"`        // is charging
	Percent       int     `json:"battery_percent"` // battery percentage
	TimeRemaining int     `json:"time_remaining"`  // minutes remaining
	PowerPlugged  bool    `json:"power_plugged"`   // plugged in
	Amps          float64 `json:"amps"`            // current amperage
	Volts         float64 `json:"volts"`           // voltage
	Watts         float64 `json:"watts"`           // watts consumed
	CycleCount    int     `json:"cycle_count"`     // battery cycle count
	Health        string  `json:"health"`          // battery health
}

// StatInfo represents the summary stat command output
type StatInfo struct {
	Platform PlatformInfo `json:"platform"`
	CPU      CPUInfo      `json:"cpu"`
	Mem      MemInfo      `json:"memory"`
	Disk     DiskInfo     `json:"disk"`     // root disk
	Uptime   string       `json:"uptime"`   // human-readable uptime
	LoadAvg  string       `json:"load_avg"` // load average (Linux)
}

// Collector is the interface for data collection
type Collector interface {
	Platform() (*PlatformInfo, error)
	CPU() (*CPUInfo, error)
	CPUPercent() (*CPUPercent, error)
	Mem() (*MemInfo, error)
	Disk() (*DiskUsage, error)
	Net() (*NetUsage, error)
	Sensors() (*SensorData, error)
	Power() (*PowerInfo, error)
}

// NewCollector creates a platform-specific collector
func NewCollector() Collector {
	switch runtime.GOOS {
	case "darwin":
		return NewDarwinCollector()
	case "linux":
		return NewLinuxCollector()
	case "windows":
		return NewWindowsCollector()
	default:
		return NewDarwinCollector() // fallback
	}
}
