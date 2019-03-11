package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 共识管理相关代码
// 目标难度值
// 代表生成的hash值需要前targetBit位为0,才能满足条件
const targetBit = 8

// 工作量证明结构
type ProofOfWork struct {
	// 对指定区块进行验证
	Block *Block
	// 目标难度哈希数值
	target *big.Int
}

// 创建一个POW对象
func NewProofWork(block *Block) *ProofOfWork {
	// 数据总长度为8位（代表原数据为256位）
	// 假设需要前两位为0， 才能满足解题条件
	// 8 - 2 = 6
	// 左移一位代表乘以2
	// a << n = a * 2^n
	// 1*2^6=64
	// 0100 0000
	// 0011 1111
	// 所以只要小于0100 0000就是满足前两位为0
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{block, target}
}

// 开始工作量证明，比较哈希
func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
	var nonce = 0       // 碰撞次数
	var hash [32]byte   // 生成的哈希
	var hashInt big.Int // 存储哈希转换之后的大数
	for {
		//生 成hash值
		dataBytes := proofOfWork.prepareData(nonce)
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		// 比较和目标生成的hash是否达标
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			// 找到了符合条件的哈希
			break
		}
		nonce++
	}
	fmt.Printf("\n碰撞次数:%d\n", nonce)
	fmt.Printf("\rhash:%x\n", hash)
	return hash[:], int64(nonce)
}

// 数据拼接函数
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	var data []byte
	//拼接所有区块属性，进行hash计算
	data = bytes.Join([][]byte{
		IntToHex(pow.Block.TimeStamp),
		IntToHex(pow.Block.Height),
		pow.Block.PreBlockHash,
		pow.Block.Data,
		IntToHex(int64(nonce)),
		IntToHex(targetBit),
	}, []byte{})
	return data
}
