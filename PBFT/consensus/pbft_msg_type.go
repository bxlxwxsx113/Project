package consensus

//客户端请求消息
type RequestMsg struct {
	//时间
	Timestamp int64 `json:"timestamp"`
	//客户端ID
	ClientID string `json:"clientID"`
	//操作
	Operation string `json:"operation"`
	//请求编号
	SequenceID int64 `json:"sequenceID"`
}

//响应消息
type ReplyMsg struct {
	//视图ID
	ViewID int64 `json:"viewID"`
	//时间
	Timestamp int64 `json:"timestamp"`
	//客户端ID
	ClientID string `json:"clientID"`
	//发送消息的节点编号
	NodeID string `json:"nodeID"`
	//操作结果
	Result string `json:"result"`
}

//序号分配消息
type PrePrePareMsg struct {
	//视图ID
	ViewID int64 `json:"viewID"`
	//请求编号
	SequenceID int64 `json:"sequenceID"`
	//哈稀
	Digest string `json:"digest"`
	//请求消息
	RequestMsg `json:"requestMsg"`
}

type VoteMsg struct {
	//视图ID
	ViewID int64 `json:"viewID"`
	//请求编号
	SequenceID int64 `json:"sequenceID"`
	//哈稀
	Digest string `json:"digest"`
	//发送消息的节点编号
	NodeID string `json:"nodeID"`
	//消息类型
	MsgType MsgType `json:"msgType"`
}

//消息类型
type MsgType int

const (
	//相互交互
	PrepareMsg MsgType = iota
	//序号确认
	CommitMsg
)
