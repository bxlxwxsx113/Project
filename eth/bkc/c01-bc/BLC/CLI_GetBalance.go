package BLC

import (
	"fmt"
	"os"
)

// 查询余额
func (cli *CLI) getBalance(from string) {
	if !dbExist() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	// 获取区块链对象
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()
	amount := blockchain.getBalance(from)
	fmt.Printf("\t地址 [%s] 余额 : %d\n", from, amount)
}
