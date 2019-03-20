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
	//Data []byte
	Txs []*Transaction
	// 区块高度，代表当前区块索引，也表示区块链中的区块数量
	Height int64
	// 在运行pow时生成hash的变化次数。也就是代表pow运算时的动态值
	Nonce int64
}

// 生产新的区块
func NewBlock(height int64, prevBlockHash []byte, txs []*Transaction) *Block {
	var block *Block
	block = &Block{
		TimeStamp:    time.Now().Unix(),
		Hash:         nil,
		PreBlockHash: prevBlockHash,
		Txs:          txs,
		Height:       height,
	}
	// block.SetHash()
	pow := NewProofWork(block)
	// 执行工作量证明
	block.Hash, block.Nonce = pow.Run()
	return block
}

// 计算区块hash
/*func (b *Block) SetHash() {
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
}*/

// 生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
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
	fmt.Printf("\tTimeStamp:%v\n", time.Unix(b.TimeStamp, 0))
	fmt.Printf("\tHeight:%d\n", b.Height)
	fmt.Printf("\tNonce:%d\n\n", b.Nonce)
	for _, tx := range b.Txs {
		fmt.Printf("\t\ttx-Hash : %x\n", tx.TxHash)
		fmt.Printf("\t\t输入...\n")
		for _, vin := range tx.Vins {
			fmt.Printf("\t\t\tvin-txHash : %x\n", vin.TxHash)
			fmt.Printf("\t\t\tvin-Vout : %x\n", vin.Vout)
			fmt.Printf("\t\t\tvin-ScriptSig : %s\n", vin.ScriptSig)
		}
		fmt.Printf("\t\t输出...\n")
		for _, vout := range tx.Vouts {
			fmt.Printf("\t\t\tvout-Value : %d\n", vout.Value)
			fmt.Printf("\t\t\tvout-ScriptPubkey : %s\n", vout.ScriptPubkey)
		}
	}
	fmt.Println("\t-------------------------------------------------------------------------")
}

// 把指定区块的所有交易序列化
func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte
	//  将区块中的所有交易哈希进行拼接
	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	// 将去块中的所有交易拼接后生新的哈希
	txsHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txsHash[:]
}
