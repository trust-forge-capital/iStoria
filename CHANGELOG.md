# Changelog

All notable changes to this project will be documented in this file.

## [v0.1.1] - 2026-03-16

### Added
- **Live 实时监控模式**：新增 `--live` 参数支持实时刷新监控
  - `--interval`：刷新间隔（毫秒，默认 1000，最小 500）
  - `--no-clear`：不清屏模式
  - 支持 `cpu`、`mem`、`disk`、`sensor`、`net` 命令的实时监控
- **单元测试覆盖**：增加测试文件，提升代码质量
  - 新增 `cmd/live_test.go`
  - 新增 `cmd/root_test.go`
  - 新增 `internal/collect/data_test.go`
- **README 优化**：补充完整的 usage 示例和安装说明

### Improved
- `--help` 输出优化
- 代码注释完善

### Fixed
- 修复 macOS 平台编译问题

## [v0.1.0] - 2026-03-15

### Added
- 首次发布：iStoria v0.1
- 跨平台硬件监控 CLI (macOS / Linux / Windows)
- 统一入口查看 CPU/内存/磁盘/网络状态
- 终端友好输出 + JSON 自动化友好输出
- Apple Silicon 支持

### Commands
- `istoria stat` - 系统摘要
- `istoria cpu` - CPU 信息
- `istoria mem` - 内存信息
- `istoria disk` - 磁盘信息
- `istoria sensor` - 传感器信息
- `istoria net` - 网络信息
- `istoria power` - 电源/电池信息
- `istoria version` - 版本信息
