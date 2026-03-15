# iStoria

适合 SSH 场景的跨平台硬件监控 CLI。

## 一句话介绍

iStoria 帮助用户在登录目标机器后，用统一命令快速查看本机硬件状态；默认输出适合人读，`--json` 适合脚本调用。

## 产品定位

- **核心场景**：SSH 登录目标机器后，在本机执行命令查看状态
- **目标平台**：macOS / Linux / Windows
- **分发方式**：单文件可执行程序
- **核心差异化**：
  - 统一入口
  - 终端友好输出
  - JSON 自动化友好
  - 持续增强 Apple Silicon 支持

## 安装

### macOS / Linux

```bash
# 下载二进制
curl -L https://github.com/yourorg/istoria/releases/latest/download/istoria -o istoria
chmod +x istoria

# 或使用 Homebrew
brew install istoria
```

### Windows

```powershell
# 使用 winget
winget install istoria

# 或下载 istoria.exe
```

## 快速开始

```bash
# 查看系统摘要
istoria stat

# 查看 CPU 信息
istoria cpu

# 查看内存信息
istoria mem

# 查看磁盘信息
istoria disk

# 查看传感器信息（需要 istats 或 lm-sensors）
istoria sensor

# 查看网络信息
istoria net

# 查看电源/电池信息
istoria power
```

## 使用示例

### 基本用法

```bash
# 查看 CPU 信息
$ istoria cpu
=== CPU Information ===
Model: Apple M4
Physical Cores: 10
Logical Threads: 20

--- Apple Silicon ---
Performance Cores: 4
Efficiency Cores: 6

--- Usage ---
Total:      15.2%
User:       9.1%
System:     6.1%
Idle:      84.8%
```

### JSON 输出

```bash
# 机器可读格式
$ istoria cpu --json
{
  "model": "Apple M4",
  "cores": 10,
  "threads": 20,
  "usage_percent": 15.2
}
```

### Live 实时监控模式

```bash
# 实时监控 CPU（每 500ms 刷新）
istoria cpu --live --interval 500

# 实时监控内存（每 1s 刷新）
istoria mem --live

# 不清屏模式（追加输出）
istoria cpu --live --no-clear
```

### 其他选项

```bash
# 禁用颜色输出
istoria cpu --no-color

# 静默模式（减少非必要输出）
istoria cpu --quiet

# 指定配置文件
istoria --config /path/to/config.yaml cpu

# 查看版本
istoria version
```

## 命令列表

| 命令 | 描述 |
|------|------|
| `istoria stat` | 系统摘要信息 |
| `istoria cpu` | CPU 信息 |
| `istoria mem` | 内存信息 |
| `istoria disk` | 磁盘信息 |
| `istoria sensor` | 传感器信息 |
| `istoria net` | 网络信息 |
| `istoria power` | 电源/电池信息 |
| `istoria version` | 版本信息 |

## 全局参数

| 参数 | 描述 |
|------|------|
| `--json` | JSON 格式输出 |
| `--no-color` | 禁用颜色 |
| `--quiet` | 静默模式 |
| `--live` | 实时刷新模式 |
| `--interval` | 刷新间隔（毫秒，默认 1000） |
| `--no-clear` | 不清屏（实时模式下） |
| `--config` | 配置文件路径 |

## Apple Silicon 支持

iStoria 针对 Apple Silicon Mac 提供了增强的监控能力：

- 性能核心 / 能效核心分离显示
- Apple 芯片专用传感器
- 优化的数据采集

## 设计原则

- 命令短、直观、适合 SSH 使用
- 默认输出优先给人看
- JSON 输出优先给脚本用
- 一屏内尽量看完主要信息
- 平台差异允许存在，但核心字段要统一

## 技术栈

- Go 1.24+
- Cobra（命令行框架）
- gopsutil（系统信息采集）

## 项目结构

```
istoria/
├── README.md
├── CHANGELOG.md
├── go.mod
├── main.go
├── cmd/              # CLI 命令
├── internal/         # 内部包
│   └── collect/      # 数据采集层
├── configs/          # 配置文件
├── scripts/          # 辅助脚本
└── docs/            # 文档
```

## 文档

- 产品需求文档：`docs/prd/istoria.md`
- 命令输出示例：`docs/prd/command-output-examples-v0.1.md`

## 当前状态

- ✅ v0.1 已发布
- ✅ 单元测试覆盖
- ✅ 三平台支持 (macOS/Linux/Windows)
- ✅ Live 实时监控模式

## 许可证

MIT License
