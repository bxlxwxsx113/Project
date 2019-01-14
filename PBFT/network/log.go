package network

import (
	"Project/PBFT/consensus"
	"fmt"
)

func LogMsg(msg interface{}) {
	fmt.Println(msg)
	switch msg.(type) {
	case *consensus.RequestMsg:
		reqMsg := msg.(*consensus.RequestMsg)
		fmt.Printf("[REQUEST] ClientID:%s, Timestamp:%d, Operateion:%s\n",
			reqMsg.ClientID, reqMsg.Timestamp, reqMsg.Operation)
	case *consensus.PrePrePareMsg:
		prePrepareMsg := msg.(*consensus.PrePrePareMsg)
		fmt.Printf("[PREPREPARE] ClientID:%s, Timestamp:%d, Operateion:%s\n",
			prePrepareMsg.RequestMsg.ClientID, prePrepareMsg.RequestMsg.Timestamp, prePrepareMsg.RequestMsg.Operation)
	case *consensus.VoteMsg:
		voteMsg := msg.(*consensus.VoteMsg)
		if voteMsg.MsgType == consensus.PrepareMsg {
			fmt.Printf("[PREPARE] NodeID:%s\n", voteMsg.NodeID)
		} else if voteMsg.MsgType == consensus.CommitMsg {
			fmt.Printf("[COMMIT] NodeID:%s\n", voteMsg.NodeID)
		}
	}
}

func LogStage(stage string, isDone bool) {
	if isDone {
		fmt.Printf("[STAGE-DONE] %s\n", stage)
	} else {
		fmt.Printf("[STAGE-BEGIN] %s\n", stage)
	}
}
