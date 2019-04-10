package BLC

import (
	"bytes"
	"math/big"
)

var b58Alphabet = []byte("" +
	"123456789" +
	"abcdefghijkmnopqsrtuvwxyz" +
	"ABCDEFGHJKLMNPQRSTUVWXYZ" +
	"")

func Base58Encode(input []byte) []byte {
	var result []byte
	// 将byte转换为big.Int
	x := big.NewInt(0).SetBytes(input)
	// 求余的基本长度58
	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	// 余数
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}
	Reverse(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for _, b := range input {
		if b == b58Alphabet[0] {
			zeroBytes++
			break
		} else {
			break
		}
	}
	// 去掉前缀
	data := input[zeroBytes:]
	for _, b := range data {
		// 得到bytes数组中指定的数字/字符第一次出现的索引
		charIndex := bytes.IndexByte(b58Alphabet, b)

		result.Mul(result, big.NewInt(58))
		// 加上余数
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	return decoded
}
