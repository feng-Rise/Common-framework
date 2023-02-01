package web

import "net/http"

//为什么http请求是指针，返回是个接口
type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
}
