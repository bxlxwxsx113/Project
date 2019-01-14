package consensus

type State struct {
	//视图ID
	ViewID int64
	//消息日志
	MsgLogs *MsgLogs
	//最后的区块ID
	LastSequenceID int64
	//当前状态
	CurrentStage Stage
}

//消息日志
type MsgLogs struct {
	ReqMsg      *RequestMsg
	PrePareMsgs map[string]*VoteMsg
	CommitMsgs  map[string]*VoteMsg
}

type Stage int

const (
	//节点已经创建，但是还未开始共识
	Idle Stage = iota
	PrePrepared
	Prepared
	Commited
)

//代表系统容纳错误节点个数
const f = 1

//创建状态
func CreateState(viewID int64, lastSequenceID int64) *State {
	return &State{
		ViewID: viewID,
		MsgLogs: &MsgLogs{
			ReqMsg:      nil,
			PrePareMsgs: make(map[string]*VoteMsg),
			CommitMsgs:  make(map[string]*VoteMsg),
		},
		LastSequenceID: lastSequenceID,
		CurrentStage:   Idle,
	}
}
