package BLC

import (
	"fmt"
	"os"
)

// 打印完整的区块信息
func (cli *CLI) printChain() {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	blockchain.PrintChain()
}
