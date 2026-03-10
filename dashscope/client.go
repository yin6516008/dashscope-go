package dashscope

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// DefaultBaseURL 默认 API 地址
	DefaultBaseURL = "https://dashscope.aliyuncs.com/api/v1"
)

// Client 百炼应用调用客户端
type Client struct {
	apiKey   string
	appID    string
	baseURL  string
	httpCli  *http.Client
}

// ClientOption 客户端配置选项
type ClientOption func(*Client)

// WithBaseURL 设置自定义 API 地址（如私网调用）
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithHTTPClient 设置自定义 HTTP 客户端
func WithHTTPClient(cli *http.Client) ClientOption {
	return func(c *Client) {
		c.httpCli = cli
	}
}

// NewClient 创建百炼应用调用客户端
func NewClient(apiKey, appID string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		appID:   appID,
		baseURL: DefaultBaseURL,
		httpCli: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// completionURL 获取 completion 接口 URL
func (c *Client) completionURL() string {
	return fmt.Sprintf("%s/apps/%s/completion", c.baseURL, c.appID)
}

// doRequest 发送 HTTP 请求
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return c.httpCli.Do(req)
}
