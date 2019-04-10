package BLC

import (
	"bytes"
)

type TxOutput struct {
	// 金额
	Value int64
	// 用户名（钱时水的，UTXO所有者）
	//ScriptPubkey string
	Ripemd160Hash []byte
}

// 身份验证
func (out *TxOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	// 把address转换成ripemd160hash
	hash160 := stringToHash160(address)
	return bytes.Compare(out.Ripemd160Hash, hash160) == 0
}

// string转hash160
func stringToHash160(address string) []byte {
	pubKeyHash := Base58Encode([]byte(address))
	hash160 := pubKeyHash[:len(pubKeyHash)-addresscCheckLength]
	return hash160
}

// 新建output对象
func NewTxOutput(value int64, address string) *TxOutput {
	txOutput := &TxOutput{}
	hash160 := stringToHash160(address)
	txOutput.Value = value
	txOutput.Ripemd160Hash = hash160
	return txOutput
}
