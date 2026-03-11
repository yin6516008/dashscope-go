# 百炼 API 响应数据分析报告

基于 curl 调用的单轮、多轮、流式 API 响应数据，对比 SDK 实现提出优化建议。

## 一、数据来源

| 类型 | 文件 | 说明 |
|------|------|------|
| 单轮 | `single_turn.json` | 非流式，has_thoughts=true |
| 多轮 | `multi_turn_1.json`, `multi_turn_2.json` | 第一轮获取 session_id，第二轮携带 |
| 流式 | `stream.txt` | SSE 格式，incremental_output=true |

## 二、响应结构分析

### 2.1 Thought 结构（思考过程）

API 实际返回的 thought 字段比 SDK 定义更丰富：

| 字段 | SDK 支持 | API 实际返回 | 说明 |
|------|----------|-------------|------|
| action | ✅ | reasoning / search_knowledgebase | 工具/动作名 |
| action_type | ✅ | reasoning / mcp | reasoning=思考, mcp=知识库检索 |
| thought | ✅ | 思考内容 | 流式时为增量片段 |
| response | ✅ | 同 thought | 流式时为增量片段 |
| observation | ✅ | 知识库检索 JSON | 嵌套 content[].text |
| action_name | ❌ | "思考过程" | 中文描述 |
| action_input_stream | ❌ | `{"query":"..."}` | 工具输入 |
| arguments | ❌ | `{"query":"..."}` | 工具参数 |

**建议**：扩展 `Thought` 结构体，增加 `ActionName`、`ActionInputStream`、`Arguments` 字段，便于高级用户获取完整信息。

### 2.2 流式响应特性

1. **增量输出**：每个 chunk 的 `thought`、`response`、`text` 为**增量片段**，非累积
2. **空 chunk**：部分 chunk 如 `id:40` 仅含 `finish_reason`，无 thoughts/text
3. **usage 延迟**：流式过程中 `usage` 常为空 `{}`，最终 chunk 可能带用量
4. **编码**：curl 保存时若未指定 UTF-8，中文可能乱码；Go SDK 直接读 HTTP body 无此问题

### 2.3 单轮/多轮共性

- `output.text`：最终回复
- `output.session_id`：会话 ID，多轮第二轮需携带
- `output.thoughts`：思考过程（has_thoughts=true 时）
- `output.doc_references`：回答来源（应用开启「展示回答来源」时）
- `usage.models`：input_tokens、output_tokens、model_id

## 三、SDK 优化建议

### 已实现 ✅

- 单轮/多轮/流式调用
- thoughts、observation 解析（observation 支持 string/object）
- 流式 buffer 256KB~4MB 应对大 observation
- formatJSONStr 递归解析 content[].text

### 待优化

1. **Thought 扩展**：增加 `action_name`、`action_input_stream`、`arguments`
2. **流式文档**：在 README/示例中说明 thoughts 为增量片段，展示时需累积
3. **SSE 解析**：对 JSON 解析失败的 data 行记录日志或返回错误，避免静默跳过
4. **data 目录**：将 `data/*.json`、`data/*.txt` 加入 `.gitignore`（若含敏感/动态数据）
