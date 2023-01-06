package demo

import (
	"errors"
	"fmt"
	"reflect"
)

func InitClientProxy(servicie interface{}) error {
	typ := reflect.TypeOf(servicie).Elem()
	value := reflect.ValueOf(servicie).Elem()
	numberFiled := value.NumField()
	for i := 0; i < numberFiled; i++ {
		fieldType := typ.Field(i)
		filedValue := value.Field(i)

		if !filedValue.CanSet() {
			continue
		}
		if fieldType.Type.Kind() != reflect.Func {
			continue
		}
		fn := reflect.MakeFunc(fieldType.Type, func(args []reflect.Value) (results []reflect.Value) {
			arg := args[1].Interface()

			req := &Requset{
				ServiceName: typ.Name(),
				MethodName:  fieldType.Name,
				Args:        arg,
			}
			fmt.Printf("%v\n", req)
			//不写下面这个 测试用例会报错
			results = []reflect.Value{reflect.New(fieldType.Type.Out(0).Elem()), reflect.ValueOf(errors.New("你好"))}
			return
		})
		filedValue.Set(fn)
	}
	return nil
}

type Requset struct {
	ServiceName string
	MethodName  string
	Args        interface{}
}
