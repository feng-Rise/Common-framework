package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestRouter_addRouter(T *testing.T) {

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		contentType := "text/plain"
		path, _ := os.Getwd()
		path = path + "/test.txt"
		buf, err := os.ReadFile(path)
		if err != nil {
			return
		}

		ctx.Writer.Header().Set("Cache-Control", "private, max-age=31536000")

		ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(buf)))

		ctx.Data(http.StatusOK, contentType, buf)
	})
	r.Run(":8080")
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

func TestPanic(t *testing.T) {

}
