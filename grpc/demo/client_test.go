package demo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
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
			err := InitClientProxy(tc.service)
			assert.Equal(t, tc.wantErr, err)
			tc.service.GetById(context.Background(), &GetByIdReq{})
		})
	}
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

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

type GetByIdReq struct {
}

type GetByIdResp struct {
}
