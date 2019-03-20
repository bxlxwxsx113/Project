package BLC

type TxInput struct {
	// 交易的哈希（不是指当前交易的哈希，而是该输入所引用的交易的哈希）
	TxHash []byte
	// 引用的上衣比交易的索引号
	Vout int
	// 用户名
	ScriptSig string
}

// 权限判断，检查引用的输出是否属于传入的地址
func (in *TxInput) UnLockWithAddress(address string) bool {
	return in.ScriptSig == address
}
