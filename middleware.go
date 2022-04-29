package main

import (
	"log"
	"net/http"
	"time"
)

type MyMiddleware func(next MyHandlerFunc) MyHandlerFunc

func LogHandler(next MyHandlerFunc) MyHandlerFunc {
	return func(c *Context) {
		// 다음 미들웨어를 시작하기 전에 현재 시간 기록
		t := time.Now()

		// 다음 미들웨어 실행
		next(c)

		// 다음 미들웨어 실행이 끝나면 시간과 로그를 남긴다.
		log.Printf("[%s] %s %v \n",
			c.Request.Method,
			c.Request.URL.String(),
			time.Now().Sub(t))
	}
}

func ErrorHandler(next MyHandlerFunc) MyHandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := err.(error).Error()
				c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
				c.ResponseWriter.Write([]byte(errMsg))
				log.Println(errMsg)
			}
		}()
		next(c)
	}
}
