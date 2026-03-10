// 百炼应用调用示例
// 使用前请设置环境变量: DASHSCOPE_API_KEY, DASHSCOPE_APP_ID
package main

import (
	"context"
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
	// 输出思考过程
	if len(resp.Output.Thoughts) > 0 {
		fmt.Println("【思考过程】")
		for _, t := range resp.Output.Thoughts {
			if t.ActionType == "reasoning" && t.Thought != "" {
				fmt.Print(t.Thought)
			}
		}
		fmt.Println("\n【回复】")
	}
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
	printThoughts(resp1.Output.Thoughts)
	fmt.Printf("助手: %s\n\n", resp1.Output.Text)

	// 第二轮（自动携带 session_id）
	resp2, err := session.Call(ctx, "精臣智慧的呢？", dashscope.WithHasThoughts(true))
	if err != nil {
		fmt.Printf("第二轮失败: %v\n", err)
		return
	}
	fmt.Printf("用户: 精臣智慧的呢？\n")
	printThoughts(resp2.Output.Thoughts)
	fmt.Printf("助手: %s\n", resp2.Output.Text)
}

func printThoughts(thoughts []dashscope.Thought) {
	if len(thoughts) == 0 {
		return
	}
	fmt.Print("【思考过程】 ")
	for _, t := range thoughts {
		if t.ActionType == "reasoning" && t.Thought != "" {
			fmt.Print(t.Thought)
		}
	}
	fmt.Println()
}

func streamOutput(client *dashscope.Client, ctx context.Context) {
	fmt.Println("用户: 巴黎出差的差旅费是多少？")
	fmt.Print("【思考过程】 ")
	thinkingDone := false
	err := client.Stream(ctx, "巴黎出差的差旅费是多少？", func(chunk *dashscope.StreamChunk) bool {
		// 先输出思考过程
		for _, t := range chunk.Output.Thoughts {
			if t.ActionType == "reasoning" && t.Thought != "" {
				fmt.Print(t.Thought)
			}
		}
		// 再输出回复
		if chunk.Output.Text != "" {
			if !thinkingDone {
				fmt.Println("\n【回复】")
				thinkingDone = true
			}
			fmt.Print(chunk.Output.Text)
		}
		return true
	}, dashscope.WithHasThoughts(true))
	if err != nil {
		fmt.Printf("\n流式调用失败: %v\n", err)
		return
	}
	fmt.Println()
}
