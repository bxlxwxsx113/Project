package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
)

// 数据库名称
const dbName = "block7.db"

// 表名称
const blockTableName = "blocks"

const lastestHash = "lastestHash"

// 区块链基本结构
type BlockChain struct {
	// 数据库对象
	DB *bolt.DB
	// 最新区块哈希
	Tip []byte
}

// 判断数据库文件是否存在
func dbExist() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		// 数据库文件不存在
		return false
	}
	return true
}

//初始化区块链
func CreateblockChainWithGenersisBlock() *BlockChain {
	if dbExist() {
		fmt.Println("创世区块已经存在...")
		os.Exit(1)
	}
	var blockHash []byte
	// 创建或打开一个数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	// 生成创世区块
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil == b {
			// 没找到表
			bucket, err := tx.CreateBucket([]byte(blockTableName))
			if nil != err {
				log.Panicf("create the bucket [%s] failed! %v\n", blockTableName, err)
			}
			genesisBlock := CreateGenesisBlock("this init of the blockChain")
			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if nil != err {
				log.Panicf("insert the genesis block to db failed! %v\n", err)
			}
			err = bucket.Put([]byte(lastestHash), genesisBlock.Hash)
			if nil != err {
				log.Panicf("insert the lastest block hash to db failed! %v\n", err)
			}
			blockHash = genesisBlock.Hash
		} else {
			fmt.Println(" find")
		}
		return nil
	})
	if err != nil {
		log.Panicf("insert the block to db failed! %v\n", err)
		return nil
	}
	return &BlockChain{db, blockHash}
}

//添加区块到区块链
func (bc *BlockChain) AddBlock(data []byte) error {
	// 获取已经存储的最后一个区块
	last_lock := bc.GetLastBlock()
	// 创建新区块
	newBlock := NewBlock(last_lock.Height+1, last_lock.Hash, data)
	// 更新数据
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockTableName))
		// 更新最后一个区块hash
		err := bucket.Put([]byte(lastestHash), newBlock.Hash)
		if nil != err {
			log.Panicf("insert the last hash to db failed! %v\n", err)
		}
		// 更新最后一个区块
		err = bucket.Put(newBlock.Hash, newBlock.Serialize())
		if nil != err {
			log.Panicf("insert the last block to db failed! %v\n", err)
		}
		bc.Tip = newBlock.Hash
		return nil
	})
	if nil != err {
		log.Panicf("update the db failed! %v\n", err)
	}
	return nil
}

//  获取最后一个区块
func (bc *BlockChain) GetLastBlock() *Block {
	var block *Block
	err := bc.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockTableName))
		if bucket == nil {
			log.Panicf("get the bucket failed! \n")
		}
		lastHash := bucket.Get([]byte(lastestHash))
		fmt.Printf("lastHash = %x\n", lastHash)
		if nil == lastHash {
			log.Panicf("get the last hash failed! \n")
		}
		blockBytes := bucket.Get(lastHash)
		fmt.Printf("blockBytes = %x\n", blockBytes)
		if blockBytes == nil {
			log.Panicf("get the lastest block failed! \n")
		}
		fmt.Println(blockBytes)
		block = DeserializeBlock(blockBytes)
		return nil
	})
	if err != nil {
		//log.Panicf("get the lastest block failed! %v\n", err)
	}
	return block
}

// 遍历数据库，输出所有区块信息

func (bc *BlockChain) PrintChain() {
	fmt.Println("区块链完整信息...")
	var curBlock *Block
	bcit := bc.Iterator()
	for {
		curBlock = bcit.Next()
		curBlock.Print()
		// 什么时候退出
		// 判断创世区块前hash是否位空
		var hashInt big.Int
		hashInt.SetBytes(curBlock.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

// 迭代器对象
func (blc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{blc.DB, blc.Tip}
}

// 返回一个blockchain对象
func BlockChainObject() *BlockChain {
	// 读取数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panicf("open the db file failed!%v\n", err)
	}
	// 获取最新区块哈希
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			tip = b.Get([]byte(lastestHash))
		}
		return nil
	})
	if err != nil {
		log.Panicf("get the last hash failed!%v\n", err)
	}
	return &BlockChain{db, tip}
}
