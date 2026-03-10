package dashscope

// CallInput 调用输入参数
type CallInput struct {
	Prompt    string   `json:"prompt,omitempty"`
	SessionID string   `json:"session_id,omitempty"`
	Messages  []Message `json:"messages,omitempty"`
	MemoryID  string   `json:"memory_id,omitempty"`
	ImageList []string `json:"image_list,omitempty"`
	BizParams any      `json:"biz_params,omitempty"`
}

// Message 对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CallParameters 调用参数
type CallParameters struct {
	IncrementalOutput bool   `json:"incremental_output,omitempty"`
	HasThoughts      bool   `json:"has_thoughts,omitempty"`
	EnableThinking   bool   `json:"enable_thinking,omitempty"`
	RagOptions       any    `json:"rag_options,omitempty"`
}

// CallRequest 完整请求体
type CallRequest struct {
	Input      CallInput      `json:"input"`
	Parameters CallParameters `json:"parameters,omitempty"`
	Debug      map[string]any `json:"debug,omitempty"`
}

// CallOutput 响应输出
type CallOutput struct {
	Text         string   `json:"text"`
	SessionID    string   `json:"session_id"`
	FinishReason string   `json:"finish_reason"`
	Thoughts     []Thought `json:"thoughts,omitempty"`
}

// Thought 思考过程（深度思考模型）
type Thought struct {
	Action     string `json:"action"`
	ActionType string `json:"action_type"`
	Thought    string `json:"thought"`
	Response   string `json:"response"`
}

// Usage 用量信息
type Usage struct {
	Models []ModelUsage `json:"models"`
}

// ModelUsage 单模型用量
type ModelUsage struct {
	InputTokens  int    `json:"input_tokens"`
	OutputTokens int    `json:"output_tokens"`
	ModelID      string `json:"model_id"`
}

// CallResponse 完整响应
type CallResponse struct {
	Output    CallOutput `json:"output"`
	Usage     Usage      `json:"usage"`
	RequestID string     `json:"request_id"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}
