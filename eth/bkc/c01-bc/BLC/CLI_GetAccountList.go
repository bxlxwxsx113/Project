package BLC

import "fmt"

// 实现通过命令获取地址列表
func (cli *CLI) GetAccount() {
	wallets := NewWallets()
	fmt.Println("帐号别表...")
	for key, _ := range wallets.Wallets {
		fmt.Printf("\t[%s]\n", key)
	}
}
