package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

// int64转换byte
func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	//根据大小端转换
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Panicf("int to []byte faild! %v\n", err)
	}
	return buffer.Bytes()
}

// 标准的JSON格式转切片
func JSONtoSlice(jsonString string) []string {
	var strSlince []string
	if err := json.Unmarshal([]byte(jsonString), &strSlince); err != nil {
		log.Printf("json to []string failed! %v\n", err)
	}
	return strSlince
}

// 反转切片
func Reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
