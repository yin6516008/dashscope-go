package dashscope

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CallOption 调用选项
type CallOption func(*CallRequest)

// WithSessionID 设置会话 ID（多轮对话）
func WithSessionID(sessionID string) CallOption {
	return func(r *CallRequest) {
		r.Input.SessionID = sessionID
	}
}

// WithMessages 使用 messages 自行管理对话历史
func WithMessages(messages []Message) CallOption {
	return func(r *CallRequest) {
		r.Input.Messages = messages
	}
}

// WithMemoryID 设置长期记忆 ID
func WithMemoryID(memoryID string) CallOption {
	return func(r *CallRequest) {
		r.Input.MemoryID = memoryID
	}
}

// WithBizParams 设置业务参数（提示词变量、插件参数等）
func WithBizParams(params any) CallOption {
	return func(r *CallRequest) {
		r.Input.BizParams = params
	}
}

// WithRagOptions 设置知识库检索选项
func WithRagOptions(opts any) CallOption {
	return func(r *CallRequest) {
		r.Parameters.RagOptions = opts
	}
}

// WithHasThoughts 设置是否返回思考过程（深度思考模型）
func WithHasThoughts(enable bool) CallOption {
	return func(r *CallRequest) {
		r.Parameters.HasThoughts = enable
	}
}

// WithEnableThinking 开启思考模式（Qwen3 等模型）
func WithEnableThinking(enable bool) CallOption {
	return func(r *CallRequest) {
		r.Parameters.EnableThinking = enable
	}
}

// Call 调用智能体应用（非流式）
func (c *Client) Call(ctx context.Context, prompt string, opts ...CallOption) (*CallResponse, error) {
	req := &CallRequest{
		Input: CallInput{Prompt: prompt},
		Parameters: CallParameters{
			IncrementalOutput: false,
		},
		Debug: map[string]any{},
	}
	for _, opt := range opts {
		opt(req)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.completionURL(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.doRequest(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return nil, fmt.Errorf("api error [%s]: %s (request_id: %s)", errResp.Code, errResp.Message, errResp.RequestID)
		}
		return nil, fmt.Errorf("request failed: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result CallResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &result, nil
}
