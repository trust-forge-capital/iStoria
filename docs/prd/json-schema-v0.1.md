# iStoria JSON 输出字段规范 v0.1

## 目标
- 为 `--json` 输出提供稳定、可脚本解析的字段约定。
- 保证核心字段跨平台尽量统一。
- 允许平台特有字段通过扩展对象暴露，避免破坏兼容性。

## 设计原则
- 顶层字段稳定。
- 字段名统一使用 `snake_case`。
- 时间字段使用 ISO 8601 UTC 字符串。
- 百分比字段统一使用 `percent` 或 `used_percent` 这类明确命名。
- 平台增强字段放在 `extra` 或模块级 `extra` 内。
- 缺失字段优先省略，不强制输出 `null`。

## 顶层结构

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
  "cpu": {},
  "memory": {},
  "disk": {},
  "network": {},
  "sensor": {},
  "power": {},
  "extra": {}
}
```

## 通用顶层字段
- `timestamp`: 采集时间，UTC ISO 8601
- `hostname`: 主机名
- `platform.os`: `darwin` / `linux` / `windows`
- `platform.arch`: `arm64` / `amd64` 等
- `platform.version`: 系统版本字符串
- `platform.family`: 展示友好的平台名称
- `command`: 当前命令，如 `stat` / `cpu` / `mem`

## `stat` 命令
`stat` 是聚合摘要，建议包含：

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

## `cpu` 模块

```json
{
  "usage_percent": 18.4,
  "cores_logical": 10,
  "cores_physical": 8,
  "load_avg": [1.02, 0.98, 0.77],
  "frequency_mhz": 3200,
  "per_core_usage_percent": [12.1, 9.8, 22.5],
  "extra": {
    "performance_core_usage_percent": [31.2, 29.8],
    "efficiency_core_usage_percent": [8.2, 6.9]
  }
}
```

字段说明：
- `usage_percent`: CPU 总占用
- `cores_logical`: 逻辑核心数
- `cores_physical`: 物理核心数（可获取时）
- `load_avg`: 1/5/15 分钟负载（平台支持时）
- `frequency_mhz`: 当前或摘要频率（可获取时）
- `per_core_usage_percent`: 每核心占用
- `extra`: 平台增强字段，例如 Apple Silicon 的性能核/能效核分组

## `memory` 模块

```json
{
  "total_bytes": 17179869184,
  "used_bytes": 9437184000,
  "free_bytes": 2147483648,
  "available_bytes": 7730941132,
  "used_percent": 54.9,
  "swap_total_bytes": 2147483648,
  "swap_used_bytes": 0,
  "extra": {
    "pressure": "normal"
  }
}
```

## `disk` 模块

```json
{
  "root": {
    "path": "/",
    "filesystem": "apfs",
    "total_bytes": 494384795648,
    "used_bytes": 224197754880,
    "free_bytes": 270187040768,
    "used_percent": 45.3
  },
  "volumes": [
    {
      "path": "/",
      "filesystem": "apfs",
      "total_bytes": 494384795648,
      "used_bytes": 224197754880,
      "free_bytes": 270187040768,
      "used_percent": 45.3
    }
  ]
}
```

## `network` 模块

```json
{
  "primary_interface": "en0",
  "interfaces": [
    {
      "name": "en0",
      "ipv4": ["192.168.1.10"],
      "ipv6": ["fe80::1"],
      "rx_bytes_per_sec": 1024,
      "tx_bytes_per_sec": 2048
    }
  ]
}
```

## `sensor` 模块

```json
{
  "temperature_celsius": 52.1,
  "temperatures": [
    {
      "name": "cpu",
      "celsius": 52.1
    }
  ],
  "fan_rpm": [1200],
  "power_watts": 18.5,
  "extra": {
    "apple_silicon": {
      "performance_cores_frequency_mhz": 3228,
      "efficiency_cores_frequency_mhz": 2064
    }
  }
}
```

## `power` 模块

```json
{
  "source": "ac",
  "battery_percent": 100,
  "charging": false,
  "cycle_count": 87,
  "health_percent": 96
}
```

## 平台扩展策略
- 顶层和核心模块字段尽量固定。
- Apple Silicon / 特定 Linux 发行版 / Windows 专属字段放到模块下的 `extra`。
- 不同平台允许子字段不同，但不要改变已存在字段语义。

## 兼容性规则
- 新增字段：允许。
- 删除字段：尽量避免。
- 修改字段语义：禁止。
- 修改字段类型：禁止。

## v0.1 强制统一字段
- `timestamp`
- `hostname`
- `platform.os`
- `platform.arch`
- `command`
- `cpu.usage_percent`
- `memory.total_bytes`
- `memory.used_bytes`
- `memory.used_percent`
- `disk.root.path`
- `disk.root.used_percent`

## 后续版本可扩展方向
- 告警 JSON 结构
- 规则定义 JSON 结构
- 历史样本输出格式
- 流式输出/持续监控模式
