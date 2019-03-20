package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	// 交易哈希，交易的唯一标识
	TxHash []byte
	// 输入列表
	Vins []*TxInput
	// 输出列表
	Vouts []*TxOutput
}

// 生成coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {
	//  输入
	txInput := &TxInput{[]byte{}, -1, "system reward"}
	// 输出
	txOutput := &TxOutput{10, address}
	txCoinbase := &Transaction{nil, []*TxInput{txInput}, []*TxOutput{txOutput}}
	txCoinbase.HashTransaction()
	return txCoinbase
}

// 生成交易哈希
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	endcode := gob.NewEncoder(&result)
	if err := endcode.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed! %v\n", err)
	}
	// 生成哈希值
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

// 生成普通转账交易
func NewSimpleTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput
	money, spendableUTXODic := bc.FindSpendableUTXO(from, int64(amount))
	fmt.Println("money : %v\n", money)
	for txHash, indexArray := range spendableUTXODic {
		txHashbytes, err := hex.DecodeString(txHash)
		if nil != err {
			log.Panicf("decode string failed!")
		}
		for _, index := range indexArray {
			txInput := &TxInput{txHashbytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}
	// 输出
	txOutput := &TxOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)
	// 输出（找零）
	txOutput = &TxOutput{money - int64(amount), from}
	txOutputs = append(txOutputs, txOutput)
	tx := &Transaction{nil, txInputs, txOutputs}
	// 生成新的交易哈希
	tx.HashTransaction()
	return tx
}

// 判断指定交易是否为coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}
