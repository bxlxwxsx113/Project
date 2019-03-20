package BLC

// utxo结构管理
type UTXO struct {
	// UTXO所对应的哈希
	TxHash []byte
	// UTXO在其所属交易输出列表的索引
	Index int
	// Output
	OutPut *TxOutput
}
