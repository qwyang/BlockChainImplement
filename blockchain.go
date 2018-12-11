package main

import (
	"bytes"
	"errors"
	"github.com/boltdb/bolt"
)

package main

import (
"bytes"
"errors"
"github.com/boltdb/bolt"
)

const (
	DataBaseFile = "blockchain.db"
	lastHashKey = "LastHashKey"
	bucketName = "blockchainBucket"
)

type BlockChain struct {
	db *bolt.DB
	lastHash []byte
}

type BlockChainIterator struct{
	db *bolt.DB
	currentHash []byte
}

type UTXO struct {
	tx *Transaction
	indexes []int64
}

func NewBlockChain() *BlockChain {
	var lastHash []byte
	db,err := bolt.Open(DataBaseFile,0600,nil)
	CheckError("NewBlockChain #1",err)
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
	CheckError("NewBlockChain #3",err)
	return &BlockChain{db,lastHash}
}

func GetBlockChainHandler() *BlockChain {
	var lastHash []byte
	db,err := bolt.Open(DataBaseFile,0600,nil)
	CheckError("GetBlockChainHandler #1",err)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			err := errors.New("Empty Database" )
			CheckError("GetBlockChainHandler #2",err)
		}
		lastHash = bucket.Get([]byte(lastHashKey))
		if lastHash == nil {
			err := errors.New("cannot find lasthashkey in db:" + string(lastHashKey))
			CheckError("GetBlockChainHandler #2",err)
		}
		return nil
	})
	CheckError("GetBlockChainHandler #3",err)
	return &BlockChain{db,lastHash}
}

func (bc *BlockChain) AddBlock(txs []*Transaction) {
	block := NewBlock(txs,bc.lastHash)
	err := bc.db.Update(func(tx *bolt.Tx)error{
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil {
			err := bucket.Put(block.Hash,block.Serialize())
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(lastHashKey),block.Hash)
			if err != nil {
				return err
			}
			bc.lastHash = block.Hash
		}else {
			err := errors.New("Fatal:AddBlock can not find the \"%s\" bucket in DB")
			CheckError("AddBlock #1",err)
		}
		return nil
	})
	CheckError("AddBlock #2",err)
}

func (bc *BlockChain) Iterator() *BlockChainIterator{
	return &BlockChainIterator{
		bc.db,bc.lastHash,
	}
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	if bytes.Equal(iter.currentHash,[]byte{}) {
		return nil
	}
	err:=iter.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		data := bucket.Get([]byte(iter.currentHash))
		if data == nil {
			err := errors.New("cannot find hashkey in db:" + string(iter.currentHash))
			CheckError("BlockChainIterator.Next #1",err)
		}
		block = Deserialize(data)
		iter.currentHash = block.PrevBlockHash
		return nil
	})
	CheckError("BlockChainIterator.Next #2",err)

	return block
}

func IsIn(set []int64,v int64) bool {
	if set == nil {
		return false
	}
	for _,x := range set {
		if v == x {
			return  true
		}
	}
	return false
}

func (bc *BlockChain)GetUTXOs(address string) []UTXO{
	utxo := []UTXO{}
	spent := make(map[string][]int64)
	iter := bc.Iterator()
	for block := iter.Next();block != nil;block = iter.Next() {
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() { //CoinBase交易没有inputs,不统计
				continue
			}
			for index, input := range tx.Inputs {//遍历每个普通交易的输入
				if input.Unlock(address) {//属于本人的花费
					spent[string(input.TxID)] = append(spent[string(input.TxID)], int64(index))
				}
			}
		}
		//遍历outputs
		for _, tx := range block.Transactions {
			u := UTXO{tx,[]int64{}}
			for index, output := range tx.Outputs {//遍历每个普通交易的输入
				if output.Unlock(address) && !IsIn(spent[string(tx.ID)],int64(index)){//属于我的且没有花费的output
					u.indexes = append(u.indexes,int64(index))
					utxo = append(utxo,u)
				}
			}
		}
	}
	return utxo
}

func (bc *BlockChain)GetSuitableUTXOs(address string, amount float64) (float64,map[string][]int64) {
	m := make(map[string][]int64)
	var money float64 = 0.0
	utxos := bc.GetUTXOs(address)
	for _,u := range utxos {
		for _,i := range u.indexes{
			money += u.tx.Outputs[i].Value
			if money < amount {
				m[string(u.tx.ID)] = append(m[string(u.tx.ID)],i)
			}else {
				goto EXIT
			}
		}
	}
EXIT:
	return money,m
}