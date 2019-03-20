package BLC

type TxOutput struct {
	// 金额
	Value int64
	// 用户名（钱时水的，UTXO所有者）
	ScriptPubkey string
}

// 身份验证
func (out *TxOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	return out.ScriptPubkey == address
}
