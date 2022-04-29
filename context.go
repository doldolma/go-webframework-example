package main

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Params         map[string]any
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type MyHandlerFunc func(*Context)

type Json map[string]any

func (c *Context) JSON(value any) {
	// Status 200
	c.ResponseWriter.WriteHeader(http.StatusOK)
	// 콘텐츠타입 JSON 으로 설정
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	// value 값을 json으로 출력
	if err := json.NewEncoder(c.ResponseWriter).Encode(value); err != nil {
		// JSON 인코딩 실패시 에러 응답
		c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.ResponseWriter.Write([]byte(err.Error()))
	}
}
