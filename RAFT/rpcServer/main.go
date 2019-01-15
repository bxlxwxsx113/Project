package main

import (
	"log"
	"net/http"
	"net/rpc"
)

type Params struct {
	Width  int
	Height int
}

type Rect struct {
}

//函数必须时可到出的，首字母大写
//第二个采纳书时返回给客户端的采纳书，必须是指针类型
//函数必须有一个返回值：error
func (r *Rect) Area(p Params, ret *int) error {
	*ret = p.Width * p.Height
	return nil
}

func (r *Rect) Perimeter(p Params, ret *int) error {
	*ret = (p.Width + p.Height) * 2
	return nil
}

func main() {
	rect := &Rect{}
	rpc.Register(rect)
	rpc.HandleHTTP()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
