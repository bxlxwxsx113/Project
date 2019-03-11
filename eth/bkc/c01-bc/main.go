package main

import "Project/eth/bkc/c01-bc/BLC"

//启动入口
func main() {
	/*bc := BLC.CreateblockChainWithGenersisBlock()
	bc.AddBlock([]byte("123"))
	bc.AddBlock([]byte("456"))
	bc.AddBlock([]byte("789"))
	bc.PrintChain()*/

	cli := BLC.CLI{}
	cli.Run()
}
