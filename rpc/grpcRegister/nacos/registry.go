package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"goStudy/rpc/grpcRegister/registry"
)

type Registry struct {
	nacosNamingClient naming_client.NamingClient
}

func (r Registry) Register(ctx context.Context, ins registry.ServiceInstance) error {

	panic("implement me")
}

func (r Registry) Unregister(ctx context.Context, ins registry.ServiceInstance) error {
	//TODO implement me
	panic("implement me")
}

func (r Registry) ListService(ctx context.Context, serviceName string) ([]registry.ServiceInstance, error) {
	//r.nacosNamingClient.SelectInstances()
	panic("implement me")
}

func (r Registry) Subscribe(serviceName string) (<-chan registry.Event, error) {
	//TODO implement me
	//r.nacosNamingClient.Subscribe()
	panic("implement me")
}

func (r Registry) Close() error {
	//TODO implement me
	panic("implement me")
}
