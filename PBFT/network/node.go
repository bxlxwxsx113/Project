package network

import (
	"Project/PBFT/consensus"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Node struct {
	//节点名称
	NodeID string
	//节点信息
	NodeTable map[string]string
	//视图
	View *View
	//节点当前状态
	CurrentState *consensus.State
	//已经提交的消息
	CommitedMsgs []*consensus.RequestMsg
	//消息缓冲区
	MsgBuffer *MsgBuff
	//消息入口
	MsgEntrance chan interface{}
	//消息传输通道
	MsgDelivery chan interface{}
	//警报
	Alarm chan bool
	// 标识是否正在进行共识
	Active bool
}

//视图
type View struct {
	ID      int64
	Primary string
}

//各类消息缓冲区
//当某个某个消息到达时，可能节点正在进行共识某个阶段
//不能及时处理消息，将消息存储在缓冲区
type MsgBuff struct {
	ReqMsgs        []*consensus.RequestMsg
	PrePrepareMsgs []*consensus.PrePrePareMsg
	PrepareMsgs    []*consensus.VoteMsg
	CommitMsgs     []*consensus.VoteMsg
}

//创建新的节点
func NewNode(nodeID string) *Node {
	const viewID = 1000
	node := &Node{
		NodeID: nodeID,
		NodeTable: map[string]string{
			"N1": "localhost:5001",
			"N2": "localhost:5002",
			"N3": "localhost:5003",
			"N4": "localhost:5004",
		},
		View: &View{
			ID:      viewID,
			Primary: "N1",
		},
		CurrentState: nil,
		CommitedMsgs: make([]*consensus.RequestMsg, 0),
		MsgBuffer: &MsgBuff{
			ReqMsgs:        make([]*consensus.RequestMsg, 0),
			PrePrepareMsgs: make([]*consensus.PrePrePareMsg, 0),
			PrepareMsgs:    make([]*consensus.VoteMsg, 0),
			CommitMsgs:     make([]*consensus.VoteMsg, 0),
		},
		MsgEntrance: make(chan interface{}, 10),
		MsgDelivery: make(chan interface{}, 10),
		Alarm:       make(chan bool),
		Active:      false,
	}

	//开启信息调度
	go node.dispatchMsg()

	//开启一个警报触发器
	go node.alarmToDispatcher()

	//开启消息解析程序
	go node.resolveMsg()
	return node
}

func (node *Node) dispatchMsg() {
	for {
		select {
		case msg := <-node.MsgEntrance:
			err := node.routeMsg(msg)
			if err != nil {
				fmt.Println(err)
			}
		case <-node.Alarm:
			err := node.routeMsgWhenAlerm()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (node *Node) routeMsgWhenAlerm() []error {
	if node.CurrentState == nil {
		if len(node.MsgBuffer.ReqMsgs) != 0 {
			//创建切片
			msgs := make([]*consensus.RequestMsg,
				len(node.MsgBuffer.ReqMsgs))
			//将缓冲区中的消息拷贝到msgs
			copy(msgs, node.MsgBuffer.ReqMsgs)
			//将缓冲其中的消息清空掉
			node.MsgBuffer.ReqMsgs = make([]*consensus.RequestMsg, 0)
			//将消息发送到传送通道
			node.MsgDelivery <- msgs
		}
		if len(node.MsgBuffer.PrePrepareMsgs) != 0 {
			//创建切片
			msgs := make([]*consensus.PrePrePareMsg,
				len(node.MsgBuffer.PrePrepareMsgs))
			//将缓冲区中的消息拷贝到msgs
			copy(msgs, node.MsgBuffer.PrePrepareMsgs)
			//将缓冲其中的消息清空掉
			node.MsgBuffer.PrePrepareMsgs = make([]*consensus.PrePrePareMsg, 0)
			//将消息发送到传送通
			node.MsgDelivery <- msgs
		}
	} else {
		switch node.CurrentState.CurrentStage {
		case consensus.PrePrepared:
			if len(node.MsgBuffer.PrepareMsgs) != 0 {
				//创建切片
				msgs := make([]*consensus.VoteMsg,
					len(node.MsgBuffer.PrepareMsgs))
				//将缓冲区中的消息拷贝到msgs
				copy(msgs, node.MsgBuffer.PrepareMsgs)
				//将缓冲其中的消息清空掉
				node.MsgBuffer.PrepareMsgs = make([]*consensus.VoteMsg, 0)
				//将消息发送到传送通
				node.MsgDelivery <- msgs
			}
		case consensus.Prepared:
			if len(node.MsgBuffer.CommitMsgs) != 0 {
				//创建切片
				msgs := make([]*consensus.VoteMsg,
					len(node.MsgBuffer.CommitMsgs))
				//将缓冲区中的消息拷贝到msgs
				copy(msgs, node.MsgBuffer.CommitMsgs)
				//将缓冲其中的消息清空掉
				node.MsgBuffer.CommitMsgs = make([]*consensus.VoteMsg, 0)
				//将消息发送到传送通
				node.MsgDelivery <- msgs
			}
		}
	}
	return nil
}

//触发告警清空缓冲区
func (node *Node) alarmToDispatcher() {
	for {
		time.Sleep(time.Second)
		node.Alarm <- true
	}
}

//处理消息
func (node *Node) resolveMsg() {
	for {
		msgs := <-node.MsgDelivery
		fmt.Println(msgs)
		switch msgs.(type) {
		case []*consensus.RequestMsg:
			errs := node.resolveRequestMsg(msgs.([]*consensus.RequestMsg))
			if len(errs) != 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
			}
		case []*consensus.PrePrePareMsg:
			//处理序号分配消息
			errs := node.resolvePreprepareMsg(msgs.([]*consensus.PrePrePareMsg))
			if len(errs) != 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
			}
		case []*consensus.VoteMsg:
			voteMsgs := msgs.([]*consensus.VoteMsg)
			if len(voteMsgs) == 0 {
				break
			}
			if voteMsgs[0].MsgType == consensus.PrepareMsg {
				errs := node.reslvePrepareMsg(voteMsgs)
				if len(errs) != 0 {
					for _, err := range errs {
						fmt.Println(err)
					}
				}
			} else if voteMsgs[0].MsgType == consensus.CommitMsg {
				errs := node.reslvePrepareMsg(voteMsgs)
				if len(errs) != 0 {
					for _, err := range errs {
						fmt.Println(err)
					}
				}
			}
		}

	}
}

func (node *Node) GetCommit(commitMsg *consensus.VoteMsg) error {
	LogMsg(commitMsg)
	//判断该消息是否被提交过
	for i := 0; i < len(node.CommitedMsgs); i++ {
		if commitMsg.SequenceID == node.CommitedMsgs[i].SequenceID {
			return nil
		}
	}
	//共识进入相互交互节阶段
	replyMsg, committedMsg, err := node.CurrentState.Commit(commitMsg)
	if err != nil {
		return nil
	}
	if replyMsg == nil {
		return nil
	}
	if committedMsg == nil {
		return errors.New("Committed message is nil")
	}
	replyMsg.NodeID = node.NodeID
	node.Broadcast(commitMsg, "/commit")
	node.CommitedMsgs = append(node.CommitedMsgs, committedMsg)
	LogStage("Commit", true)
	node.Reply(replyMsg)
	LogStage("Reply", true)
	node.StateInit()
	return nil
}

func (node *Node) StateInit() {
	node.Active = false
	node.CurrentState = nil
	node.Alarm <- true
	node.MsgBuffer = &MsgBuff{
		ReqMsgs:        make([]*consensus.RequestMsg, 0),
		PrePrepareMsgs: make([]*consensus.PrePrePareMsg, 0),
		PrepareMsgs:    make([]*consensus.VoteMsg, 0),
		CommitMsgs:     make([]*consensus.VoteMsg, 0),
	}
	fmt.Println("New consensus begin!")
}

func (node *Node) Reply(replyMsg *consensus.ReplyMsg) error {
	bytes, err := json.Marshal(replyMsg)
	if err != nil {
		return err
	}
	send(node.NodeTable[node.View.Primary]+"replay", bytes)
	return nil
}

func (node *Node) resolveCommitMsg(msgs []*consensus.VoteMsg) []error {
	errs := make([]error, 0)
	for _, commitMsg := range msgs {
		err := node.GetCommit(commitMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}

func (node *Node) reslvePrepareMsg(msgs []*consensus.VoteMsg) []error {
	errs := make([]error, 0)
	for _, prepareMsg := range msgs {
		err := node.GetPrepare(prepareMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}

func (node *Node) GetPrepare(prepareMsg *consensus.VoteMsg) error {
	LogMsg(prepareMsg)
	//判断该消息是否被提交过
	for i := 0; i < len(node.CommitedMsgs); i++ {
		if prepareMsg.SequenceID == node.CommitedMsgs[i].SequenceID {
			return nil
		}
	}
	//共识进入相互交互节阶段
	commitMsg, err := node.CurrentState.PrePare(prepareMsg)
	if err != nil {
		return nil
	}
	if prepareMsg == nil {
		return nil
	}
	commitMsg.NodeID = node.NodeID
	LogStage("Prepare", true)
	node.Broadcast(commitMsg, "/commit")
	LogStage("Commit", false)
	return nil
}

//处理序号分配消息
func (node *Node) resolvePreprepareMsg(msgs []*consensus.PrePrePareMsg) []error {
	errs := make([]error, 0)
	for _, prePrepareMsg := range msgs {
		err := node.GetPreprepare(prePrepareMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}

//发送相互确认消息
func (node *Node) GetPreprepare(prePrepareMsg *consensus.PrePrePareMsg) error {
	LogMsg(prePrepareMsg)
	//判断该消息是否被提交过
	for i := 0; i < len(node.CommitedMsgs); i++ {
		if prePrepareMsg.SequenceID == node.CommitedMsgs[i].SequenceID {
			return nil
		}
	}
	fmt.Println("In pre-prepare")
	err := node.createStateForNewConsensus()
	if err != nil {
		return err
	}

	prepareMsg, err := node.CurrentState.PrePrePare(prePrepareMsg)
	if err != nil {
		return nil
	}
	if prepareMsg == nil {
		return nil
	}
	prepareMsg.NodeID = node.NodeID
	LogStage("Pre-prepare", true)
	node.Broadcast(prepareMsg, "/prepare")
	LogStage("Prepare", false)
	return nil
}

func (node *Node) resolveRequestMsg(msgs []*consensus.RequestMsg) []error {
	errs := make([]error, 0)
	for _, reqMsg := range msgs {
		err := node.GetReq(reqMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

//发送序号分配信息
func (node *Node) GetReq(reqMsg *consensus.RequestMsg) error {
	//打印消息
	LogMsg(reqMsg)
	err := node.createStateForNewConsensus()
	if err != nil {
		return nil
	}
	prePrepareMsg, err := node.CurrentState.StartConsensus(reqMsg)
	if err != nil {
		return nil
	}
	LogStage(fmt.Sprintf("leader prepare to send preprepareMsg(ViewID:%d, SequenceID:%d, Digest:%s, Request)",
		node.CurrentState.ViewID, prePrepareMsg.SequenceID, prePrepareMsg.Digest), false)
	if prePrepareMsg != nil {
		node.Broadcast(prePrepareMsg, "/preprepare")
		LogStage("Pre-prepare", true)
	}
	return nil
}

//将序号分配请求进行广播
func (node *Node) Broadcast(msg interface{}, path string) map[string]error {
	//创建map，用于存储错误消息
	errMap := make(map[string]error)
	for nodeID, url := range node.NodeTable {
		if nodeID == node.NodeID {
			continue
		}
		//序列化消息
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			errMap[nodeID] = err
			continue
		}
		fmt.Printf("send to %s\n", nodeID)
		send(url+path, jsonMsg)
	}
	if len(errMap) == 0 {
		return nil
	}
	return errMap
}

// 为当前节点即将进行的共识创建新的状态
func (node *Node) createStateForNewConsensus() error {
	//判断当前节点是否正在进行共识
	if node.Active {
		return errors.New("another consuses is on going!")
	}
	node.Active = true
	var lastSequenceID int64
	if len(node.CommitedMsgs) == 0 {
		lastSequenceID = -1
	} else {
		lastSequenceID = node.CommitedMsgs[len(node.CommitedMsgs)-1].SequenceID + 1
	}
	fmt.Println(len(node.CommitedMsgs))
	node.CurrentState = consensus.CreateState(node.View.ID, lastSequenceID)
	LogStage("Create the replica status ", true)
	return nil
}

//从消息入口把消息取出发送到消息传输通道,统一处理
func (node *Node) routeMsg(msg interface{}) []error {
	switch msg.(type) {
	case *consensus.RequestMsg:
		{
			fmt.Println(msg)
			if node.Active == false {
				//创建切片
				msgs := make([]*consensus.RequestMsg,
					len(node.MsgBuffer.ReqMsgs)+1)
				//将缓冲区中的消息拷贝到msgs
				copy(msgs, node.MsgBuffer.ReqMsgs)
				//将新到的消息追加到msgs中
				msgs = append(msgs, msg.(*consensus.RequestMsg))
				//将缓冲其中的消息清空掉
				node.MsgBuffer.ReqMsgs = make([]*consensus.RequestMsg, 0)
				//将消息发送到传送通道
				node.MsgDelivery <- msgs
			} else {
				//此时该节点正在进行共识，所以该消息暂时不能处理
				node.MsgBuffer.ReqMsgs = append(node.MsgBuffer.ReqMsgs,
					msg.(*consensus.RequestMsg))
			}
		}
	case *consensus.PrePrePareMsg:
		{
			fmt.Println(msg)
			if node.Active == false {
				//创建切片
				msgs := make([]*consensus.PrePrePareMsg,
					len(node.MsgBuffer.PrePrepareMsgs)+1)
				//将缓冲区中的消息拷贝到msgs
				copy(msgs, node.MsgBuffer.PrePrepareMsgs)
				//将新到的消息追加到msgs中
				msgs = append(msgs, msg.(*consensus.PrePrePareMsg))
				//将缓冲其中的消息清空掉
				node.MsgBuffer.PrePrepareMsgs = make([]*consensus.PrePrePareMsg, 0)
				//将消息发送到传送通
				node.MsgDelivery <- msgs
			} else {
				//此时该节点正在进行共识，所以该消息暂时不能处理
				node.MsgBuffer.PrePrepareMsgs = append(node.MsgBuffer.PrePrepareMsgs,
					msg.(*consensus.PrePrePareMsg))
			}
		}
	case *consensus.VoteMsg:
		if msg.(*consensus.VoteMsg).MsgType == consensus.PrepareMsg {
			if node.CurrentState == nil || node.CurrentState.CurrentStage != consensus.Prepared {
				node.MsgBuffer.PrepareMsgs = append(node.MsgBuffer.PrepareMsgs, msg.(*consensus.VoteMsg))
			} else {
				//创建切片
				msgs := make([]*consensus.VoteMsg,
					len(node.MsgBuffer.PrepareMsgs)+1)
				//将缓冲区中的消息拷贝到msgs
				copy(msgs, node.MsgBuffer.PrepareMsgs)
				//将新到的消息追加到msgs中
				msgs = append(msgs, msg.(*consensus.VoteMsg))
				//将缓冲其中的消息清空掉
				node.MsgBuffer.PrepareMsgs = make([]*consensus.VoteMsg, 0)
				//将消息发送到传送通
				node.MsgDelivery <- msgs
			}
		} else if msg.(*consensus.VoteMsg).MsgType == consensus.CommitMsg {
			if node.CurrentState == nil || node.CurrentState.CurrentStage != consensus.Prepared {
				node.MsgBuffer.CommitMsgs = append(node.MsgBuffer.CommitMsgs, msg.(*consensus.VoteMsg))
			} else {
				//创建切片
				msgs := make([]*consensus.VoteMsg,
					len(node.MsgBuffer.PrepareMsgs)+1)
				//将缓冲区中的消息拷贝到msgs
				copy(msgs, node.MsgBuffer.CommitMsgs)
				//将新到的消息追加到msgs中
				msgs = append(msgs, msg.(*consensus.VoteMsg))
				//将缓冲其中的消息清空掉
				node.MsgBuffer.CommitMsgs = make([]*consensus.VoteMsg, 0)
				//将消息发送到传送通
				node.MsgDelivery <- msgs
			}
		}
	}
	return nil
}

func (node *Node) GetReply(msg *consensus.ReplyMsg) {
	fmt.Println("Result: %s by %s\n", msg.Result, msg.NodeID)
}
