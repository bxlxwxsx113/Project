package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"log"
)

// 版本
const version = byte(0x00)

// 校验和长度
const addresscCheckLength = 4

type Wallet struct {
	// 私钥
	PrivateKey ecdsa.PrivateKey
	// 公钥
	PublicKey []byte
}

// 创建一个钱包
func NewWallet() *Wallet {
	privateKey, pubKey := newKeyPair()
	return &Wallet{privateKey, pubKey}
}

// 生成私钥-公钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 得到一个椭圆
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicf("ecdsa generate key failed!%v\n", err)
	}
	// 生成公钥
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, pubKey
}

// 实现双哈希
func Ripemd160Hash(pubKey []byte) []byte {
	// 1.sha256
	hash256 := sha256.New()
	hash256.Write(pubKey)
	hash := hash256.Sum(nil)
	// 2.ripmd160
	rmd160 := ripemd160.New()
	rmd160.Write(hash)
	return rmd160.Sum(nil)
}

// 通过钱包获取地址
func (w *Wallet) GetAddress() []byte {
	// 1.获取ripemd160结果
	ripedmd160 := Ripemd160Hash(w.PublicKey)
	// 2.获取到version,  加入到hash中
	//version_ripemd160Hash := append([]byte{version}, ripedmd160...)
	version_ripemd160Hash := ripedmd160
	// 3.生成校验和
	checkSumBytes := CheckSum(version_ripemd160Hash)
	// 4.拼接校验和
	addressBytes := append(version_ripemd160Hash, checkSumBytes...)
	// 5.base58编码
	b58Bytes := Base58Encode(addressBytes)
	fmt.Printf("%s\n", b58Bytes)
	return b58Bytes
}

// 生成校验和
func CheckSum(payload []byte) []byte {
	first_hsah := sha256.Sum256(payload)
	second_hash := sha256.Sum256(first_hsah[:])
	return second_hash[:addresscCheckLength]
}

// 判断地址有效性
func IsValidForAddress(address []byte) bool {
	// 1.通过base58解码
	version_pubkey_checkSumByte := Base58Decode(address)
	// 2.拆分，进行校验和的校验
	checkSumBytes := version_pubkey_checkSumByte[len(version_pubkey_checkSumByte)-addresscCheckLength:]
	version_ripemd160 := version_pubkey_checkSumByte[:len(version_pubkey_checkSumByte)-addresscCheckLength]
	// 3.比较
	checkBytes := CheckSum(version_ripemd160)
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}
