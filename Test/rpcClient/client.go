package main

import (
	rpclib "Project/Test/rpc"
	"fmt"
	"net"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)
	var result float64
	err = client.Call("DemoService.Div", rpclib.Args{10, 4}, &result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
