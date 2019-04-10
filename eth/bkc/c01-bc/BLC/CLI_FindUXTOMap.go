package BLC

import "fmt"

func (cli *CLI) TestFIndUTXOMap() {
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()
	utxoMap := blockchain.FindUTXOMap()
	for key, value := range utxoMap {
		fmt.Printf("key : %v\n", key)
		for _, out := range value.TXoutputs {
			fmt.Printf("pubKey:%x, value : %d\n", out.Ripemd160Hash, out.Value)
		}
	}
	fmt.Printf("utxoMap:%v\n", utxoMap)
}

// 重置UTXO table
func (cli *CLI) TestResetUTXO() {
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()
	utxoSet := UTXOSet{BlockChain: blockchain}
	utxoSet.RestUTXOSet()
}
