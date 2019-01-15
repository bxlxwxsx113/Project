package main

import (
	"fmt"
	"github.com/insionng/macross/libraries/gommon/log"
	"math/rand"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

//模拟三节点的分布式选举

//定义常量3
const raftCount = 3

//声明leader
type Leader struct {
	//任期
	Term int
	//领导编号
	LeaderId int
}

var leader = Leader{0, -1}

//声明raft节点
type Raft struct {
	mu          sync.Mutex
	me          int //节点编号
	currentTerm int //当前任期
	votedFor    int //为哪个节点投票
	//0:follower  1:candidate   2:leader
	state           int   //当前节点的状态
	lastMessageTime int64 //发送最后一条消息的时间
	currentLeader   int   //当前节点的领导

	//消息通道
	message chan bool
	//选举通道
	electCh chan bool
	//心跳信号
	heartBeat chan bool
	//心跳返回信号
	hearbeatRe chan bool
	//超时时间
	timeout int
}

func main() {
	//创建三个节点
	for i := 0; i < raftCount; i++ {
		Make(i)
	}

	rpc.Register(new(Raft))
	rpc.HandleHTTP()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//创建节点
func Make(me int) *Raft {
	rf := &Raft{}
	rf.me = me
	//现在是创建节点，还没开始投票
	rf.votedFor = -1
	//选举尚未开始，大家都是follower
	rf.state = 0
	rf.timeout = 0
	//最初没有领导
	rf.currentLeader = -1
	rf.setTerm(0)

	rf.electCh = make(chan bool)
	rf.message = make(chan bool)
	rf.heartBeat = make(chan bool)
	rf.hearbeatRe = make(chan bool)

	//随机出种子
	rand.Seed(time.Now().UnixNano())
	//选举
	go rf.election()

	go rf.sendLeaderHeartBeat()

	return rf
}

func (rf *Raft) sendLeaderHeartBeat() {
	for {
		select {
		case <-rf.heartBeat:
			rf.sendAppendEntriesImpl()
		}
	}
}

func (rf *Raft) sendAppendEntriesImpl() {
	//判断当前是否是leader节点
	if rf.currentLeader == rf.me {
		var success_count = 0

		for i := 0; i < raftCount; i++ {
			if i != rf.me {
				go func() {
					rp, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
					if err != nil {
						log.Fatal(err)
					}
					var ok = false
					err = rp.Call("Raft.Communication", Param{"hello"}, &ok)
					if err != nil {
						log.Fatal(err)
					}
					if ok {
						rf.hearbeatRe <- true
					}
				}()
			}
		}
		for i := 0; i < raftCount; i++ {
			select {
			case ok := <-rf.hearbeatRe:
				if ok {
					success_count++
					if success_count > raftCount/2 {
						fmt.Println("投票选举成功，校验心跳成功")
						fmt.Println("程序结束")
					}
				}
			}
		}
	}
}

func (rf *Raft) setTerm(term int) {
	rf.currentTerm = term
}

//设置节点的选举
func (rf *Raft) election() {
	var result bool = false
	for {
		timeout := randRange(150, 300)
		select {
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			fmt.Println("当前节点状态为: ", rf.state)
		}
		for !result {
			result = rf.election_one_rand(&leader)
		}
	}
}

func randRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

//将节点的状态修改为candidate的状态
func (rf *Raft) becomeCandidate() {
	rf.state = 1
	rf.setTerm(rf.currentTerm + 1)
	rf.votedFor = rf.me
	rf.currentLeader = -1
}

//选leader
func (rf *Raft) election_one_rand(leader *Leader) bool {
	rf.mu.Lock()
	rf.becomeCandidate()
	rf.mu.Unlock()

	fmt.Println("start electing leader")
	for {
		for i := 0; i < raftCount; i++ {
			if i != rf.me {
				go func() {
					if leader.LeaderId < 0 {
						rf.electCh <- true
					}
				}()
			}
		}
		vote := 0
		sucess := false
		triggerHearbeat := false
		for i := 0; i < raftCount; i++ {
			select {
			case ok := <-rf.electCh:
				if ok {
					vote++
					sucess = vote > raftCount/2
					if sucess && !triggerHearbeat {

						triggerHearbeat = true
						rf.mu.Lock()
						rf.becomeLeader()
						rf.mu.Unlock()
						rf.heartBeat <- true
						fmt.Println(rf.me, "号节点成为leader!")
						fmt.Println("leader发送心跳信号")
						return true
					}
				}
			}
		}
		if vote >= raftCount/2 || rf.currentLeader > -1 {
			break
		} else {
			select {
			case <-time.After(time.Duration(10) * time.Millisecond):
				fmt.Println("选举超时")
			}
		}
	}
	return false
}

func (rf *Raft) becomeLeader() {
	rf.state = 2
	rf.currentLeader = rf.me
}

type Param struct {
	Msg string
}

func (r *Raft) Communication(p Param, a *bool) error {
	fmt.Println(p.Msg)
	*a = true
	return nil
}
