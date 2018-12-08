package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

func PrintBlock(block *Block){
	fmt.Printf("Hash:0x%x\n", block.Hash)
	fmt.Printf("Data:%s\n", block.Data)
	fmt.Printf("version:%d\n", block.Version)
	fmt.Printf("PrevBlockHash:0x%x\n", block.PrevBlockHash)
	fmt.Printf("TimeStamp:%d\n", block.TimeStamp)
	fmt.Printf("TargetBits:0x%x\n", block.TargetBits)
	fmt.Printf("Nonce:%d\n", block.Nonce)
	fmt.Printf("MerkleRoot:0x%x\n", block.MerkleRoot)
	pow := NewPow(block)
	fmt.Printf("isValid:%v\n", pow.IsValid())
}

type CLI struct {
	bc *BlockChain
}

func NewCLI(bc *BlockChain)*CLI{
	return &CLI{bc}
}

func (cli *CLI) usage()string{
	var buffer bytes.Buffer
	_,err := fmt.Fprintf(&buffer,"Usage:\n")
	CheckError("CLI.usage #0",err)
	_,err = fmt.Fprintf(&buffer,"%s addblock --data DATA\n",os.Args[0])
	CheckError("CLI.usage #1",err)
	_,err = fmt.Fprintf(&buffer,"%s listblock\n",os.Args[0])
	CheckError("CLI.usage #2",err)
	return string(buffer.Bytes())
}

func (cli *CLI) AddBlock(data string){
	cli.bc.AddBlock(data)
}

func (cli *CLI) ListBlock() {
	iter := cli.bc.Iterator()
	var i int = 1
	for b := iter.Next();b != nil;b = iter.Next() {
		fmt.Printf("=================block:%d===================\n", i)
		PrintBlock(b)
		i++
	}
}

func (cli *CLI) run() {
	cmdAdd := flag.NewFlagSet("addblock",flag.ExitOnError)
	cmdList := flag.NewFlagSet("listblock",flag.ExitOnError)
	cmdAddDataParam := cmdAdd.String("data","","Transactions Data String")
	if len(os.Args) < 2 {
		fmt.Printf("%s\n",cli.usage())
		os.Exit(-1)
	}
	switch os.Args[1] {
	case "addblock":
		err := cmdAdd.Parse(os.Args[2:])
		CheckError("CLI.run #1",err)
		if cmdAdd.Parsed(){
			if *cmdAddDataParam != "" {
				cli.AddBlock(*cmdAddDataParam)
			}
		}
	case "listblock":
		err := cmdList.Parse(os.Args[2:])
		CheckError("CLI.run #1",err)
		if cmdList.Parsed(){
			cli.ListBlock()
		}
	default:
		fmt.Printf("%s\n",cli.usage())
		os.Exit(1)
	}
}
