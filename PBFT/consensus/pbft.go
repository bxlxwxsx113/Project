package consensus

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type PBFT interface {
	//开始共识
	StartConsensus(request *RequestMsg) (*PrePrePareMsg, error)
	//序号分配
	PrePrePare(prePrepareMsg *PrePrePareMsg) (*VoteMsg, error)
	//相互交互
	PrePare(prePareMsg *VoteMsg) (*VoteMsg, error)
	//序号确认
	Commit(commitMsg *VoteMsg) (*ReplyMsg, *RequestMsg, error)
}

func (state *State) StartConsensus(request *RequestMsg) (*PrePrePareMsg, error) {
	sequenceID := time.Now().Unix()
	if state.LastSequenceID != -1 {
		for state.LastSequenceID >= sequenceID {
			sequenceID += 1
		}
	}
	//为请求编号
	request.SequenceID = sequenceID
	//将请求写入到消息日志中
	state.MsgLogs.ReqMsg = request
	//计算请求hash
	digest, err := digest(interface{}(request))
	if err != nil {
		return nil, err
	}
	//修改当前状态位序号分配阶段
	state.CurrentStage = PrePrepared
	prePrepareMsg := &PrePrePareMsg{
		ViewID:     state.ViewID,
		SequenceID: sequenceID,
		RequestMsg: *request,
		Digest:     digest,
	}
	return prePrepareMsg, nil
}

func (state *State) PrePrePare(prePrepareMsg *PrePrePareMsg) (*VoteMsg, error) {
	state.MsgLogs.ReqMsg = &prePrepareMsg.RequestMsg
	//验证消息合法性
	if !state.verifyMsg(prePrepareMsg.ViewID, prePrepareMsg.SequenceID, prePrepareMsg.Digest) {
		return nil, errors.New("pre-prepare message is error")
	}
	//修改当前节点的状态
	state.CurrentStage = PrePrepared
	return &VoteMsg{
		ViewID:     state.ViewID,
		SequenceID: prePrepareMsg.SequenceID,
		Digest:     prePrepareMsg.Digest,
		MsgType:    PrepareMsg,
	}, nil
}

//相互交互
func (state *State) PrePare(prePareMsg *VoteMsg) (*VoteMsg, error) {
	//消息验证
	if !state.verifyMsg(prePareMsg.ViewID, prePareMsg.SequenceID, prePareMsg.Digest) {
		return nil, errors.New("prepare message is error")
	}
	state.MsgLogs.PrePareMsgs[prePareMsg.NodeID] = prePareMsg
	fmt.Printf("[Prepare-Vote]:%d\n", len(state.MsgLogs.PrePareMsgs))
	if state.prepared() {
		return nil, nil
	}
	state.CurrentStage = Prepared
	return &VoteMsg{
		ViewID:     state.ViewID,
		SequenceID: prePareMsg.SequenceID,
		Digest:     prePareMsg.Digest,
		MsgType:    CommitMsg,
	}, nil
}

//序号确认
func (state *State) Commit(commitMsg *VoteMsg) (*ReplyMsg, *RequestMsg, error) {
	if !state.verifyMsg(commitMsg.ViewID, commitMsg.SequenceID, commitMsg.Digest) {
		return nil, nil, errors.New("prepare message is error")
	}
	state.MsgLogs.CommitMsgs[commitMsg.NodeID] = commitMsg
	fmt.Printf("[Commit-Vote]:%d\n", len(state.MsgLogs.CommitMsgs))
	if !state.commited() {
		return nil, nil, nil
	}
	state.CurrentStage = Commited
	return &ReplyMsg{
		ViewID:    commitMsg.ViewID,
		Timestamp: state.MsgLogs.ReqMsg.Timestamp,
		ClientID:  state.MsgLogs.ReqMsg.ClientID,
		Result:    "Excuted",
	}, state.MsgLogs.ReqMsg, nil
}

func (state *State) commited() bool {
	if state.prepared() {
		return false
	}
	if len(state.MsgLogs.CommitMsgs) < 2*f {
		return false
	}
	return true
}

//检查相互交互阶段信息的合法性
func (state *State) prepared() bool {
	if state.MsgLogs.ReqMsg == nil {
		return false
	}
	if len(state.MsgLogs.PrePareMsgs) < 2*f {
		return false
	}
	return true
}

//验证消息的合法性
func (state *State) verifyMsg(viewID int64, sequenceID int64, digestGot string) bool {
	if state.ViewID != viewID {
		return false
	}
	if state.LastSequenceID != -1 {
		if state.LastSequenceID > sequenceID {
			return false
		}
	}
	if digest, _ := digest(state.MsgLogs.ReqMsg); digest != digestGot {
		return false
	}
	return true
}

func digest(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return Hash(bytes), nil
}
