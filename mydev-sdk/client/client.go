package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mydev/mydev-api/response"
)

const (
	headerTrace = "X-B3-TraceId"
	headerSpan  = "X-B3-SpanId"
)

// Client SDK 客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient 创建 SDK 客户端
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// CallResult 调用结果
type CallResult struct {
	Code int
	Msg  string
	Data interface{}
}

// HealthCheck 健康检查
func (c *Client) HealthCheck(traceId, spanId string) (*CallResult, error) {
	url := c.baseURL + "/srv/v1/hc"
	return c.doPost(url, nil, traceId, spanId)
}

// Divide 除法运算
func (c *Client) Divide(a, b float64, traceId, spanId string) (*CallResult, error) {
	url := c.baseURL + "/srv/v1/divide"

	reqBody := response.DivideRequest{
		A: a,
		B: b,
	}

	bodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	return c.doPost(url, bodyData, traceId, spanId)
}

// doPost 执行 POST 请求
func (c *Client) doPost(url string, body []byte, traceId, spanId string) (*CallResult, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置 headers
	req.Header.Set("Content-Type", "application/json")
	if traceId != "" {
		req.Header.Set(headerTrace, traceId)
	}
	if spanId != "" {
		req.Header.Set(headerSpan, spanId)
	}

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp response.Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &CallResult{
		Code: apiResp.Code,
		Msg:  apiResp.Msg,
		Data: apiResp.Data,
	}, nil
}
