package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

// 实现一个最基本的区块结构
type Block struct {
	// 区块时间戳，代表区块产生时间
	TimeStamp int64
	// 当前区块hash
	Hash []byte
	// 前一区块hash
	PreBlockHash []byte
	// 交易数据
	Data []byte
	// 区块高度，代表当前区块索引，也表示区块链中的区块数量
	Height int64
	// 在运行pow时生成hash的变化次数。也就是代表pow运算时的动态值
	Nonce int64
}

// 生产新的区块
func NewBlock(height int64, prevBlockHash []byte, data []byte) *Block {
	var block *Block
	block = &Block{
		TimeStamp:    time.Now().Unix(),
		Hash:         nil,
		PreBlockHash: prevBlockHash,
		Data:         data,
		Height:       height,
	}
	// block.SetHash()
	pow := NewProofWork(block)
	// 执行工作量证明
	block.Hash, block.Nonce = pow.Run()
	return block
}

// 计算区块hash
func (b *Block) SetHash() {
	timeStampBytes := IntToHex(b.TimeStamp)
	heightBytes := IntToHex(b.Height)
	// 拼接所有区块属性，进行hash计算
	blockBytes := bytes.Join([][]byte{
		timeStampBytes,
		heightBytes,
		b.PreBlockHash,
		b.Data,
	}, []byte{})
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}

// 生成创世区块
func CreateGenesisBlock(data string) *Block {
	fmt.Println("12")
	return NewBlock(1, nil, []byte(data))
}

// 区块结构序列化
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result) // 新建encode对象
	if err := encoder.Encode(b); nil != err {
		log.Panicf("serialize the block to []byte failed! &v\n", err)
	}
	return result.Bytes()
}

// 反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&block); nil != err {
		log.Panicf("deserialize the []byte to block failed! &v\n", err)
	}
	return &block
}

func (b *Block) Print() {
	fmt.Printf("\tHash:%x\n", b.Hash)
	fmt.Printf("\tPreBlockHash:%x\n", b.PreBlockHash)
	fmt.Printf("\tTimeStamp:%v\n", b.TimeStamp)
	fmt.Printf("\tData:%v\n", b.Data)
	fmt.Printf("\tHeight:%d\n", b.Height)
	fmt.Printf("\tNonce:%d\n\n", b.Nonce)
}
