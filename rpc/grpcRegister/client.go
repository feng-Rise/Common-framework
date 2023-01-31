package grpcRegister

import (
	"google.golang.org/grpc"
	"gostudy/grpc/grpcRegister/proto"
	"gostudy/grpc/grpcRegister/registry"
)

func main() {
	var r registry.Registry
	proto.NewUserServiceClient()
	grpc.Dial()
}
