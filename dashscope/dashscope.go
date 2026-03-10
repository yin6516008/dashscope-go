// Package dashscope 阿里云百炼智能体应用调用 SDK
//
// 支持功能：
//   - 单轮/多轮对话
//   - 流式输出
//   - 长期记忆
//   - 知识库检索
//   - 自定义业务参数
//
// 使用示例：
//
//	client := dashscope.NewClient(apiKey, appID)
//	resp, err := client.Call(ctx, "你好")
//	// 多轮对话
//	session := dashscope.NewSession(client)
//	resp1, _ := session.Call(ctx, "你是谁？")
//	resp2, _ := session.Call(ctx, "你有什么技能？")
//	// 流式输出
//	client.Stream(ctx, "写一首诗", func(chunk *dashscope.StreamChunk) bool {
//	    fmt.Print(chunk.Output.Text)
//	    return true
//	})
package dashscope
