package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

func PrintBlock(block *Block){
	fmt.Printf("Hash:0x%x\n", block.Hash)
	fmt.Printf("Transactions: %v\n", block.Transactions)
	fmt.Printf("version:%d\n", block.Version)
	fmt.Printf("PrevBlockHash:0x%x\n", block.PrevBlockHash)
	fmt.Printf("TimeStamp:%d\n", block.TimeStamp)
	fmt.Printf("TargetBits:0x%x\n", block.TargetBits)
	fmt.Printf("Nonce:%d\n", block.Nonce)
	fmt.Printf("MerkleRoot:0x%x\n", block.MerkleRoot)
	pow := NewPow(block)
	fmt.Printf("isValid:%v\n", pow.IsValid())
}

type CLI struct {}

func NewCLI()*CLI{
	return &CLI{}
}

func (cli *CLI) usage()string{
	var buffer bytes.Buffer
	_,err := fmt.Fprintf(&buffer,"Usage:\n")
	CheckError("CLI.usage #0",err)
	_,err = fmt.Fprintf(&buffer,"%s newchaine\n",os.Args[0])
	CheckError("CLI.usage #1",err)
	_,err = fmt.Fprintf(&buffer,"%s addblock\n",os.Args[0])
	CheckError("CLI.usage #2",err)
	_,err = fmt.Fprintf(&buffer,"%s listblock\n",os.Args[0])
	CheckError("CLI.usage #3",err)
	return string(buffer.Bytes())
}

func (cli *CLI) AddBlock(){
	bc := NewBlockChain()
	bc.AddBlock(nil)
}

func (cli *CLI) ListBlock() {
	bc := NewBlockChain()
	iter := bc.Iterator()
	var i int = 1
	for b := iter.Next();b != nil;b = iter.Next() {
		fmt.Printf("=================block:%d===================\n", i)
		PrintBlock(b)
		i++
	}
}

func (cli *CLI) NewChain(){
	var lastHash []byte
	db,err := bolt.Open(DataBaseFile,0600,nil)
	CheckError("CLI.NewChain #1",err)
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil{
			lastHash = bucket.Get([]byte(lastHashKey))
		} else {
			transx := NewCoinbaseTx("",GenesisBlockInfo)
			block := NewGenesisBlock(transx)
			bucket, err = tx.CreateBucket([]byte(bucketName))
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(lastHashKey), block.Hash)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(block.Hash), block.Serialize())
			if err != nil {
				return err
			}
			lastHash = block.Hash
		}
		return nil
	})
	CheckError("CLI.NewChain #1",err)
}

func (cli *CLI) run() {
	cmdCreateChain := flag.NewFlagSet("newchain",flag.ExitOnError)
	cmdAdd := flag.NewFlagSet("addblock",flag.ExitOnError)
	cmdList := flag.NewFlagSet("listblock",flag.ExitOnError)
	if len(os.Args) < 2 {
		fmt.Printf("%s\n",cli.usage())
		os.Exit(-1)
	}
	switch os.Args[1] {
	case "newchain":
		err := cmdCreateChain.Parse(os.Args[2:])
		CheckError("CLI.run #1",err)
		if cmdCreateChain.Parsed(){
			cli.NewChain()
		}
	case "addblock":
		err := cmdAdd.Parse(os.Args[2:])
		CheckError("CLI.run #1",err)
		if cmdAdd.Parsed(){
				cli.AddBlock()
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
