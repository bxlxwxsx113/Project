package network

import (
	"Project/PBFT/consensus"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//发送消息
func send(url string, msg []byte) {
	buffer := bytes.NewBuffer(msg)
	http.Post("http://"+url, "application/json", buffer)
}

type Server struct {
	url  string
	node *Node
}

func NewServer(nodeID string) *Server {
	node := NewNode(nodeID)
	server := &Server{
		node.NodeTable[nodeID],
		node,
	}
	server.SetRouter()
	return server
}

func (server *Server) SetRouter() {
	http.HandleFunc("/req", server.getReg)
	http.HandleFunc("/preprepare", server.GetPreprepare)
	http.HandleFunc("/prepare", server.GetPreprepare)
	http.HandleFunc("/commit", server.GetCommit)
	http.HandleFunc("/replay", server.GetReply)

}

func (server *Server) getReg(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.RequestMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	}
	server.node.MsgEntrance <- &msg
}

func (server *Server) GetPreprepare(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.PrePrePareMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	}
	server.node.MsgEntrance <- &msg
}

func (server *Server) GetReply(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.ReplyMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	}
	server.node.GetReply(&msg)
}

func (server *Server) GetCommit(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.VoteMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	}
	server.node.MsgEntrance <- &msg
}

//开始服务
func (server *Server) Start() {
	fmt.Printf("Server will be start at %s...\n", server.url)
	err := http.ListenAndServe(server.url, nil)
	if err != nil {
		fmt.Println(err)
	}
}
