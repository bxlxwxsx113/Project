package BLC

import "fmt"

func (cli *CLI) CreateWallets() {
	// 创建一个钱包集合
	// 在钱包文件已经存在的情况下，现获取数据
	wallets := NewWallets()
	wallets.CreateWallet()
	fmt.Printf("wallets ： %s\n", wallets)
}
