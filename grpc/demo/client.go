package demo

import (
	"context"
	"encoding/json"
	"reflect"
)

func InitClientProxy(servicie Service, p Proxy) error {
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
			ctx, ok := args[0].Interface().(context.Context)
			if !ok {
				panic("")
			}

			data, _ := json.Marshal(arg)
			req := &Requset{
				ServiceName: servicie.Name(),
				MethodName:  fieldType.Name,
				Data:        data,
			}

			//发送请求   有个接口  不希望在这里用具体的TCP操作
			resp, err := p.Invoke(ctx, req)

			//第一个返回值  真的返回值
			out := fieldType.Type.Out(0).Elem()
			if err != nil {
				results = append(results, reflect.Zero(out))
				results = append(results, reflect.ValueOf(err))
				return
			}
			//必须转成interface再unmarshal，不然不会解析，坑。。。
			first := reflect.New(out).Interface()
			//涉及序列化协议的转化 resp.data => first 用json做序列化

			err = json.Unmarshal(resp.Data, first)

			results = append(results, reflect.ValueOf(first))
			if err != nil {
				results = append(results, reflect.ValueOf(err))
			} else {
				results = append(results, reflect.Zero(reflect.TypeOf(new(error)).Elem()))
			}

			//fmt.Printf("%v\n", req)
			////不写下面这个 测试用例会报错
			//results = []reflect.Value{reflect.New(fieldType.Type.Out(0).Elem()), reflect.ValueOf(errors.New("你好"))}
			return
		})
		filedValue.Set(fn)
	}
	return nil
}

type Service interface {
	Name() string
}
