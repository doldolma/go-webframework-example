package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.Int("p", 8888, "리스닝할 포트를 입력해주세요.")
	addr := fmt.Sprintf(":%d", *port)

	server := NewServer()
	router := server.router

	// 라우터 그룹생성 /test 일때 -> 테스트 라우터로 이동해서 다시 분기
	groupRouter := router.Group("/test")

	// 결과적으로 /test/doldol 로 요청을 날려야 이 응답을 받을 수 있다.
	groupRouter.GET("/doldol", func(context *Context) {
		context.JSON(Json{
			"message": "GOOD",
		})
	})

	// 에러핸들링 미들웨어 추가
	server.Use(ErrorHandler)

	router.GET("/헬로우", func(c *Context) {
		c.ResponseWriter.Write([]byte("안녕하세요"))
	})

	router.GET("/헬로우/:이름", func(c *Context) {
		이름 := c.Params["이름"]
		인사말 := fmt.Sprintf("hello %s", 이름)
		c.JSON(Json{
			"message": 인사말,
		})
	})

	server.Run(addr)
}
