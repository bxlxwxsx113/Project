package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"
)

type Transaction struct {
	// 交易哈希，交易的唯一标识
	TxHash []byte
	// 输入列表
	Vins []*TxInput
	// 输出列表
	Vouts []*TxOutput
}

// 生成coinbase交易
func NewCoinbaseTransaction(address string) *Transaction {
	//  输入
	txInput := &TxInput{[]byte{}, -1, nil, nil}
	// 输出
	txOutput := NewTxOutput(10, address)
	txCoinbase := &Transaction{nil, []*TxInput{txInput}, []*TxOutput{txOutput}}
	txCoinbase.HashTransaction()
	return txCoinbase
}

// 生成交易哈希
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	endcode := gob.NewEncoder(&result)
	if err := endcode.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed! %v\n", err)
	}
	// 添加时间戳标识，否则会导致所有coinbase交易哈希一样
	tm := time.Now().Unix()
	txHashBytes := bytes.Join([][]byte{result.Bytes(), IntToHex(tm)}, []byte{})
	// 生成哈希值
	hash := sha256.Sum256(txHashBytes)
	tx.TxHash = hash[:]
}

// 生成普通转账交易
func NewSimpleTransaction(from, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput

	money, spendableUTXODic := bc.FindSpendableUTXO(from, int64(amount), txs)
	fmt.Println("money : %v\n", money)
	// 获取到钱包集合
	wallets := NewWallets()
	// 查找对应的钱包结构
	wallet, ok := wallets.Wallets[from]
	if !ok {
		log.Panicf("get address [%s] failed! \n", from)
	}
	for txHash, indexArray := range spendableUTXODic {
		txHashbytes, err := hex.DecodeString(txHash)
		if nil != err {
			log.Panicf("decode string failed!")
		}
		for _, index := range indexArray {
			txInput := &TxInput{txHashbytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}
	// 输出
	txOutput := NewTxOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)
	// 输出（找零）
	txOutput = NewTxOutput(money-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)
	tx := &Transaction{nil, txInputs, txOutputs}
	// 生成新的交易哈希
	tx.HashTransaction()
	// 交易签名
	bc.SignTransaction(wallet.PrivateKey, tx)
	return tx
}

// 判断指定交易是否为coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

func (tx *Transaction) Sign(priavteKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	// 处理输入
	for _, vin := range tx.Vins {
		// 输入引用的交易
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Printf("ERROR:Prev transaction is not corret!\n")
		}
	}
	// 提取需要签名的属性
	txCopy := tx.TrimmedCopy()
	for vin_id, vin := range txCopy.Vins {
		// 获取关联交易
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// 找到发送者(当前输入引用的哈希就是其索引用的输出的哈希值)
		txCopy.Vins[vin_id].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		// 调用核心签名函数
		r, s, err := ecdsa.Sign(rand.Reader, &priavteKey, txCopy.TxHash)
		if err != nil {
			log.Panicf("sign to tx [%x] failed!%v\n", txCopy.TxHash, err)
		}
		// 组成签名
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[vin_id].Signature = signature
	}
	//
}

// 交易拷贝，生成一个专用签名的副本
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput
	for _, vin := range tx.Vins {
		inputs = append(inputs, &TxInput{vin.TxHash, vin.Vout, nil, nil})
	}
	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TxOutput{vout.Value, vout.Ripemd160Hash})
	}
	txCopy := Transaction{tx.TxHash, inputs, outputs}
	return txCopy
}

// 设置用于签名交易的哈希
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// 交易序列化
func (tx *Transaction) Serialize() []byte {
	var result bytes.Buffer
	// 新建encoder对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); nil != err {
		log.Panicf("serialize the tx to byte failed!%v\n", err)
	}
	return result.Bytes()
}

//  交易签名验证
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}
	// 检查能否找到交易
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("VERIFY ERROR : transaction is not correct!\n")
		}
	}
	// 获取相同的交易副本
	txCopy := tx.TrimmedCopy()
	// 使用相同的椭圆获取密钥对
	curve := elliptic.P256()

	// 遍历tx的输入，对每比输入pinyon的输出进行验证呢个
	for vinId, vin := range tx.Vins {
		tx := prevTxs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[vinId].PublicKey = tx.Vouts[vin.Vout].Ripemd160Hash
		// 需要验证的数据
		txCopy.TxHash = txCopy.Hash()
		// 获取r，s
		// 签名时一个数字对，r和s就代表签名
		// 将signature中的r，s抽取出来
		// r,s 长度相等
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		// 获取公钥
		// 公钥是由X,Y坐标组合
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(pubKeyLen / 2)])
		y.SetBytes(vin.PublicKey[(pubKeyLen / 2):])
		// 组装成原始公钥
		rawpubKey := ecdsa.PublicKey{curve, &x, &y}
		if !ecdsa.Verify(&rawpubKey, txCopy.TxHash, &r, &s) {
			return false
		}
	}
	return true
}
