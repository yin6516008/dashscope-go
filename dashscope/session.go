package dashscope

import "context"

// Session 多轮对话会话，自动维护 session_id
type Session struct {
	client    *Client
	sessionID string
}

// NewSession 创建多轮对话会话
func NewSession(client *Client) *Session {
	return &Session{client: client}
}

// SessionID 返回当前会话 ID
func (s *Session) SessionID() string {
	return s.sessionID
}

// Call 发送消息并获取回复（非流式）
func (s *Session) Call(ctx context.Context, prompt string, opts ...CallOption) (*CallResponse, error) {
	opts = append([]CallOption{WithSessionID(s.sessionID)}, opts...)
	resp, err := s.client.Call(ctx, prompt, opts...)
	if err != nil {
		return nil, err
	}
	if resp.Output.SessionID != "" {
		s.sessionID = resp.Output.SessionID
	}
	return resp, nil
}

// Stream 流式发送消息并获取回复
func (s *Session) Stream(ctx context.Context, prompt string, callback StreamCallback, opts ...CallOption) error {
	opts = append([]CallOption{WithSessionID(s.sessionID)}, opts...)
	return s.client.Stream(ctx, prompt, func(chunk *StreamChunk) bool {
		if chunk.Output.SessionID != "" {
			s.sessionID = chunk.Output.SessionID
		}
		return callback(chunk)
	}, opts...)
}
