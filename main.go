package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var port = flag.Int("p", 8888, "리스닝할 포트를 입력해주세요.")
	addr := fmt.Sprintf(":%d", *port)

	log.Printf("Server is listening %d", *port)

	server := NewServer()
	router := server.router

	// 에러핸들링 미들웨어 추가
	server.Use(ErrorHandler)

	router.HandleFunc("GET", "/헬로우/:이름", func(c *Context) {
		이름 := c.Params["이름"]
		인사말 := fmt.Sprintf("hello %s", 이름)
		c.JSON(Json{
			"message": 인사말,
		})
	})

	server.Run(addr)
}
