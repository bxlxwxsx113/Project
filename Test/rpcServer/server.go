package main

import (
	rpclib "Project/Test/rpc"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

//{"method":"DemoService.Div","params":[{"A":3,"B":4},"id":1]}

func main() {
	rpc.Register(rpclib.DemoService{})
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Println(err)
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}
