package web

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	//通过实现server 来启动服务，灵活性更强
	Start(addr string) error
	addRoute(method string, path string, handler HandleFunc)
}

var _ Server = &HttpServer{}

type HttpServer struct {
	router
	mdls []Middleware
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		router: newRouter(),
	}
}
func (s *HttpServer) Use(mdls ...Middleware) {
	if s.mdls == nil {
		s.mdls = mdls
		return
	}
	s.mdls = append(s.mdls, mdls...)
}

/*
这里就是做实际请求的入口
构建起 Web 框架的上下文
查找路由树，并执行命中路由的代码
*/
func (s *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 最后一个应该是 HTTPServer 执行路由匹配，执行用户代码
	root := s.Server
	// 从后往前组装
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}
	// 第一个应该是回写响应的
	// 因为它在调用next之后才回写响应，
	// 所以实际上 flashResp 是最后一个步骤
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			//s.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

func (s *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HttpServer) GET(path string, handlers HandleFunc) {
	s.addRoute(http.MethodGet, path, handlers)

}

func (s *HttpServer) POST(path string, handlers HandleFunc) {
	s.addRoute(http.MethodPost, path, handlers)
}
func (s *HttpServer) Server(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}
	ctx.PathParams = mi.pathParams
	mi.n.handler(ctx)
}
