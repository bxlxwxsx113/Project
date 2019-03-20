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
	fmt.Printf("\tcreateblockchain -address address-- 创建区块链\n")
	// 添加区块
	fmt.Printf("\taddblock -data DATA -- 添加区块\n")
	// 打印完整的区块信息
	fmt.Printf("\tprintchain -- 输出区块链信息\n")
	// 转账
	fmt.Printf("\tsend -from FROM -to To -amount AMOUNT -- 转账\n")
	// 查询余额
	fmt.Println("\tgetbalance -address FROM -- 查询指定地址余额")

	fmt.Println("\t转账参数说明：")
	fmt.Println("\t\t-from FROM -- 转账源地址")
	fmt.Println("\t\t-to TO -- 转账源地址")
	fmt.Println("\t\t-amout AMOUT -- 转账源地址")

	fmt.Println("\t\t-address -- 要查询的地址")
}

// 校验参数数量
func (cli *CLI) IsValidArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

// 添加区块
func (cli *CLI) addBlock(txs []*Transaction) {
	if !dbExist() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	blockchain.AddBlock(txs)
}

// 启动命令运行
func (cli *CLI) Run() {
	// 检测参数数量
	cli.IsValidArgs()
	// 新建相关命令
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBLCWithGenesisBlcokCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("etbalance", flag.ExitOnError)
	// 命令判断
	flagAddBlockArg := addBlockCmd.String("data", "send 100 btc to player", "添加区块")
	flagCreateBlockChain := createBLCWithGenesisBlcokCmd.String("address", "", "接受创世奖励")
	flagSendFrom := sendCmd.String("from", "", "源地址")
	flagSendTo := sendCmd.String("to", "", "目标地址")
	flagSendAmount := sendCmd.String("amount", "", "转账数额")
	flagBalanceAddress := getBalanceCmd.String("address", "", "查询的地址")
	switch os.Args[1] {
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse Balance failed! %v\n", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse transaction failed! %v\n", err)
		}
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse addBlock failed! %v\n", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse printchain failed! %v\n", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisBlcokCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse createblockchain failed! %v\n", err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	// 查询余额
	if getBalanceCmd.Parsed() {
		if *flagBalanceAddress == "" {
			fmt.Println("请输入查询地址")
			os.Exit(1)
		}
		cli.getBalance(*flagBalanceAddress)
	}

	// ./bc.exe send -from "[\"A\"]" -to "[\"C\"]" -amount "[\"10\"]"
	// 发起交易
	if sendCmd.Parsed() {
		if *flagSendFrom == "" {
			fmt.Println("源地址不能为空")
			cli.PrintUsage()
			os.Exit(1)
		}
		if *flagSendTo == "" {
			fmt.Println("目标地址不能为空")
			cli.PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmount == "" {
			fmt.Println("转账金额不能为空")
			cli.PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("FROM:[%s] TO:[%s] VALUE:[%s]\n",
			JSONtoSlice(*flagSendFrom),
			JSONtoSlice(*flagSendTo),
			JSONtoSlice(*flagSendAmount))
		cli.send(JSONtoSlice(*flagSendFrom), JSONtoSlice(*flagSendTo), JSONtoSlice(*flagSendAmount))
	}

	// 添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			cli.PrintUsage()
			os.Exit(1)
		}
		cli.addBlock([]*Transaction{})
	}
	// 输出区块链信息
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	// 创建区块链
	if createBLCWithGenesisBlcokCmd.Parsed() {
		if *flagCreateBlockChain == "" {
			cli.PrintUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateBlockChain)
	}
}
