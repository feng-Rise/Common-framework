package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

/*
	uintptr 其实就是一个数字 代表内存地址
	不要用uintptr保存地址，gc的时候没被清除的变量可能会被复制一份到另一个地址上，uintptr保存的地址再操作可能导致崩溃
	用unsafe.Pointer 来保存，gc时发生变化，会对指向的地址做修正

*/
type FieldAccessor interface {
	Field(field string) (int, error)
	SetField(field string, val int) error
}

type UnsafeAccessor struct {
	fields     map[string]FieldMeta
	entityAddr unsafe.Pointer
}

//val.UnsafeAddr()  用这个会panic  val是个指针,不能拿地址???
func NewUnsafeAccessor(entity interface{}) (*UnsafeAccessor, error) {
	if entity == nil {
		return nil, errors.New("invalid entity")
	}
	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid entity")
	}
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("val is not ptr")
	}

	fields := make(map[string]FieldMeta, typ.Elem().NumField())
	elemType := typ.Elem()
	for i := 0; i < elemType.NumField(); i++ {
		fd := elemType.Field(i)
		fields[fd.Name] = FieldMeta{offset: fd.Offset}
	}
	return &UnsafeAccessor{entityAddr: unsafe.Pointer(val.Elem().UnsafeAddr()), fields: fields}, nil
}

func (u *UnsafeAccessor) Field(field string) (int, error) {
	fdMeta, ok := u.fields[field]
	if !ok {
		return 0, fmt.Errorf("invalid field %s", field)
	}
	ptr := unsafe.Pointer(uintptr(u.entityAddr) + fdMeta.offset)
	if ptr == nil {
		return 0, fmt.Errorf("invalid address of the field: %s", field)
	}
	res := *(*int)(ptr)
	return res, nil
}

func (u *UnsafeAccessor) SetField(field string, val int) error {
	fdMeta, ok := u.fields[field]
	if !ok {
		return fmt.Errorf("invalid field %s", field)
	}
	ptr := unsafe.Pointer(uintptr(u.entityAddr) + fdMeta.offset)
	if ptr == nil {
		return fmt.Errorf("invalid address of the field: %s", field)
	}
	*(*int)(ptr) = val
	return nil
}

func (u *UnsafeAccessor) FieldAny(filed string) (interface{}, error) {
	meta, ok := u.fields[filed]
	if !ok {
		return 0, errors.New("不存在的字段")
	}

	res := reflect.NewAt(meta.typ, unsafe.Pointer(uintptr(u.entityAddr)+meta.offset))
	return res.Interface(), nil
}

type FieldMeta struct {
	typ    reflect.Type //用于fieldany
	offset uintptr
}
