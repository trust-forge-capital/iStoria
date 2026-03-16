# iStoria Live 模式交互规范 v0.1

## 目标
- 为实时刷新模式提供统一交互规范。
- 满足 SSH 场景下持续观察当前状态的需求。
- 保持实现轻量，不演变成复杂 TUI。

## 支持范围（v0.1）

### 第一优先级
- `istoria net --live`
- `istoria stat --live`

### 第二优先级
- `istoria cpu --live`
- `istoria sensor --live`

### 暂不支持
- `mem --live` 单独优先实现不是必须
- `disk --live` 单独优先实现不是必须
- 复杂交互式面板
- 鼠标操作
- 多窗口布局

## 基本参数

### `--live`
开启实时刷新模式。

### `--interval <duration>`
设置刷新间隔。

建议：
- 默认值：`1s`
- 最小值：`500ms`
- 推荐范围：`1s` ~ `3s`

说明：
- 过低刷新频率会增加资源消耗。
- 当用户传入低于最小值时，可自动提升到 `500ms`，并给出提示。

### `--no-clear`
开启后不清屏，而是持续追加输出。

适用场景：
- 调试
- 日志采样
- 终端不支持清屏控制的环境

默认行为：
- 未指定时，每次刷新前清屏重绘

## 终端行为规范

### 默认行为
- 进入 live 模式后，输出一屏摘要信息
- 每个刷新周期清屏并重绘
- 顶部显示 refresh 信息
- 使用 `Ctrl+C` 退出

### 顶部状态行建议

```text
[refresh: 1s | ctrl+c to stop]
```

对于网络命令可显示更多上下文：

```text
[refresh: 1s | interface: en0 | ctrl+c to stop]
```

### 退出方式
- 标准退出：`Ctrl+C`
- 收到中断信号后应优雅退出，不打印冗长堆栈
- 退出码建议为 `0`（用户主动中断）

## 输出策略

### `istoria net --live`
目标：尽量贴近 `vnstat -l -i eth0` 的体验。

建议展示：
- 当前 RX 速率
- 当前 TX 速率
- 累计 RX/TX
- 当前 IP
- 主接口名

不建议第一版展示过多字段，避免刷屏。

### `istoria stat --live`
目标：持续观察整机状态。

建议展示：
- CPU 摘要
- 内存摘要
- 根分区摘要
- 网络摘要
- 传感器摘要
- 电源摘要（如可用）

原则：
- 一屏内能看完
- 不在 live 模式默认展开每核心全量数据

## JSON 策略

### v0.1 建议
- `--live` 与 `--json` 可以同时使用
- 但第一版不做复杂流式 JSON 协议

建议行为：
- 每次刷新输出一行 JSON（JSON Lines / JSONL）
- 配合 `--no-clear` 使用效果最好

示例：

```bash
istoria net --live --json --no-clear
```

输出：

```json
{"timestamp":"2026-03-14T12:35:01Z","command":"net","network":{"rx_bytes_per_sec":1024,"tx_bytes_per_sec":2048}}
{"timestamp":"2026-03-14T12:35:02Z","command":"net","network":{"rx_bytes_per_sec":980,"tx_bytes_per_sec":2210}}
```

### 非推荐组合
- `--live --json` 且默认清屏

原因：
- 不利于脚本读取
- 终端体验和机器读取目标冲突

产品建议：
- 当检测到 `--live --json` 且未显式指定 `--no-clear` 时，可以自动切换为 no-clear，或给出警告提示。

## 刷新性能原则
- live 模式是“持续观察”，不是“高频采样分析器”。
- 第一版优先稳定、低开销。
- 若某些传感器采集成本高，可降低该模块内部刷新精度，或在文档中说明。

## 错误处理

### 单次采集失败
- 不要直接退出 live 模式
- 在对应区域显示：

```text
sensor: unavailable
```

或

```text
network: read error
```

### 连续失败
- 可继续重试
- 但应在顶部状态区提示错误状态

## 平台差异策略
- 不要求所有平台 live 模式字段完全一致
- 但 `net --live` 和 `stat --live` 的核心体验要尽量一致
- Apple Silicon 增强字段允许仅在支持平台展示

## 验收标准
1. [ ] `istoria net --live` 可按默认 1 秒刷新显示当前网络速率。
2. [ ] `istoria stat --live` 可按固定间隔刷新整机摘要。
3. [ ] `Ctrl+C` 可稳定退出。
4. [ ] `--interval` 可调整刷新频率。
5. [ ] `--no-clear` 可切换为追加输出模式。
6. [ ] `--live --json --no-clear` 可输出 JSON Lines。

## 产品结论
- v0.1 的 live 模式应定位为“轻量实时刷新”。
- 不追求完整交互式 dashboard。
- 不追求替代 htop/btop 这类全屏终端面板。
- 核心是：在 SSH 场景下，以最小学习成本持续看当前状态。
