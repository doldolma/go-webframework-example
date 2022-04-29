package dolmago

import (
	"log"
	"net/http"
)

type MyServer struct {
	Router       *MyRouter
	middlewares  []MyMiddleware
	startHandler MyHandlerFunc
}

func NewServer() *MyServer {
	router := &MyRouter{
		handlers: make(map[string]map[string]MyHandlerFunc),
		basePath: "/",
	}
	server := &MyServer{Router: router}
	server.Router.server = server

	// 미들웨어 배열 등록
	server.middlewares = []MyMiddleware{
		LogHandler,
	}
	return server
}

func (server *MyServer) Use(middlewares ...MyMiddleware) {
	server.middlewares = append(server.middlewares, middlewares...)
}

func (server *MyServer) Run(addr string) {
	// 가장 먼저 실행되는 startHandler 를 라우터 핸들러 함수로 지정
	server.startHandler = server.Router.handler()

	// 등록된 미들웨어를 라우터 핸들러 앞에 하나씩 추가
	// 최종적으로 startHandler 는 LogHandler(startHandler()) 형태가 된다.
	for i := len(server.middlewares) - 1; i >= 0; i-- {
		server.startHandler = server.middlewares[i](server.startHandler)
	}

	// 웹 서버 시작
	log.Printf("Server is listening %s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		panic(err)
	}
}

func (server *MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Context 생성
	c := &Context{
		Params:         make(map[string]any),
		ResponseWriter: w,
		Request:        r,
	}
	for k, v := range r.URL.Query() {
		c.Params[k] = v[0]
	}
	server.startHandler(c)
}
