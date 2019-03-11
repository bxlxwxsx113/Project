package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// 迭代器管理
type BlockChainIterator struct {
	db          *bolt.DB // 数据库对象
	CurrentHash []byte   //当前区块hash
}

// 实现迭代函数
func (bcit *BlockChainIterator) Next() *Block {
	var block *Block
	err := bcit.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//  获取指定hash的区块
			curentBlockBytes := b.Get(bcit.CurrentHash)
			block = DeserializeBlock(curentBlockBytes)
			//  更新迭代器的hash
			bcit.CurrentHash = block.PreBlockHash
		}
		return nil
	})
	if nil != err {
		log.Panicf("iterator the db dailed!%v\n", err)
	}
	return block
}
