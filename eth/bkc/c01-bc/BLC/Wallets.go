package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

// 钱包集合文件
const walletFile = "Wallets.bat"

// 钱包集合管理
type Wallets struct {
	Wallets map[string]*Wallet
}

// 初始化钱包集合
func NewWallets() *Wallets {
	// 先从文件中获取钱包信息
	// 1.判断文件是否存在
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets
	}
	// 2.文件存在读取内容
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panicf("read the file content failed! %v\n", err)
	}
	var wallets *Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panicf("decod the file content failed! %v\n", err)
	}
	return wallets
}

//将新生成的钱包加入到集合
func (wallets Wallets) CreateWallet() {
	// 1.创建钱包
	wallet := NewWallet()
	// 2.添加
	wallets.Wallets[string(wallet.GetAddress())] = wallet
	// 3.保存到文件
	wallets.SaveWallets()
}

// 持久化钱包信息
func (wallets Wallets) SaveWallets() {
	var content bytes.Buffer
	// 注册，使用regesiter后，可以直接对内部的指定的接口进行编码
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&wallets)
	if err != nil {
		log.Panicf("encode ths struct of wallets failed! %v\n", err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panicf("write the content of wallets failed! %v\n", err)
	}
}
