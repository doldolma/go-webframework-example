package dolmago

import (
	"log"
	"net/http"
	"path"
	"strings"
)

type MyRouter struct {
	// http 메소드, URL 경로의 키와 핸들러 함수를 값으로 가진 2차원 맵,
	// { "GET": { "/path": 함수 } }
	handlers map[string]map[string]MyHandlerFunc
	basePath string
	server   *MyServer
}

// HandleFunc 핸들러 등록 메소드
func (router *MyRouter) HandleFunc(method, urlPattern string, handler MyHandlerFunc) {
	_, isExists := router.handlers[method]
	if !isExists {
		// 등록된 메소드가 없으면 새로 등록
		router.handlers[method] = make(map[string]MyHandlerFunc)
	}
	// handlers에 method와 URL 패턴과 핸들러 함수 등록
	router.handlers[method][urlPattern] = handler
	log.Printf("경로 등록완료 [%s]%s", method, path.Join(router.basePath, urlPattern))
}

// 요청 정보에 맞는 핸들러를 실행시키는 미들웨어(핸들러) 리턴
func (router *MyRouter) handler() MyHandlerFunc {
	return func(c *Context) {
		for pattern, handler := range router.handlers[c.Request.Method] {
			if ok, params := match(router.calculateAbsolutePath(pattern), c.Request.URL.Path); ok {
				// Context 생성
				context := Context{
					Params:         make(map[string]any),
					ResponseWriter: c.ResponseWriter,
					Request:        c.Request,
				}
				// Path 파라미터 컨텍스트에 저장
				for key, value := range params {
					context.Params[key] = value
				}
				// 요청 URL에 해당하는 handler 수행
				handler(&context)
				return
			}
		}
		http.NotFound(c.ResponseWriter, c.Request)
		return
	}
}

func (router *MyRouter) GET(urlPattern string, handler MyHandlerFunc) {
	router.HandleFunc(http.MethodGet, urlPattern, handler)
}

func (router *MyRouter) POST(urlPattern string, handler MyHandlerFunc) {
	router.HandleFunc(http.MethodPost, urlPattern, handler)
}

func (router *MyRouter) PUT(urlPattern string, handler MyHandlerFunc) {
	router.HandleFunc(http.MethodPut, urlPattern, handler)
}

func (router *MyRouter) PATCH(urlPattern string, handler MyHandlerFunc) {
	router.HandleFunc(http.MethodPatch, urlPattern, handler)
}

func (router *MyRouter) DELETE(urlPattern string, handler MyHandlerFunc) {
	router.HandleFunc(http.MethodDelete, urlPattern, handler)
}

// 모든 http 메소드에 매칭
func (router *MyRouter) ANY(urlPattern string, handler MyHandlerFunc) {
	router.HandleFunc(http.MethodGet, urlPattern, handler)
	router.HandleFunc(http.MethodPost, urlPattern, handler)
	router.HandleFunc(http.MethodPut, urlPattern, handler)
	router.HandleFunc(http.MethodPatch, urlPattern, handler)
	router.HandleFunc(http.MethodHead, urlPattern, handler)
	router.HandleFunc(http.MethodOptions, urlPattern, handler)
	router.HandleFunc(http.MethodDelete, urlPattern, handler)
	router.HandleFunc(http.MethodConnect, urlPattern, handler)
	router.HandleFunc(http.MethodTrace, urlPattern, handler)
}

// 핸들러의 pattern과 URL PATH가 일치하는지 체크
func match(pattern, path string) (bool, map[string]string) {

	if pattern == path {
		return true, nil
	}
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	// 개수가 일치하지 않으면 false
	if len(patterns) != len(paths) {
		return false, nil
	}

	// 패턴에 일치하는 URL의 매개변수를 담기 위한 params 맵
	params := make(map[string]string)

	// "/" 로 구분된 Pattern / Path 의 각 문자열을 비교
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
			// 패턴과 패스의 부분 문자열이 일치하면 바로 다음 루프 수행
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			// 패턴이 ‘:’ 문자로 시작하면 params에 URL params를 담은 후 다음 루프 수행
			params[patterns[i][1:]] = paths[i]
		case patterns[i] == "*":
			// 패턴에 *이 있으면 무조건 매칭 성공
			return true, params
		default:
			// 일치하는 경우가 없으면 false를 반환
			return false, nil
		}
	}
	return true, params
}

// 라우터를 리턴하는 함수 (라우터 안에 라우터 구현할 때)
func (router *MyRouter) Group(path string, handlers ...MyHandlerFunc) *MyRouter {
	subRouter := &MyRouter{
		handlers: make(map[string]map[string]MyHandlerFunc),
		basePath: router.calculateAbsolutePath(path),
	}
	router.ANY(path+"/*", subRouter.handler())
	router.ANY(path, subRouter.handler())
	return subRouter
}

func (router *MyRouter) calculateAbsolutePath(relativePath string) string {
	if router.basePath == "" {
		return relativePath
	}
	return path.Join(router.basePath, relativePath)
}
