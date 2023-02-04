package web

import (
	"fmt"
	"testing"
)

func TestRouter_addRouter(T *testing.T) {

}

type testHandle func(ctx Context)

type testmiddle func(next testHandle) testHandle
type TeststructA struct {
	mdls []testmiddle
}

func aMiddle(ctx Context) {
	fmt.Println("aaaa")
}

func bMiddle(next testHandle) testHandle {
	return func(ctx Context) {
		next(ctx)
		fmt.Println("b")
	}
}
func TestMiddleWare(t *testing.T) {
	root := aMiddle
	astruct := &TeststructA{
		mdls: make([]testmiddle, 0),
	}
	astruct.mdls = append(astruct.mdls, bMiddle)
	root = astruct.mdls[0](root)
	var m testmiddle = func(next testHandle) testHandle {
		return func(ctx Context) {
			next(ctx)
			fmt.Println("123")
		}
	}
	root = m(root)
	var ctx Context
	root(ctx)
}

type logtest struct {
	logfunc func(val string)
}

func Testlog(t *testing.T) {
	a := &logtest{}
	a.logfunc("123")
}
