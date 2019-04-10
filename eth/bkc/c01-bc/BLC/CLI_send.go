package BLC

import (
	"fmt"
	"os"
)

// 发起交易
func (cli *CLI) send(from, to, amount []string) {
	if !dbExist() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	// 获取区块链对象
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()
	// 通过传入的交易生成新的区块
	blockchain.MineNewBlock(from, to, amount)

	utxoSet := &UTXOSet{BlockChain: blockchain}
	utxoSet.Update()
}
