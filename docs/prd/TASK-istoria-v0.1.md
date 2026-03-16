# TASK: iStoria v0.1

## Summary
交付一个适合 SSH 场景的跨平台硬件监控 CLI：iStoria。第一阶段聚焦“统一查看本机状态”，优先覆盖 CPU、内存、磁盘、传感器与 JSON 输出。

## Spec位置
- `projects/istoria/docs/prd/istoria.md`

## 优先级
- **P0**
- **期望完成时间**：先完成技术评估与任务拆分；随后进入 v0.1 开发

## 验收标准
1. [ ] 可在 macOS / Linux / Windows 运行单文件版本。
2. [ ] 支持 `istoria stat` 输出机器摘要信息。
3. [ ] 支持 `istoria cpu` / `mem` / `disk` / `sensor` 输出核心模块信息。
4. [ ] 核心命令支持 `--json` 输出。
5. [ ] 支持 `istoria stat --live` 或 `istoria net --live` 的实时刷新模式。
6. [ ] Apple Silicon 存在增强字段入口或额外指标支持。

## 范围说明
### 本期必须做
- 命令骨架搭建
- 基础采集层接口设计
- `stat/cpu/mem/disk/sensor`
- `--json`

### 本期应该评估
- `net`
- `power`
- `--live` / `--watch` 轻量实时刷新模式
- Apple Silicon 增强指标的首版实现边界

### 本期不做
- 多机集中控制台
- GUI / TUI
- 云账户体系
- 完整规则/告警外发通道

## 风险与依赖
- 三平台指标采集能力差异
- Apple Silicon 指标需要平台特化实现
- 字段统一与 JSON 稳定性需要早期设计

## 下一步
- `ASSIGN: leader`
- 请 Leader 评估技术可行性、拆分任务，并协调 Developer / QA / Release 进入执行准备。
