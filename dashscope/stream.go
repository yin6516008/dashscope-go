package dashscope

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// StreamChunk 流式响应块
type StreamChunk struct {
	Output    CallOutput `json:"output"`
	Usage     Usage      `json:"usage"`
	RequestID string     `json:"request_id"`
}

// StreamCallback 流式输出回调，返回 false 可提前终止
type StreamCallback func(chunk *StreamChunk) bool

// Stream 流式调用智能体应用
func (c *Client) Stream(ctx context.Context, prompt string, callback StreamCallback, opts ...CallOption) error {
	req := &CallRequest{
		Input: CallInput{Prompt: prompt},
		Parameters: CallParameters{
			IncrementalOutput: true,
		},
		Debug: map[string]any{},
	}
	for _, opt := range opts {
		opt(req)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.completionURL(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("X-DashScope-SSE", "enable")

	resp, err := c.doRequest(ctx, httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return parseSSEStream(resp.Body, callback)
}

// parseSSEStream 解析 SSE 流式响应
// 格式: 每行 data: 后跟一个完整 JSON 对象（知识库检索时 observation 可能较大）
func parseSSEStream(r io.Reader, callback StreamCallback) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 256*1024), 4*1024*1024) // 256KB 初始，最大 4MB

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" || data == "[DONE]" {
			continue
		}
		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if !callback(&chunk) {
			return nil
		}
	}
	return scanner.Err()
}
