# 阿里云百炼应用调用 Go SDK

基于 [阿里云百炼文档](https://help.aliyun.com/zh/model-studio/call-single-agent-application/) 封装的智能体应用调用 SDK。

## 安装

```bash
go get github.com/yin6516008/dashscope-go
```

## 功能特性

- ✅ 单轮对话
- ✅ 多轮对话（session_id 云端存储）
- ✅ 流式输出（SSE）
- ✅ 长期记忆（memory_id）
- ✅ 知识库检索（rag_options）
- ✅ 自定义业务参数（biz_params）

## 快速开始

### 1. 获取凭证

- **API Key**：前往 [密钥管理](https://bailian.console.aliyun.com/) 创建
- **应用 ID**：前往 [应用管理](https://bailian.console.aliyun.com/) 创建智能体应用，在应用卡片复制 APP_ID

### 2. 配置环境变量

```bash
export DASHSCOPE_API_KEY=sk-xxx
export DASHSCOPE_APP_ID=your-app-id
```

### 3. 运行示例

```bash
go run ./cmd/example/
```

## 使用示例

### 单轮对话

```go
import "github.com/yin6516008/dashscope-go/dashscope"

client := dashscope.NewClient(apiKey, appID)
resp, err := client.Call(ctx, "你是谁？")
fmt.Println(resp.Output.Text)
```

### 多轮对话

```go
session := dashscope.NewSession(client)
resp1, _ := session.Call(ctx, "你是谁？")
resp2, _ := session.Call(ctx, "你有什么技能？")  // 自动携带 session_id
```

### 流式输出

```go
client.Stream(ctx, "写一首诗", func(chunk *dashscope.StreamChunk) bool {
    fmt.Print(chunk.Output.Text)
    return true  // 返回 false 可提前终止
})
```

### 使用 messages 自行管理历史

```go
messages := []dashscope.Message{
    {Role: "user", Content: "你好"},
    {Role: "assistant", Content: "你好！有什么可以帮助你的？"},
    {Role: "user", Content: "推荐一部电影"},
}
resp, err := client.Call(ctx, "", dashscope.WithMessages(messages))
```

### 知识库检索与回答来源

```go
resp, err := client.Call(ctx, "请推荐一款3000元以下的手机",
    dashscope.WithRagOptions(map[string]any{
        "pipeline_ids": []string{"YOUR_PIPELINE_ID"},
    }),
)
// 回答来源：需在应用内开启「展示回答来源」
for _, ref := range resp.Output.DocReferences {
    fmt.Printf("来源: %s - %s\n", ref.FileName, ref.Content)
}
// 检索过程：WithHasThoughts(true) 时，thoughts 中 action_type=agentRag 的 observation
```

## 项目结构

```
dashscope-go/
├── dashscope/          # SDK 包
│   ├── client.go      # 客户端
│   ├── types.go       # 类型定义
│   ├── call.go        # 非流式调用
│   ├── stream.go      # 流式调用
│   └── session.go     # 多轮对话会话
├── cmd/example/       # 示例程序
└── README.md
```
