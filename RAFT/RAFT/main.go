package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)

//节点信息
type nodeInfo struct {
	//节点id
	id string
	//端口号
	port string
}

type Raft struct {
	//节点信息
	node nodeInfo
	//锁
	mu sync.Mutex
	//当前节点编号
	me int
	//当前任期
	currentTerm int
	//为谁投票
	votedFor int
	//当前节点状态(0:follower,1:candidate,2:leader)
	state int
	//超时时间
	timeout int
	//当前领导
	currentLeader int
	//该节点最后一次处理数据的时间
	lastMessageTime int64
	//消息通道
	message chan bool
	//选举通道
	eclectCh chan bool
	//心跳通道
	heartbeat chan bool
	//子节点给主节点返回的心跳信号
	//主节点通过心跳的形式告诉子节点他还活着，以此维持自己的领导地位
	heartbeatRe chan bool
}

//声明leader结构体
type Leader struct {
	//任期
	Term int
	//领导编号
	LeaderId int
}

//设置节点个数
const raftCount = 2

var leader = Leader{0, -1}

//创建map，用于存储缓存信息
var bufferMessage = make(map[string]string)

//处理数据库信息
var mysqlMessage = make(map[string]string)

//操作消息数组下标
var messageId = 1

//用nodeTable存储每个节点中的键值对，键为节点名称，值为节点的端口号
var nodeTable map[string]string

func main() {
	//终端接收来的是数组
	if len(os.Args) > 1 {
		//接收终端输入的信息
		userId := os.Args[1]
		//字符串转换整型
		id, _ := strconv.Atoi(userId)
		fmt.Println(id)
		//定义节点id和端口号
		nodeTable = map[string]string{
			"1": ":8000",
			"2": ":8001",
		}
		//封装nodeInfo对象
		node := nodeInfo{id: userId, port: nodeTable[userId]}

		rf := Make(id) //创建节点对象
		//确保每个新建立的节点都有端口对应
		//127.0.0.1:8000
		rf.node = node
		//注册rpc
		go func() {
			//注册rpc，为了实现远程链接
			rf.raftRegisterRPC(node.port)
		}()
		if userId == "1" {
			go func() {
				//创建 getRequest() 回调方法
				http.HandleFunc("/req", rf.getRequest)
				fmt.Println("监听8080")
				if err := http.ListenAndServe(":8080", nil); err != nil {
					fmt.Println(err)
					return
				}
			}()
		}
	}
	for {
	}
}

var clientWriter http.ResponseWriter

func (rf *Raft) getRequest(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	//http://localhost:8080/req?age=20
	if len(request.Form["age"]) > 0 {
		clientWriter = writer
		fmt.Println("主节点广播客户端请求age:", request.Form["age"][0])

		param := Param{Msg: request.Form["age"][0], MsgId: strconv.Itoa(messageId)}
		messageId++
		if leader.LeaderId == rf.me {
			rf.sendMessageToOtherNodes(param)
		} else {
			//将消息转发给leader
			leaderId := nodeTable[strconv.Itoa(leader.LeaderId)]
			//连接远程rpc服务
			rpc, err := rpc.DialHTTP("tcp", "127.0.0.1"+leaderId)
			if err != nil {
				log.Fatal("\nrpc转发连接server错误:", leader.LeaderId, err)
			}
			var bo = false
			//首先给leader传递
			err = rpc.Call("Raft.ForwardingMessage", param, &bo)
			if err != nil {
				log.Fatal("\nrpc转发调用server错误:", leader.LeaderId, err)
			}
		}
	}
}

func (rf *Raft) sendMessageToOtherNodes(param Param) {
	bufferMessage[param.MsgId] = param.Msg
	// 只有领导才能给其它服务器发送消息
	if rf.currentLeader == rf.me {
		var success_count = 0
		fmt.Printf("领导者发送数据中 。。。\n")
		go func() {
			rf.broadcast(param, "Raft.LogDataCopy", func(ok bool) {
				//需要其它服务端回应
				rf.message <- ok
			})
		}()

		for i := 0; i < raftCount-1; i++ {
			fmt.Println("等待其它服务端回应")
			select {
			case ok := <-rf.message:
				if ok {
					success_count++
					if success_count >= raftCount/2 {
						rf.mu.Lock()
						rf.lastMessageTime = milliseconds()
						mysqlMessage[param.MsgId] = bufferMessage[param.MsgId]
						delete(bufferMessage, param.MsgId)
						if clientWriter != nil {
							fmt.Fprintf(clientWriter, "OK")
						}
						fmt.Printf("\n领导者发送数据结束\n")
						rf.mu.Unlock()
					}
				}
			}
		}
	}
}

//注册Raft对象，注册后的目的为确保每个节点（raft) 可以远程接收
func (node *Raft) raftRegisterRPC(port string) {
	//注册一个服务器
	rpc.Register(node)
	//把服务绑定到http协议上
	rpc.HandleHTTP()
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("注册rpc服务失败", err)
	}
}

//创建节点对象
func Make(me int) *Raft {
	//创建raft对象，初始化相关信息
	rf := &Raft{}
	rf.me = me
	//不为任何人投票
	rf.votedFor = -1
	//0 follower ,1 candidate ,2 leader
	rf.state = 0
	rf.timeout = 0
	//当前没有领导
	rf.currentLeader = -1
	//设置任期
	rf.setTerm(0)

	//初始化消息通道
	rf.message = make(chan bool)
	//初始化选举通道
	rf.eclectCh = make(chan bool)
	//初始化心跳通道
	rf.heartbeat = make(chan bool)
	//初始化子节点给主节点返回的心跳信号通道
	rf.heartbeatRe = make(chan bool)

	//选举leader
	go rf.election()
	//每个节点都有心跳功能
	go rf.sendLeaderHeartBeat()

	return rf
}

//选举成功后，应该广播所有的节点，本节点成为了leader
func (rf *Raft) sendLeaderHeartBeat() {
	for {
		select {
		//从心跳通道中获取数据
		case <-rf.heartbeat:
			rf.sendAppendEntriesImpl()
		}
	}
}

func (rf *Raft) sendAppendEntriesImpl() {
	//判断当前领导是否是自己
	if rf.currentLeader == rf.me {

		go func() {
			//封装信息
			param := Param{Msg: "leader heartbeat",
				Arg: Leader{rf.currentTerm, rf.me}}
			//广播
			rf.broadcast(param, "Raft.Heartbeat", func(ok bool) {
				rf.heartbeatRe <- ok
			})
		}()
		var success_count = 0
		for i := 0; i < raftCount-1; i++ {
			select {
			//从节点返回的信息
			case ok := <-rf.heartbeatRe:
				if ok {
					//成功返回信息，success_count加一
					success_count++
					if success_count >= raftCount/2 {
						rf.mu.Lock()
						//获取当前时间
						rf.lastMessageTime = milliseconds()
						fmt.Println("接收到了子节点们的返回信息")
						rf.mu.Unlock()
					}
				}
			}
		}
	}
}

//生成[min,max]之间的随机值
func randomRange(min, max int64) int64 {
	//设置随机时间
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

//获得当前时间（毫秒）
func milliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//选举leader
func (rf *Raft) election() {
	var result bool
	//每隔一段时间发一次心跳
	for {
		//生成[min,max]之间的随机值
		timeout := randomRange(1500, 3000)
		//设置该节点最后一次处理消息的时间
		rf.lastMessageTime = milliseconds()

		select {
		//间隔时间为1500-3000ms的随机值
		case <-time.After(time.Duration(timeout) * time.Millisecond):
		}

		result = false
		for !result {
			//选择leader
			result = rf.election_one_round(&leader)
		}
	}
}

func (rf *Raft) election_one_round(args *Leader) bool {
	//已经有了leader，并且不是自己，那么return
	if args.LeaderId > -1 && args.LeaderId != rf.me {
		fmt.Printf("%d已是leader，终止%d选举\n", args.LeaderId, rf.me)
		return true
	}
	//超时时间
	var timeout int64 = 2000
	//投票数量
	var vote int
	//用于是否开始心跳信号的方法
	var triggerHeartbeat bool
	//当前时间戳对应的毫秒
	last := milliseconds()
	success := false
	rf.mu.Lock()
	//要想参加选举需要先成为候选人
	rf.becomeCandidate()
	rf.mu.Unlock()
	fmt.Printf("candidate=%d start electing leader\n", rf.me)
	for {
		fmt.Printf("candidate=%d send request vote to server\n", rf.me)
		go func() {
			//广播消息
			rf.broadcast(Param{Msg: "send request vote"}, "Raft.ElectingLeader", func(ok bool) {
				//无论成功失败都需要发送到通道 避免堵塞
				rf.eclectCh <- ok
			})
		}()
		//投票数量设置为0
		vote = 0
		triggerHeartbeat = false
		for i := 0; i < raftCount-1; i++ {
			fmt.Printf("candidate=%d waiting for select for i=%d\n", rf.me, i)
			select {
			//从选举通道中获取信息
			case ok := <-rf.eclectCh:
				if ok {
					//投票数量加一
					vote++
					//如果当前所得票数大于节点的一半或者当前已经存在领导
					success = vote >= raftCount/2 || rf.currentLeader > -1
					if success && !triggerHeartbeat {
						fmt.Println("okok", args)
						//开始心跳
						triggerHeartbeat = true
						rf.mu.Lock()
						//成为领导者
						rf.becomeLeader()
						//当前任期加一
						args.Term = rf.currentTerm + 1
						//领导编号设置为自己
						args.LeaderId = rf.me
						rf.mu.Unlock()
						fmt.Printf("candidate=%d becomes leader\n", rf.currentLeader)
						//向心跳通道中发送true
						rf.heartbeat <- true
					}
				}
			}
			fmt.Printf("candidate=%d complete for select for i=%d\n", rf.me, i)
		}
		//如果时间超时或者投票数量大于节点的一半或者已经存在领导
		if (timeout+last < milliseconds()) || (vote >= raftCount/2 || rf.currentLeader > -1) {
			break
		} else {
			select {
			case <-time.After(time.Duration(5000) * time.Millisecond):
			}
		}
	}
	fmt.Printf("candidate=%d receive votes status=%t\n", rf.me, success)
	return success
}

//成为领导者
func (rf *Raft) becomeLeader() {
	rf.state = 2
	fmt.Println(rf.me, "成为了leader")
	rf.currentLeader = rf.me
}

//设置发送参数的数据类型
type Param struct {
	//消息
	Msg string
	//消息编号
	MsgId string
	//领导
	Arg Leader
}

//广播消息
func (rf *Raft) broadcast(msg Param, path string, fun func(ok bool)) {

	//设置不要自己给自己广播
	for nodeID, port := range nodeTable {
		//如果节点的id号和自己的id号相等，进行下一次循环
		if nodeID == rf.node.id {
			continue
		}

		//链接远程rpc
		rp, err := rpc.DialHTTP("tcp", "127.0.0.1"+port)
		//处理错误
		if err != nil {
			fun(false)
			continue
		}

		var bo = false
		//远程调用
		err = rp.Call(path, msg, &bo)
		//处理错误
		if err != nil {
			fun(false)
			continue
		}
		fun(bo)
	}

}

//成为候选人
func (rf *Raft) becomeCandidate() {
	//自己的状态为追随者或者当前没有领导
	if rf.state == 0 || rf.currentLeader == -1 {
		//候选人状态
		rf.state = 1
		//为自己投票
		rf.votedFor = rf.me
		//当前任前加一
		rf.setTerm(rf.currentTerm + 1)
		//当前没有领导
		rf.currentLeader = -1
	}
}

//设置任期
func (rf *Raft) setTerm(term int) {
	rf.currentTerm = term
}

//Rpc处理
func (rf *Raft) ElectingLeader(param Param, a *bool) error {
	//给leader投票
	*a = true
	rf.lastMessageTime = milliseconds()
	return nil
}

func (rf *Raft) Heartbeat(param Param, a *bool) error {
	//参数中的任期小于当前的任期
	fmt.Println("\nrpc:heartbeat:", rf.me, param.Msg)
	if param.Arg.Term < rf.currentTerm {
		*a = false
	} else {
		leader = param.Arg
		fmt.Printf("%d收到leader%d心跳\n", rf.currentLeader, leader.LeaderId)
		*a = true
		rf.mu.Lock()
		rf.currentLeader = leader.LeaderId
		rf.votedFor = leader.LeaderId
		rf.state = 0
		rf.lastMessageTime = milliseconds()
		fmt.Printf("server = %d learned that leader = %d\n", rf.me, rf.currentLeader)
		rf.mu.Unlock()
	}
	return nil
}

//连接到leader节点
func (rf *Raft) ForwardingMessage(param Param, a *bool) error {
	fmt.Println("\nrpc:forwardingMessage:", rf.me, param.Msg)

	rf.sendMessageToOtherNodes(param)

	*a = true
	rf.lastMessageTime = milliseconds()

	return nil
}

//接收leader传过来的日志
func (r *Raft) LogDataCopy(param Param, a *bool) error {
	fmt.Println("\nrpc:LogDataCopy:", r.me, param.Msg)
	bufferMessage[param.MsgId] = param.Msg
	*a = true
	return nil
}
