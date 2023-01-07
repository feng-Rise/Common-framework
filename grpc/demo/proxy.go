package demo

import "context"

type Proxy interface {
	Invoke(ctx context.Context, req *Requset) (*Reponse, error)
}
type Requset struct {
	ServiceName string
	MethodName  string
	Data        []byte
	//Args        interface{}
}

type Reponse struct {
	Data  []byte
	Error string
	Meta  map[string]string
}

type UserServiceClient struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserServiceClient) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Name string `json:"name"`
}
