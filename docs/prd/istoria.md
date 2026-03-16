# iStoria PRD

## 需求背景
- 问题描述：用户在 SSH 登录不同机器后，通常需要组合使用 `top`、`free`、`df`、`vm_stat`、`sensors`、`system_profiler` 等多种命令查看硬件状态，输出分散、跨平台不统一、学习成本高。
- 目标用户：开发者、运维工程师、高级 Mac 用户、管理多台机器的技术用户。
- 价值主张：通过一个单文件 CLI，在 macOS / Linux / Windows 上统一查看本机硬件状态；默认输出适合终端阅读，JSON 输出适合脚本；并逐步增强规则、告警与 Apple Silicon 支持。

## 功能描述
- 核心功能：
  - 提供统一的跨平台硬件监控 CLI。
  - 支持 CPU、内存、磁盘、网络、传感器、电源信息查看。
  - 支持默认终端输出与 `--json` 输出。
  - 支持实时刷新显示当前硬件状态（live 模式），适合 SSH 场景下持续观察。
  - 为后续规则/告警能力预留统一数据结构。
- 用户流程：
  1. 下载单文件可执行程序。
  2. 将可执行文件放到目标机器。
  3. SSH 登录目标机器。
  4. 执行 `istoria stat` / `istoria cpu` / `istoria mem` / `istoria disk` / `istoria sensor` 查看状态。
  5. 在自动化场景中使用 `--json` 输出接入脚本。
- 界面要求：
  - 默认输出紧凑、清晰、适合 SSH 终端。
  - 重点先展示摘要，再按命令进入细分模块。
  - JSON 字段命名保持稳定，优先兼容脚本调用。
  - 平台差异字段允许扩展，但核心字段需统一。

## MVP 范围（v0.1）
- 必须实现：
  - `istoria stat`
  - `istoria cpu`
  - `istoria mem`
  - `istoria disk`
  - `istoria sensor`
  - `--json`
- 应该实现：
  - `istoria net`
  - `istoria power`
  - `--live` / `--watch` 实时刷新模式（至少支持 `stat` 与 `net`）
- 暂不实现：
  - 规则引擎完整能力
  - 告警推送通道
  - 历史趋势图
  - TUI / GUI
  - 多机集中控制台

## Apple Silicon 支持策略
- 第一版重点增强以下能力：
  - CPU 总体占用
  - 性能核 / 能效核相关信息（若可稳定获取）
  - 频率信息
  - 温度传感器
  - 风扇信息（若设备支持）
  - 功耗相关字段（若可稳定获取）
- 约束说明：
  - 不承诺所有机型字段完全一致。
  - 以“尽可能提供 Apple Silicon 增强指标”为产品原则。

## 命令设计原则
- 命令短、直观、适合 SSH 使用。
- 默认输出面向人类阅读。
- `--json` 面向脚本与自动化。
- 单个命令优先一屏读完主要信息。
- 提供轻量实时刷新能力，满足持续观察场景，但不演变为复杂 TUI。
- 避免用户为同一个问题在多个命令之间频繁跳转。

## 验收标准
1. [ ] 用户可在 macOS / Linux / Windows 上运行单文件版本。
2. [ ] `istoria stat` 可输出当前机器摘要信息。
3. [ ] `istoria cpu`、`istoria mem`、`istoria disk`、`istoria sensor` 可输出对应模块核心信息。
4. [ ] 核心命令支持 `--json` 输出。
5. [ ] `istoria stat --live` 或 `istoria net --live` 可按固定间隔持续刷新显示当前状态。
6. [ ] Apple Silicon 机器可展示额外专属字段或增强指标入口。
7. [ ] 常见 SSH 使用场景下，用户可在 1 分钟内理解主要输出结构并完成基本查看。

## 非目标（明确不做）
1. 不做第一版 GUI 桌面应用。
2. 不做第一版云账户体系。
3. 不做第一版多机集中式控制台。

## 优先级
- P0（必须）

## 成功指标
- 首次使用者在 5 分钟内完成安装并跑通至少一个核心命令。
- 用户能用 `istoria stat` 替代一部分分散的系统信息命令。
- JSON 输出能稳定接入脚本，不需要针对不同平台大幅改写调用逻辑。
- Apple Silicon 用户能感知到相比通用工具更好的指标覆盖。

## 风险与依赖
- 技术风险：
  - 三平台硬件指标获取能力不一致。
  - Apple Silicon 传感器、频率、功耗等能力可能需要平台特化实现。
  - Windows 与 Linux/macOS 的字段统一成本较高。
- 外部依赖：
  - Go 跨平台构建链路。
  - 系统信息采集库。
  - 各平台系统 API / 系统命令。

## 开发交接（Handoff）
1. Summary：交付一个适合 SSH 场景的跨平台硬件监控 CLI：iStoria，v0.1 先解决“统一查看本机状态”。
2. Spec 位置：`projects/istoria/docs/prd/istoria.md`
3. 优先级：P0，先完成命令骨架与基础数据采集方案。
4. 验收标准：支持 `stat/cpu/mem/disk/sensor` 与 `--json`，在三平台可运行，Apple Silicon 有增强字段入口。
5. 下一步：`ASSIGN: leader`，请评估技术可行性并拆分开发任务。
