package demo

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"
	"reflect"
)

type Server struct {
	services map[string]Service
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]Service),
	}
}
func (s *Server) Register(service Service) error {
	s.services[service.Name()] = service
	return nil
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		coon, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			err := s.handleCoon(coon)
			if err != nil {
				coon.Close()
				return
			}
		}()

	}

}

func (s *Server) handleCoon(coon net.Conn) error {
	for {
		//读请求
		//执行
		//写回响应
		lenBytes := make([]byte, LengthBytes)
		_, err := coon.Read(lenBytes)
		if err != nil {
			return err
		}
		length := binary.BigEndian.Uint64(lenBytes)
		reqMsg := make([]byte, length)
		_, err = coon.Read(reqMsg)
		if err != nil {
			return err
		}

		req := Requset{}
		json.Unmarshal(reqMsg, &req)
		log.Printf("%v", req)

		//通过解析出来的req找到本地服务，调用的方法，入参，之后调用
		service, ok := s.services[req.ServiceName]
		if !ok {
			//编码返回一个错误信息
			return errors.New("can not find service")
		}
		//和grpc有区别
		method := reflect.ValueOf(service).MethodByName(req.MethodName)
		//把参数传进来
		ctx := context.Background()
		//methodReq 就是第二个参数
		/*type GetByIdReq struct {
			Id int
		}*/

		/*这里有个问题 req unmarshal解析出来的args 不是传参的结构体而是个map key是id，value是值，
		引为args是个interface，服务端这边是不知道它具体的类型的
		传过来的args 是  &GetByIdReq{Id: 14}，要把解析出来的arg也就是map传的参数赋值给methodreq
		解决方法是把args interface类型 换成[]byte，用json直接反序列化
		*/

		methodReq := reflect.New(method.Type().In(1))
		json.Unmarshal(req.Data, methodReq.Interface())
		resp := method.Call([]reflect.Value{reflect.ValueOf(ctx), methodReq.Elem()})

		methodResp := resp[0].Interface()
		//TODO 先不管error
		//methodErr := resp[0].Interface()

		//编码

		data, _ := json.Marshal(methodResp)
		data = EncodeMsg(data)
		_, err = coon.Write(data)
		if err != nil {
			return err
		}
		return nil
	}
}
