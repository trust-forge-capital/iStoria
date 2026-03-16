# Spec change: iStoria live mode

## 变更摘要
- 为 iStoria 增加 **实时刷新显示当前硬件状态** 的能力。
- 需求来源：用户希望具备类似 `vnstat -l -i eth0` 的持续观察体验。
- 产品决策：采用 **`--live` 轻量实时刷新模式**，而不是复杂 TUI 或后台 daemon。

## 变更原因
- 当前 `stat/cpu/mem/disk/net/sensor` 主要满足“单次查看”。
- SSH 场景下，用户常需要持续观察网络速率、CPU 波动、温度变化。
- 若仍要求用户频繁重复执行命令，体验明显不如 live 模式直观。

## 变更范围
### 新增能力
- `istoria net --live`
- `istoria stat --live`
- `--interval <duration>`
- `--no-clear`

### 新增交互规则
- 默认刷新间隔：`1s`
- 最小刷新间隔：`500ms`
- 默认清屏重绘
- `Ctrl+C` 优雅退出
- 支持 `--live --json --no-clear` 输出 JSON Lines

## 优先级变化
- `live mode` 从“后续可选”提升为 **v0.1 应该实现**。
- 在 live 相关功能中，优先级顺序：
  1. `istoria net --live`
  2. `istoria stat --live`
  3. `cpu --live`
  4. `sensor --live`

## 受影响文档
- `docs/prd/istoria.md`
- `docs/prd/roadmap-v0.1.md`
- `docs/prd/command-output-examples-v0.1.md`
- `docs/prd/TASK-istoria-v0.1.md`
- `docs/prd/live-mode-spec-v0.1.md`

## 对开发的影响
- 需要在命令层增加 live 参数处理
- 需要在输出层增加“清屏重绘 / 追加输出”双模式
- 需要在采集层支持循环采样和节流控制
- 需要明确 JSONL 输出约束

## 对 QA 的影响
- 需要新增 live 模式用例：
  - 刷新间隔正确
  - Ctrl+C 正常退出
  - no-clear 行为正确
  - JSONL 输出可解析

## 产品结论
- 此变更不改变 iStoria 的核心定位。
- 该能力会增强 iStoria 在 SSH 场景下的竞争力，并更贴近“持续观察当前状态”的真实使用方式。
