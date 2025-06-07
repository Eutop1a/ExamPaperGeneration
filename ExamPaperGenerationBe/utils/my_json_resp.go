package utils

import (
	"encoding/json"
	"time"
)

// Response 定义响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// MakeResp 生成通用响应
func MakeResp(code int, msg string, data interface{}) string {
	r := Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	respBytes, _ := json.Marshal(r)
	return string(respBytes)
}

// Make200Resp 生成 200 状态码的响应
func Make200Resp(msg string, data interface{}) string {
	return MakeResp(200, msg, data)
}

// Make500Resp 生成 500 状态码的响应
func Make500Resp(msg string) string {
	data := "request parameter error, log time:" + time.Now().Format(time.RFC3339Nano)
	return MakeResp(500, msg, data)
}

// Make403Resp 生成 403 状态码的响应
func Make403Resp(msg string) string {
	data := "permission denied, log time:" + time.Now().Format(time.RFC3339Nano)
	return MakeResp(403, msg, data)
}

// Make400Resp 生成 400 状态码的响应
func Make400Resp(msg string) string {
	data := "bad request, log time:" + time.Now().Format(time.RFC3339Nano)
	return MakeResp(400, msg, data)
}
