package web

import (
	"net/http"
	"net/url"
)

//为什么http请求是指针，返回是个接口
type Context struct {
	Req *http.Request
	// Resp 原生的 ResponseWriter。当你直接使用 Resp 的时候，
	// 那么相当于你绕开了 RespStatusCode 和 RespData。
	// 响应数据直接被发送到前端，其它中间件将无法修改响应
	// 其实我们也可以考虑将这个做成私有的
	Resp       http.ResponseWriter
	PathParams map[string]string

	// 缓存的响应部分
	// 这部分数据会在最后刷新
	RespStatusCode int
	RespData       []byte

	// 命中的路由
	MatchedRoute string

	// 万一将来有需求，可以考虑支持这个，但是需要复杂一点的机制
	// Body []byte 用户返回的响应
	// Err error 用户执行的 Error

	// 缓存的数据
	cacheQueryValues url.Values
}
