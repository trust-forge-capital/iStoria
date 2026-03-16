# iStoria 命令输出样例 v0.1

> 目的：统一产品与开发对命令输出体验的预期。

## 1. `istoria stat`

### 终端输出示例

```text
Host:      mac-mini
Platform:  macOS 15.0 (darwin/arm64)
Uptime:    3d 4h 12m

CPU:       18.4%   load 1.02 0.98 0.77
Memory:    9.4 GB / 16.0 GB (54.9%)
Disk(/):   224.2 GB / 494.4 GB (45.3%)
Network:   en0  rx 1.0 KB/s  tx 2.0 KB/s
Sensor:    temp 52.1°C  fan 1200 rpm
Power:     AC  battery 100%
```

### JSON 输出示例

```json
{
  "timestamp": "2026-03-14T06:44:00Z",
  "hostname": "mac-mini",
  "platform": {
    "os": "darwin",
    "arch": "arm64",
    "version": "15.0",
    "family": "macOS"
  },
  "command": "stat",
  "cpu": {
    "usage_percent": 18.4,
    "cores_logical": 10,
    "load_avg": [1.02, 0.98, 0.77]
  },
  "memory": {
    "total_bytes": 17179869184,
    "used_bytes": 9437184000,
    "used_percent": 54.9,
    "swap_used_bytes": 0
  },
  "disk": {
    "root": {
      "path": "/",
      "total_bytes": 494384795648,
      "used_bytes": 224197754880,
      "used_percent": 45.3
    }
  },
  "network": {
    "primary_interface": "en0",
    "rx_bytes_per_sec": 1024,
    "tx_bytes_per_sec": 2048,
    "ipv4": ["192.168.1.10"]
  },
  "sensor": {
    "temperature_celsius": 52.1,
    "fan_rpm": [1200]
  },
  "power": {
    "source": "ac",
    "battery_percent": 100,
    "charging": false
  }
}
```

## 2. `istoria cpu`

### 终端输出示例

```text
CPU Usage:     18.4%
Logical Cores: 10
Physical Cores: 8
Load Avg:      1.02 0.98 0.77
Frequency:     3200 MHz

Per Core:
  core0  12.1%
  core1   9.8%
  core2  22.5%
  core3  14.7%
```

### Apple Silicon 增强展示示例

```text
Clusters:
  performance  31.2% 29.8%
  efficiency    8.2%  6.9%  7.1%  5.4%
```

## 3. `istoria mem`

```text
Memory Total:   16.0 GB
Memory Used:     9.4 GB (54.9%)
Memory Free:     2.0 GB
Available:       7.2 GB
Swap Used:       0 B
Pressure:        normal
```

## 4. `istoria disk`

```text
Mount   FS    Total      Used       Free       Used%
/       apfs  494.4 GB   224.2 GB   270.2 GB   45.3%
```

## 5. `istoria sensor`

```text
Temperature:
  CPU        52.1°C
  GPU        49.3°C

Fans:
  fan0       1200 rpm

Power:
  Package    18.5 W
```

### Apple Silicon 增强展示示例

```text
Apple Silicon:
  perf_freq_mhz   3228
  eff_freq_mhz    2064
```

## 6. `istoria net`

```text
Primary Interface: en0
IPv4: 192.168.1.10
RX:   1.0 KB/s
TX:   2.0 KB/s
```

## 7. `istoria power`

```text
Source:         AC
Battery:        100%
Charging:       false
Cycle Count:    87
Health:         96%
```

## 8. `istoria stat --live`

### 终端输出示例

```text
[refresh: 1s | ctrl+c to stop]
Host:      mac-mini
Platform:  macOS 15.0 (darwin/arm64)
Uptime:    3d 4h 12m

CPU:       21.8%   load 1.12 1.01 0.80
Memory:    9.6 GB / 16.0 GB (56.0%)
Disk(/):   224.2 GB / 494.4 GB (45.3%)
Network:   en0  rx 18.2 KB/s  tx 42.7 KB/s
Sensor:    temp 54.3°C  fan 1298 rpm
Power:     AC  battery 100%
```

建议参数：
- `--live`：开启实时刷新
- `--interval 1s`：刷新间隔，默认 1 秒
- `--no-clear`：调试模式，不清屏，仅追加输出

## 9. `istoria net --live`

### 终端输出示例

```text
[refresh: 1s | interface: en0 | ctrl+c to stop]
RX:   18.2 KB/s
TX:   42.7 KB/s

Total RX:  1.2 GB
Total TX:  842.1 MB
IPv4:      192.168.1.10
```

说明：
- 这是最接近 `vnstat -l -i eth0` 的模式。
- 第一版建议优先实现 `net --live`，其次实现 `stat --live`。

## 输出约束
- 默认输出尽量一屏内读完。
- `stat` 优先摘要，不展开所有细节。
- 详细信息放在子命令中。
- live 模式优先纯文本刷新，不引入复杂 TUI 交互。
- JSON 字段以 `json-schema-v0.1.md` 为准。
- 当平台不支持某字段时，默认输出不强行展示占位值。
