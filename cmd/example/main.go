// 百炼应用调用示例
// 使用前请设置环境变量: DASHSCOPE_API_KEY, DASHSCOPE_APP_ID
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yin6516008/dashscope-go/dashscope"
)

func main() {
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	appID := os.Getenv("DASHSCOPE_APP_ID")
	if apiKey == "" || appID == "" {
		fmt.Println("请设置环境变量: DASHSCOPE_API_KEY, DASHSCOPE_APP_ID")
		fmt.Println("示例: export DASHSCOPE_API_KEY=sk-xxx")
		fmt.Println("      export DASHSCOPE_APP_ID=your-app-id")
		os.Exit(1)
	}

	client := dashscope.NewClient(apiKey, appID)

	// 每个操作使用独立 context，避免前一轮超时影响后续调用
	fmt.Println("========== 1. 单轮对话 ==========")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 90*time.Second)
	singleTurn(client, ctx1)
	cancel1()

	fmt.Println("\n========== 2. 多轮对话 ==========")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 90*time.Second)
	multiTurn(client, ctx2)
	cancel2()

	fmt.Println("\n========== 3. 流式输出 ==========")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 90*time.Second)
	streamOutput(client, ctx3)
	cancel3()
}

func singleTurn(client *dashscope.Client, ctx context.Context) {
	resp, err := client.Call(ctx, "考勤异常怎么办？", dashscope.WithHasThoughts(true))
	if err != nil {
		fmt.Printf("调用失败: %v\n", err)
		return
	}
	// 输出思考过程和知识库检索（与 Python SDK 一致）
	printThoughtsAndObservations(resp.Output.Thoughts)
	printDocReferences(resp.Output.DocReferences)
	fmt.Printf("回复: %s\n", resp.Output.Text)
	fmt.Printf("session_id: %s\n", resp.Output.SessionID)
}

func multiTurn(client *dashscope.Client, ctx context.Context) {
	session := dashscope.NewSession(client)

	// 第一轮
	resp1, err := session.Call(ctx, "精臣科技的开票信息", dashscope.WithHasThoughts(true))
	if err != nil {
		fmt.Printf("第一轮失败: %v\n", err)
		return
	}
	fmt.Printf("用户: 精臣科技的开票信息\n")
	printThoughtsAndObservations(resp1.Output.Thoughts)
	printDocReferences(resp1.Output.DocReferences)
	fmt.Printf("助手: %s\n\n", resp1.Output.Text)

	// 第二轮（自动携带 session_id）
	resp2, err := session.Call(ctx, "精臣智慧的呢？", dashscope.WithHasThoughts(true))
	if err != nil {
		fmt.Printf("第二轮失败: %v\n", err)
		return
	}
	fmt.Printf("用户: 精臣智慧的呢？\n")
	printThoughtsAndObservations(resp2.Output.Thoughts)
	printDocReferences(resp2.Output.DocReferences)
	fmt.Printf("助手: %s\n", resp2.Output.Text)
}

// printThoughtsAndObservations 输出思考过程和知识库检索（与 Python SDK 一致）
// 参考: https://help.aliyun.com/zh/model-studio/call-single-agent-application/
func printThoughtsAndObservations(thoughts []dashscope.Thought) {
	for i, t := range thoughts {
		if t.Thought != "" {
			fmt.Printf("[思考 %d] %s\n", i+1, t.Thought)
		}
		if t.Observation != nil {
			formatted := formatJSONStr(t.Observation)
			fmt.Printf("[观察 %d]\n%s\n", i+1, formatted)
		}
	}
}

// formatJSONStr 尝试解析并格式化 JSON，包括嵌套的 content[].text（知识库检索结构）
func formatJSONStr(v any) string {
	var obj any
	switch val := v.(type) {
	case string:
		if err := json.Unmarshal([]byte(val), &obj); err != nil {
			return val
		}
	default:
		obj = v
	}
	// 递归格式化 content 中嵌套的 text
	if m, ok := obj.(map[string]any); ok {
		if content, ok := m["content"].([]any); ok {
			for _, item := range content {
				if itemMap, ok := item.(map[string]any); ok {
					if text, ok := itemMap["text"]; ok {
						if textStr, ok := text.(string); ok {
							var nested any
							if err := json.Unmarshal([]byte(textStr), &nested); err == nil {
								itemMap["text"] = nested
							}
						}
					}
				}
			}
		}
	}
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

// printDocReferences 输出知识库回答来源（需在应用内开启「展示回答来源」）
func printDocReferences(refs []dashscope.DocReference) {
	if len(refs) == 0 {
		return
	}
	fmt.Println("【回答来源】")
	for i, r := range refs {
		fmt.Printf("  [%d] ", i+1)
		if r.FileName != "" {
			fmt.Printf("文件: %s ", r.FileName)
		}
		if r.Title != "" {
			fmt.Printf("标题: %s ", r.Title)
		}
		if r.DocID != "" {
			fmt.Printf("(doc_id: %s) ", r.DocID)
		}
		if r.Content != "" {
			// 内容可能较长，截断显示
			content := r.Content
			if len(content) > 200 {
				content = content[:200] + "..."
			}
			fmt.Printf("\n    内容: %s", content)
		}
		fmt.Println()
	}
}

func streamOutput(client *dashscope.Client, ctx context.Context) {
	fmt.Println("用户: 巴黎出差的差旅费是多少？")
	var lastDocRefs []dashscope.DocReference
	err := client.Stream(ctx, "巴黎出差的差旅费是多少？", func(chunk *dashscope.StreamChunk) bool {
		// 思考过程和知识库检索（与 Python SDK 一致）
		if chunk.Output.Thoughts != nil {
			for i, t := range chunk.Output.Thoughts {
				if t.Thought != "" {
					fmt.Printf("[思考 %d] %s\n", i+1, t.Thought)
				}
				if t.Observation != nil {
					formatted := formatJSONStr(t.Observation)
					fmt.Printf("[观察 %d]\n%s\n", i+1, formatted)
				}
			}
		}
		// 收集回答来源
		if len(chunk.Output.DocReferences) > 0 {
			lastDocRefs = chunk.Output.DocReferences
		}
		// 回复
		if chunk.Output.Text != "" {
			fmt.Print(chunk.Output.Text)
		}
		return true
	}, dashscope.WithHasThoughts(true))
	if err != nil {
		fmt.Printf("\n流式调用失败: %v\n", err)
		return
	}
	printDocReferences(lastDocRefs)
	fmt.Println()
}
