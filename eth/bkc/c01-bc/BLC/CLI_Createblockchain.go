package BLC

// 初始化区块链
func (cli *CLI) createBlockChain(address string) {
	blockchain := CreateblockChainWithGenersisBlock(address)
	defer blockchain.DB.Close()
	// 设置uxto重置操作
	uxtoSet := &UTXOSet{blockchain}
	uxtoSet.RestUTXOSet()
}
