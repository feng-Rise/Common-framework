package demo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestInitClientProxy(t *testing.T) {
	testCases := []struct {
		name    string
		service *UserService
		wantErr error
	}{
		{
			name:    "user_service",
			service: &UserService{},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &MockProxy{
				result: []byte(`{"name": "feng"}`),
			}
			err := InitClientProxy(tc.service, p)
			assert.Equal(t, tc.wantErr, err)
			resp, err := tc.service.GetById(context.Background(), &GetByIdReq{Id: 14})
			require.Nil(t, err)

			//断言 p的数据
			assert.Equal(t, &Requset{
				ServiceName: "UserService",
				MethodName:  "GetById",
				Args:        &GetByIdReq{Id: 14},
			}, p.req)
			assert.Equal(t, &GetByIdResp{
				Name: "feng",
			}, resp)
		})
	}
}

type MockProxy struct {
	req    *Requset
	result []byte
}

func (m *MockProxy) Invoke(ctx context.Context, req *Requset) (*Reponse, error) {
	m.req = req

	return &Reponse{
		Data: m.result,
	}, nil
}

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Name string `json:"name"`
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
