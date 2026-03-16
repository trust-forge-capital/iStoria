# iStoria 开发拆分建议：live mode v0.1

## 目标
以最小实现成本，为 iStoria 增加 SSH 场景友好的实时刷新能力。

## 总体拆分顺序
1. 先打通 `net --live`
2. 再复用能力到 `stat --live`
3. 最后视进度扩展到 `cpu --live` / `sensor --live`

这样可以优先交付最接近 `vnstat -l -i eth0` 的核心价值。

## 开发任务拆分

### TASK-LIVE-001：命令参数层支持
**目标**：所有核心命令可识别 live 相关参数。

建议内容：
- 增加全局或子命令参数：
  - `--live`
  - `--interval`
  - `--no-clear`
- 统一参数校验：
  - 默认 `1s`
  - 最小 `500ms`
- 明确 `--live --json` 的行为约束

**验收**：
- 命令能正确解析参数
- 非法 interval 有清晰报错或自动兜底

---

### TASK-LIVE-002：输出层 live renderer
**目标**：建立轻量实时输出机制。

建议内容：
- 实现清屏重绘模式
- 实现 no-clear 追加模式
- 增加顶部状态栏输出
- 统一 Ctrl+C 优雅退出处理

**验收**：
- 默认模式可清屏刷新
- `--no-clear` 可持续追加
- 用户 Ctrl+C 退出不报错

---

### TASK-LIVE-003：`istoria net --live`
**目标**：优先交付最核心 live 场景。

建议内容：
- 每个刷新周期输出：
  - RX 速率
  - TX 速率
  - 累计 RX/TX
  - 主接口
  - IP
- 支持 `--json --no-clear` 输出 JSONL

**验收**：
- 默认 1 秒刷新
- 实时数值变化可观察
- JSONL 可被脚本逐行解析

---

### TASK-LIVE-004：`istoria stat --live`
**目标**：支持整机摘要实时观察。

建议内容：
- 基于已有 `stat` 数据模型复用
- 每轮刷新显示：
  - CPU 摘要
  - Memory 摘要
  - Disk 摘要
  - Network 摘要
  - Sensor 摘要
  - Power 摘要
- 保证单屏可读

**验收**：
- 刷新过程中布局稳定
- 核心摘要信息可持续观察

---

### TASK-LIVE-005：JSONL 与脚本兼容
**目标**：明确 live + json 的机器可读输出。

建议内容：
- `--live --json --no-clear` 输出 JSON Lines
- 每行包含 timestamp 与 command
- 字段结构兼容 `json-schema-v0.1.md`

**验收**：
- 可用 `jq` / Python / shell 逐行消费
- 刷新不中断 JSON 解析

---

### TASK-LIVE-006：QA 用例补充
**目标**：保证 live 模式稳定性。

建议测试：
- interval 边界值
- Ctrl+C 退出
- no-clear 行为
- JSONL 可解析
- 异常采集时不直接崩溃
- macOS/Linux/Windows 基本行为一致性

## 推荐排期

### 第一阶段（最小可用）
- TASK-LIVE-001
- TASK-LIVE-002
- TASK-LIVE-003

### 第二阶段（完整 v0.1 体验）
- TASK-LIVE-004
- TASK-LIVE-005

### 第三阶段（收尾）
- TASK-LIVE-006

## 产品建议
- `net --live` 是 live 模式 MVP 中的 MVP，必须优先。
- `stat --live` 是体验增强项，但仍建议纳入 v0.1。
- 不建议在第一版 live 模式中加入复杂热键、交互切换、列排序等能力。

## 交付判断
若时间紧张，最小上线组合应为：
- `istoria net --live`
- `--interval`
- `--no-clear`
- `Ctrl+C` 退出
- `--live --json --no-clear` JSONL
