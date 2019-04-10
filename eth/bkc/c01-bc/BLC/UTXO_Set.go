package BLC

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// UTXO持久化管理
const utxoTableName = "utxoTable"

// 生成UTXOSet结构（保存指定的区块所有的 UTXO）
type UTXOSet struct {
	BlockChain *BlockChain
}

// 将UTXO进行序列化
func (txOutputs *TXOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(txOutputs); nil != err {
		log.Panicf("serialize the utxo table failed! %v\n", err)
	}
	return result.Bytes()
}

// 反序列化
func DeserializeTXOuputs(txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panicf("deserilize UTXO table failed! %v\n", err)
	}
	return &txOutputs
}

// 重置
func (utxoSet *UTXOSet) RestUTXOSet() {
	// 在第一次创建的时候更新utxo table
	utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if nil != err {
				log.Panicf("delete the utxo table error! &v\n", err)
			}
		}
		bucket, _ := tx.CreateBucket([]byte(utxoTableName))
		if nil != bucket {
			// 查找所有UTXO
			txOutputsMap := utxoSet.BlockChain.FindUTXOMap()
			for keyHash, outputs := range txOutputsMap {
				// 存入utxo table
				txHash, _ := hex.DecodeString(keyHash)
				err := bucket.Put(txHash, outputs.Serialize())
				if nil != err {
					log.Panicf("put the utxo into table failed! %v\n", err)
				}
			}
		}
		return nil
	})
}

// 查找
func (utxoSet *UTXOSet) GetBalance(address string) int64 {

	UTXOS := utxoSet.FindUTXOWithAddress(address)
	var amount int64
	for _, utxo := range UTXOS {
		fmt.Printf("utxo-txHash:%x\n", utxo.TxHash)
		fmt.Printf("utxo-Index:%d\n", utxo.Index)
		fmt.Printf("utxo-Ripemd160Hash:%x\n", utxo.OutPut.Ripemd160Hash)
		fmt.Printf("utxo-Value:%d\n", utxo.OutPut.Value)
		amount += utxo.OutPut.Value
	}
	return amount
}

// 通过utxotale找指定地址的UTXO
func (utxoSet *UTXOSet) FindUTXOWithAddress(address string) []*UTXO {
	var utxos []*UTXO
	// 查找数据库表
	utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			// 游标
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("k : %x\n", k)
				txOutputs := DeserializeTXOuputs(v)
				for _, utxo := range txOutputs.TXoutputs {
					if utxo.UnLockScriptPubkeyWithAddress(address) {
						utxo_signle := UTXO{OutPut: utxo}
						utxos = append(utxos, &utxo_signle)
					}
				}
			}
		}
		return nil
	})
	return utxos
}

// 更新
func (utxoSet *UTXOSet) Update() {
	// 获取最新区块
	lastest_block := utxoSet.BlockChain.Iterator().Next()
	db := utxoSet.BlockChain.DB
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			for _, tx := range lastest_block.Txs {
				if !tx.IsCoinbaseTransaction() {
					// 遍历交易输入
					for _, vin := range tx.Vins {
						// 需要更新的输出
						updateoutputs := TXOutputs{}
						outputBytes := b.Get(vin.TxHash)
						outs := DeserializeTXOuputs(outputBytes)
						for outIdx, out := range outs.TXoutputs {
							if vin.Vout != outIdx {
								updateoutputs.TXoutputs = append(updateoutputs.TXoutputs, out)
							}
						}
						//  如果交易没有需要更新的UTXO，删除交易
						if len(updateoutputs.TXoutputs) == 0 {
							b.Delete(vin.TxHash)
						} else {
							// 存入数据库
							b.Put(vin.TxHash, updateoutputs.Serialize())
						}

					}
				}
				newOutputs := TXOutputs{}
				newOutputs.TXoutputs = append(newOutputs.TXoutputs, tx.Vouts...)
				b.Put(tx.TxHash, newOutputs.Serialize())
			}
		}
		return nil
	})
	if err != nil {
		log.Panicf("update utxo failed! %v\n", err)
	}
}
