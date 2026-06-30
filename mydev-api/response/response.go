package response

import "encoding/json"

// Response 统一响应格式
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewResponse 创建响应
func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// Success 成功响应（无数据）
func Success() *Response {
	return &Response{
		Code: 0,
		Msg:  "Success.",
		Data: []interface{}{},
	}
}

// SuccessWithData 成功响应（带数据）
func SuccessWithData(data interface{}) *Response {
	return &Response{
		Code: 0,
		Msg:  "Success.",
		Data: data,
	}
}

// Error 错误响应
func Error(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: []interface{}{},
	}
}

// ToJSON 转换为 JSON 字节
func (r *Response) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
