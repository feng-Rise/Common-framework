package demo

import "context"

type Proxy interface {
	Invoke(ctx context.Context, req *Requset) (*Reponse, error)
}
type Requset struct {
	ServiceName string
	MethodName  string
	Args        interface{}
}

type Reponse struct {
	Data  []byte
	Error string
	Meta  map[string]string
}
