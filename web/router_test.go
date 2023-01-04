package web

import (
	"fmt"
	"github.com/mattn/go-gtk/gtk"
	"math/rand"
	"testing"
	"time"
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

func TestPanic(t *testing.T) {

}
func TestGame(t *testing.T) {

	const (
		width      = 10
		height     = 20
		blockSize  = 30
		blockColor = "#00FF00"
	)

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("俄罗斯方块")
	window.Connect("destroy", gtk.MainQuit)

	board := gtk.NewDrawingArea()
	board.SetSizeRequest(width*blockSize, height*blockSize)
	window.Add(board)

	rand.Seed(time.Now().UnixNano())

	// insert a block at a random location
	//x := width / 2
	y := 0
	board.Connect("expose-event", func() {
		win := board.GetWindow().Show
		win()
		//board.GetWindow().SetBackground(gtk.NewGdkColor(blockColor))
		//board.GetWindow().Rectangle(board.GetStyle().BlackGC(), true, blockSize*x, blockSize*y, blockSize, blockSize)
	})

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			y++
			board.QueueDraw()
		}
	}()

	window.ShowAll()
	gtk.Main()
}
