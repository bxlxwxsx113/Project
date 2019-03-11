package BLC

import (
	"bytes"
	"encoding/binary"
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
