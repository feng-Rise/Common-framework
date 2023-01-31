package RpcProxyMode

import (
	"context"
	"encoding/json"
	"gostudy/grpc/RpcProxyMode/message"
	"net"
	"reflect"
)

type Server struct {
	services map[string]reflectionStub
}

func NewServer() *Server {
	res := &Server{
		services: map[string]reflectionStub{},
	}
	return res
}

func (s *Server) Register(service Service) error {
	s.services[service.Name()] = reflectionStub{
		value: reflect.ValueOf(service),
	}
	return nil
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			// 考虑输出日志，然后返回
			// return
			continue
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				// 这里考虑输出日志
				conn.Close()
				return
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		reqMsg, err := ReadMsg(conn)
		if err != nil {
			return err
		}
		req := message.DecodeReq(reqMsg)

		resp := &message.Response{
			Version:    req.Version,
			Compressor: req.Compressor,
			Serializer: req.Serializer,
			MessageId:  req.MessageId,
		}
		// 可以考虑找到本地的服务，然后发起调用
		service, ok := s.services[req.ServiceName]
		if !ok {
			// 返回客户端一个错误信息
			resp.Error = []byte("找不到服务")
			resp.SetHeadLength()
			_, err = conn.Write(message.EncodeResp(resp))
			if err != nil {
				return err
			}
			continue
		}

		ctx := context.Background()
		// 在这里可以检测超时
		//var cancel func() = func() {}
		//for key, value := range req.Meta {
		//	if key == "timeout" {
		//		deadline, err := strconv.ParseInt(value, 10, 64)
		//		if err != nil {
		//			// 返回客户端一个错误信息
		//			resp.Error = []byte(err.Error())
		//			resp.SetHeadLength()
		//			_, err = conn.Write(message.EncodeResp(resp))
		//			if err != nil {
		//
		//				return err
		//			}
		//			continue
		//		}
		//		ctx, cancel = context.WithDeadline(ctx, time.)
		//	} else {
		//		ctx = context.WithValue(ctx, key, value)
		//	}
		//}

		data, err := service.invoke(ctx, req)

		// 在这里可以检测超时
		resp.SetHeadLength()
		resp.BodyLength = uint32(len(data))
		resp.Data = data
		data = message.EncodeResp(resp)
		// 在这里检测超时
		_, err = conn.Write(data)
		if err != nil {
			return err
		}
	}
}

type reflectionStub struct {
	value reflect.Value
}

func (s *reflectionStub) invoke(ctx context.Context, req *message.Request) ([]byte, error) {
	method := s.value.MethodByName(req.MethodName)
	inType := method.Type().In(1)
	in := reflect.New(inType.Elem())
	err := json.Unmarshal(req.Data, in.Interface())
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), in})

	if len(res) > 1 && !res[1].IsZero() {
		return nil, res[1].Interface().(error)
	}
	return json.Marshal(res[0].Interface())
}
