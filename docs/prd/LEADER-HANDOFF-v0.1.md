# LEADER HANDOFF: iStoria v0.1

## Summary
iStoria 是一个适合 SSH 场景的跨平台硬件监控 CLI。v0.1 目标是统一查看本机状态，并补充轻量实时刷新能力，优先解决 SSH 登录后快速观察当前机器状态的问题。

## Spec位置
- `projects/go/istoria/docs/prd/istoria.md`
- `projects/go/istoria/docs/prd/json-schema-v0.1.md`
- `projects/go/istoria/docs/prd/roadmap-v0.1.md`
- `projects/go/istoria/docs/prd/command-output-examples-v0.1.md`
- `projects/go/istoria/docs/prd/competitive-analysis-v0.1.md`
- `projects/go/istoria/docs/prd/live-mode-spec-v0.1.md`
- `projects/go/istoria/docs/prd/SPEC-CHANGE-live-mode.md`
- `projects/go/istoria/docs/prd/dev-breakdown-live-mode-v0.1.md`

## 优先级
- **P0**
- **期望完成时间**：先完成技术评估与任务拆分，再进入 v0.1 开发

## 验收标准
1. [ ] 可在 macOS / Linux / Windows 运行单文件版本。
2. [ ] 支持 `istoria stat` 输出整机摘要。
3. [ ] 支持 `istoria cpu` / `mem` / `disk` / `sensor` 输出核心模块信息。
4. [ ] 核心命令支持 `--json` 输出。
5. [ ] 支持 `istoria net --live` 的实时刷新模式。
6. [ ] 支持 `istoria stat --live` 或至少完成其技术方案与接口复用设计。
7. [ ] Apple Silicon 存在增强字段入口或额外指标支持。

## 关键产品决策
- SSH-first，本机执行，不做远程集中式监控平台。
- 单文件分发，跨平台统一命令体验。
- 输出同时服务于“人类阅读”和“自动化脚本”。
- live 模式采用 `--live` 轻量刷新方案，不做复杂 TUI。
- live 模式优先级：
  1. `net --live`
  2. `stat --live`

## 建议拆分顺序
### 第一阶段（MVP 核心）
- 命令骨架与参数层
- 基础采集层接口
- `stat/cpu/mem/disk/sensor`
- `--json`

### 第二阶段（live mode MVP）
- `--live`
- `--interval`
- `--no-clear`
- `net --live`
- JSONL 输出（`--live --json --no-clear`）

### 第三阶段（体验完善）
- `stat --live`
- Apple Silicon 增强字段整理
- QA 回归与边界用例

## 风险与依赖
- 三平台指标采集能力差异大。
- Apple Silicon 传感器和频率映射需要平台特化实现。
- live 模式需要兼顾低开销、可读性和 JSON 输出一致性。

## 推荐任务流转
- `ASSIGN: leader`
- 请 Leader：
  1. 评估技术可行性
  2. 拆分 Developer 任务
  3. 准备 QA 验收范围
  4. 确认 Release 所需构建产物与命名规范

## 建议发送给 Leader 的简版消息

```text
TASK: iStoria v0.1
优先级：P0
期望完成时间：先完成技术评估与任务拆分，再进入开发

Summary:
交付一个适合 SSH 场景的跨平台硬件监控 CLI：iStoria。v0.1 聚焦统一查看本机状态，并增加轻量实时刷新能力（类似 vnstat -l 的体验）。

Spec位置：
- projects/go/istoria/docs/prd/istoria.md
- projects/go/istoria/docs/prd/live-mode-spec-v0.1.md
- projects/go/istoria/docs/prd/dev-breakdown-live-mode-v0.1.md

关键验收标准：
- 支持 stat/cpu/mem/disk/sensor
- 支持 --json
- 支持 net --live
- 尽量支持 stat --live
- Apple Silicon 有增强字段入口

下一步：
ASSIGN: leader
请评估并拆分 Developer / QA / Release 执行任务。
```
