package RpcProxyMode

import (
	"context"
	"github.com/silenceper/pool"
	"gostudy/grpc/RpcProxyMode/message"
	"net"
	"time"
)

type Client struct {
	coonPool pool.Pool
}

func NewClient(addr string) *Client {
	pool, _ := pool.NewChannelPool(&pool.Config{
		InitialCap: 10,
		MaxCap:     100,
		MaxIdle:    50,
		Factory: func() (interface{}, error) {
			return net.Dial("tcp", addr)
		},
		IdleTimeout: time.Minute,
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()

		},
	})
	return &Client{
		coonPool: pool,
	}
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	coon, err := c.coonPool.Get()
	if err != nil {
		return nil, err
	}
	//发送请求
	data := message.EncodeReq(req)
	_, err = coon.(net.Conn).Write(data)
	if err != nil {
		return nil, err
	}
	//读取响应
	//
	respMsg, err := ReadMsg(coon.(net.Conn))
	// 还可以在这里检测超时
	if err != nil {
		return nil, err
	}
	return message.DecodeResp(respMsg), nil
}
