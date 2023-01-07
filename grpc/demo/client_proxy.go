package demo

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/silenceper/pool"
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

func (c *Client) Invoke(ctx context.Context, req *Requset) (*Reponse, error) {
	coon, err := c.coonPool.Get()
	if err != nil {
		return nil, err
	}
	//发送请求
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	data = EncodeMsg(data)
	_, err = coon.(net.Conn).Write(data)
	if err != nil {
		return nil, err
	}
	//读取响应
	lenBytes := make([]byte, LengthBytes)
	_, err = coon.(net.Conn).Read(lenBytes)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint64(lenBytes)
	respMsg := make([]byte, length)
	_, err = coon.(net.Conn).Read(respMsg)
	if err != nil {
		return nil, err
	}
	return &Reponse{
		Data: respMsg,
	}, nil
}
