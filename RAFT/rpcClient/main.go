package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Params struct {
	Width  int
	Height int
}

func main() {
	para := &Params{2, 3}
	rp, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	var ret = 0
	err = rp.Call("Rect.Area", para, &ret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("面积 = %d\n", ret)
	err = rp.Call("Rect.Perimeter", para, &ret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("周长 = %d\n", ret)
}
