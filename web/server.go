package web

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	//通过实现server 来启动服务，灵活性更强
	start(addr string) error
	AddRoute(method string, path string, handlers ...HandleFunc)
}

var _ Server = &HttpServer{}

type HttpServer struct {
	router
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		router: newRouter(),
	}
}

/*这里就是做实际请求的入口
构建起 Web 框架的上下文
查找路由树，并执行命中路由的代码*/
func (s *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.Server(ctx)
}

func (s *HttpServer) start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HttpServer) AddRoute(method string, path string, handlers ...HandleFunc) {
	//TODO implement me
	panic("implement me")
}

func (s *HttpServer) GET(path string, handlers HandleFunc) {
	s.AddRoute(http.MethodGet, path, handlers)
}

func (s *HttpServer) POST(path string, handlers HandleFunc) {
	s.AddRoute(http.MethodPost, path, handlers)
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
