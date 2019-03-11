package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

// 用法展示
func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	// 初始化区块链
	fmt.Printf("\tcreateblockchain -- 创建区块链\n")
	// 添加区块
	fmt.Printf("\taddblock -data DATA -- 添加区块\n")
	// 打印完整的区块信息
	fmt.Printf("\tprintchain -- 输出区块链信息\n")
}

// 校验参数数量
func (cli *CLI) IsValidArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

// 初始化区块链
func (cli *CLI) createBlockChain() {
	CreateblockChainWithGenersisBlock()
}

// 添加区块
func (cli *CLI) addBlock(data string) {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	blockchain.AddBlock([]byte(data))
}

// 打印完整的区块信息
func (cli *CLI) printChain() {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	blockchain.PrintChain()
}

// 启动命令运行
func (cli *CLI) Run() {
	// 检测参数数量
	cli.IsValidArgs()
	// 新建相关命令
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBLCWithGenesisBlcokCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	// 命令判断
	flagAddBlockArg := addBlockCmd.String("data", "send 100 btc to player", "添加区块")
	switch os.Args[1] {
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse addBlock failed! %v", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse printchain failed! %v", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisBlcokCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse createblockchain failed! %v", err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	// 添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			cli.PrintUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockArg)
	}
	// 输出区块链信息
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	// 创建区块链
	if createBLCWithGenesisBlcokCmd.Parsed() {
		cli.createBlockChain()
	}
}
