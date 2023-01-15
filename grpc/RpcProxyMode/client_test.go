package RpcProxyMode

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gostudy/grpc/RpcProxyMode/message"
	"log"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	cli := NewClient(":8081")
	us := &UserServiceClient{}
	err := InitClientProxy(us, cli)
	require.NoError(t, err)
	resp, err := us.GetById(context.Background(), &GetByIdReq{Id: 15})
	log.Printf("%v\n", resp)
}
func TestInitClientProxy(t *testing.T) {
	testCases := []struct {
		name    string
		service *UserServiceClient
		wantErr error
	}{
		{
			name:    "user_service",
			service: &UserServiceClient{},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &MockProxy{
				result: []byte(`{"Id": "14"}`),
			}
			err := InitClientProxy(tc.service, p)
			assert.Equal(t, tc.wantErr, err)
			resp, err := tc.service.GetById(context.Background(), &GetByIdReq{Id: 14})
			require.Nil(t, err)

			//断言 p的数据
			assert.Equal(t, &message.Request{
				ServiceName: "user-service",
				MethodName:  "GetById",
				Data:        []byte(`{"Id":"14"}`),
			}, p.req)
			assert.Equal(t, &GetByIdResp{
				Name: "feng",
			}, resp)
		})
	}
}

type MockProxy struct {
	req    *message.Request
	result []byte
}

func (m *MockProxy) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	m.req = req

	return &message.Response{
		Data: m.result,
	}, nil
}

func TestMakeFunc(t *testing.T) {

	// Create a function value that increments its argument by one.
	increment := func(x int) int {
		return x + 1
	}

	// Create a function value that wraps the increment function and intercepts calls to it.
	wrapped := reflect.MakeFunc(reflect.TypeOf(increment), func(args []reflect.Value) []reflect.Value {
		// Extract the first (and only) argument from the args slice.
		x := args[0].Int()

		// Call the increment function and return the result.
		result := increment(int(x))
		return []reflect.Value{reflect.ValueOf(result)}
	})

	// Call the wrapped function and print the result.
	result := wrapped.Call([]reflect.Value{reflect.ValueOf(5)})[0].Int()
	fmt.Println(result) // Output: 6

}
