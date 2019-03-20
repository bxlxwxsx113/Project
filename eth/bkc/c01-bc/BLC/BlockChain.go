package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
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
func CreateblockChainWithGenersisBlock(address string) *BlockChain {
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

			txCoinbase := NewCoinbaseTransaction(address)
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
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
func (bc *BlockChain) AddBlock(txs []*Transaction) error {
	// 获取已经存储的最后一个区块
	last_lock := bc.GetLastBlock()
	// 创建新区块
	newBlock := NewBlock(last_lock.Height+1, last_lock.Hash, txs)
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
		log.Panicf("get the lastest block failed! %v\n", err)
	}
	return block
}

func (bc *BlockChain) InsertBlock(block *Block) {
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			err := b.Put(block.Hash, block.Serialize())
			if nil != err {
				log.Panicf("update the new block to db failed! %v\n", err)
			}
			// 更新最新区块哈希
			err = b.Put([]byte(lastestHash), block.Hash)
			if nil != err {
				log.Panicf("update the lastest hash to db failed! %v\n", err)
			}
		}
		return nil
	})
	bc.Tip = block.Hash
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

// 实现挖矿
// 通过接受指定交易，进行打包确认，最总生成新的区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	var txs []*Transaction
	value, _ := strconv.Atoi(amount[0])
	// 生成新的交易
	tx := NewSimpleTransaction(from[0], to[0], value, bc)
	txs = append(txs, tx)
	// 从数据库获取最新的区块
	block := bc.GetLastBlock()
	// 生成新的区块
	block = NewBlock(block.Height+1, block.Hash, txs)
	bc.InsertBlock(block)
}

//查找指定地址的已花费输出
func (bc *BlockChain) SpentOutPut(address string) map[string][]int {
	bcit := bc.Iterator()
	// 已花费的输出缓存
	spentTXOutputs := make(map[string][]int)
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			// 排除coinbase交易
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					if in.UnLockWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						// 添加到以花费输出的缓存
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}
		}

		// 退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return spentTXOutputs
}

// 获取指定地址的UTO
func (bc *BlockChain) UnUTXOS(address string) []*UTXO {
	var unUTXOS []*UTXO
	// 1.遍历区块链查找与address相关的所有交易
	//获取迭代器对象
	bcit := bc.Iterator()
	// 2.遍历交易中的每笔交易的输出列表
	// 3.查找已花费输出
	// key: 每个input索引用的交易哈希
	// value: 索引用的输出的索引列表
	spentTXOutputs := bc.SpentOutPut(address)
	// 4.遍历每个输出与花费输出列表进行对比
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			// 查找输出
		work:
			for index, vout := range tx.Vouts {
				if vout.UnLockScriptPubkeyWithAddress(address) {
					// 判断当前vout是否存在于sentTXOutput缓存中
					// 状态变量，通过该变量该output是否已经被花费掉
					if len(spentTXOutputs) != 0 {
						var isSpentTXoutputs bool
						for txHash, indexArray := range spentTXOutputs {
							for _, i := range indexArray {
								if i == index && txHash == hex.EncodeToString(tx.TxHash) {
									isSpentTXoutputs = true
									continue work
								}
							}
						}
						// 如果说spentTXOutputs都遍历完成后，仍然没有修改状态变量
						// 则说明到那个亲Vout不存在于spentTXOutputs中
						if !isSpentTXoutputs {
							utxo := &UTXO{tx.TxHash, index, vout}
							unUTXOS = append(unUTXOS, utxo)
						}
					} else {
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
						break
					}
				}
			}
		}
		// 退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return unUTXOS
}

// 查询指定地址的余额
func (bc *BlockChain) getBalance(address string) int64 {
	// 查找指定地址的所有UTXO
	utoxs := bc.UnUTXOS(address)
	var amount int64
	for _, utxo := range utoxs {
		// 获取每个utxo的value
		amount += utxo.OutPut.Value
	}
	return amount
}

// 转账时长找找可用的UTXO就返回
func (bc *BlockChain) FindSpendableUTXO(from string, amount int64) (int64, map[string][]int) {
	spendableUTXO := make(map[string][]int)
	var value int64
	utxo := bc.UnUTXOS(from)
	for _, utxo := range utxo {
		value += utxo.OutPut.Value
		// 计算哈希
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}
	//资金不足的情况
	if value < amount {
		fmt.Printf("%s 余额不足\n", from)
		os.Exit(1)
	}
	return value, spendableUTXO
}
