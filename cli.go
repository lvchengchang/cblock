package main

import (
	"fmt"
	"os"
	"strconv"
)

type Cli struct {
	bc *BlockChain
}

const USAGE = `
	addBlock --data DATA "add data to blockchain"
	printChain           "print all blockchain data"
	getBalance --address ADDRESS "获取指定地址余额"
	send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
`

func (cli *Cli) Run() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf(USAGE)
		return
	}

	cmd := args[1]
	switch cmd {
	case "addBlock":
		if len(args) == 4 && args[2] == "--data" {
			data := args[3]
			cli.AddBlock(data)
		} else {
			fmt.Println("添加区块参数失败,请重试")
		}
	case "printChain":
		cli.PrintBlockChain()
	case "getBalance":
		fmt.Println("获取余额")
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Printf("转账开始...\n")
		if len(args) != 7 {
			fmt.Printf("参数个数错误，请检查！\n")
			fmt.Printf(USAGE)
			return
		}
		//./block send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64) //知识点，请注意
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	default:
		fmt.Printf(USAGE)
	}
}

func (cli *Cli) GetBalance(address string) {
	utxos := cli.bc.FindUTXOs(address)

	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("%s 的余额是 %f", address, total)
}

func (cli *Cli) Send(from, to string, amount float64, miner, data string) {
	// 生成一个挖矿数据
	coinbase := NewCoinBaseTx(miner, data)
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		return
	}

	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	fmt.Println("转账成功")
}
