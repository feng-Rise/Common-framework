package RpcProxyMode

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestClient_Inovke(t *testing.T) {
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
			p := NewClient(":8081")
			err := InitClientProxy(tc.service, p)
			assert.Equal(t, tc.wantErr, err)
			resp, err := tc.service.GetById(context.Background(), &GetByIdReq{Id: 14})
			require.Nil(t, err)
			fmt.Printf("%v\n", resp)
			assert.Equal(t, &GetByIdResp{
				Name: "feng",
			}, resp)
		})
	}
}

type A struct{}

func (a *A) DoSomething() {
	fmt.Println("Doing something in A...")
}

func TestDemo(t *testing.T) {
	aValue := reflect.ValueOf(&A{})
	method := aValue.MethodByName("DoSomething")
	newFunc := reflect.MakeFunc(method.Type(), func(args []reflect.Value) (results []reflect.Value) {
		fmt.Println("Doing something in new function...")
		return nil
	})
	method.Set(newFunc)
	method.Call(nil)
}
