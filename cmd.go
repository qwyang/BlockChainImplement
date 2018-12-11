package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

func PrintBlock(block *Block){
	fmt.Printf("Hash:0x%x\n", block.Hash)
	fmt.Printf("Transactions: %s\n", block.Transactions)
	fmt.Printf("version:%d\n", block.Version)
	fmt.Printf("PrevBlockHash:0x%x\n", block.PrevBlockHash)
	fmt.Printf("TimeStamp:%d\n", block.TimeStamp)
	fmt.Printf("TargetBits:0x%x\n", block.TargetBits)
	fmt.Printf("Nonce:%d\n", block.Nonce)
	fmt.Printf("MerkleRoot:0x%x\n", block.MerkleRoot)
	pow := NewPow(block)
	fmt.Printf("isValid:%v\n", pow.IsValid())
}

func (cli *CLI) NewTransaction(from,to string,amount float64){
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	tx := NewTransaction(from,to,amount,bc)
	if tx != nil {
		bc.SaveTx(tx,NewTransactionBucket)
	}
}

func (cli *CLI) Mine(address string){
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	//get transactions
	var totalTxs []*Transaction = []*Transaction{}
	tx := NewCoinbaseTx(address,"CoinBase")
	//bc.SaveTx(tx,NewTransactionBucket)
	totalTxs = append(totalTxs,tx)
	txs := bc.GetNewTxs(10)
	totalTxs = append(totalTxs,txs...)
	bc.AddBlock(totalTxs)
}

func (cli *CLI) ListBlock() {
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	iter := bc.Iterator()
	var i int = 1
	for b := iter.Next();b != nil;b = iter.Next() {
		fmt.Printf("=================block:%d===================\n", i)
		PrintBlock(b)
		i++
	}
}

func checkDatabaseExist()bool{
	_,err := os.Stat(DataBaseFile)
	if os.IsNotExist(err){
		return false
	}
	return true
}


func (cli *CLI) NewChain(address string){
	if checkDatabaseExist() {
		fmt.Printf("BlockChain Already Exist.\n")
		return
	}
	transx := NewCoinbaseTx(address,GenesisBlockInfo)
	block := NewGenesisBlock(transx)
	db,err := bolt.Open(DataBaseFile,0600,nil)
	CheckError("CLI.NewChain #1",err)
	defer db.Close()
	bc := &BlockChain{db,block.Hash}
	//bc.SaveTx(transx,NewTransactionBucket)
	err = db.Update(func(tx *bolt.Tx)error{
		_,err := tx.CreateBucketIfNotExists([]byte(NewTransactionBucket))
		if err != nil {
			return err
		}
		return nil
	})
	CheckError("CLI.NewChain #1",err)
	bc.SaveBlock(block)
}

func (cli *CLI) GetUTXO (address string){
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	fmt.Printf("------------------%s utxo-----------------------\n",address)
	for _,u := range bc.GetUTXOs(address){
		fmt.Printf("transaction:%s,unused outputs:%v\n",u.tx,u.indexes)
	}
}

func (cli *CLI) cmdBalance(address string){
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	balance := bc.GetBalance(address)
	fmt.Printf("the balance of %s is：%0.4f\n",address,balance)
}