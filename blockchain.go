package main

import (
	"bytes"
	"errors"
	"github.com/boltdb/bolt"
)

const (
	DataBaseFile         = "blockchain.db"
	LastBlockHashKey     = "LastHashKey"
	BlockBucketName      = "blockchainBucket"
	NewTransactionBucket = "NewTransactionBucket"
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
		bucket := tx.Bucket([]byte(BlockBucketName))
		if bucket != nil{
			lastHash = bucket.Get([]byte(LastBlockHashKey))
		} else {
			transx := NewCoinbaseTx("",GenesisBlockInfo)
			block := NewGenesisBlock(transx)
			bucket, err = tx.CreateBucket([]byte(BlockBucketName))
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(LastBlockHashKey), block.Hash)
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
		bucket := tx.Bucket([]byte(BlockBucketName))
		if bucket == nil {
			err := errors.New("Empty Database" )
			CheckError("GetBlockChainHandler #2",err)
		}
		lastHash = bucket.Get([]byte(LastBlockHashKey))
		if lastHash == nil {
			err := errors.New("cannot find lasthashkey in db:" + string(LastBlockHashKey))
			CheckError("GetBlockChainHandler #2",err)
		}
		return nil
	})
	CheckError("GetBlockChainHandler #3",err)
	return &BlockChain{db,lastHash}
}

func (bc *BlockChain) AddBlock(txs []*Transaction) {
	block := NewBlock(txs,bc.lastHash)
	bc.SaveBlock(block)
	for _,tx := range block.Transactions {
		bc.RemoveNewTx(tx.ID)
	}
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
		bucket := tx.Bucket([]byte(BlockBucketName))
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
			if tx.IsCoinbase() { //CoinBase交易没有inputs,不统计
				continue
			}
			for index, input := range tx.Inputs {//遍历每个普通交易的输入
				//fmt.Printf("input:%v\n",input)
				if input.Unlock(address) {//属于本人的花费
					spent[string(input.TxID)] = append(spent[string(input.TxID)], int64(index))
					//fmt.Printf("spent:%v\n",spent)
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
	//fmt.Printf("%s,%v\n","GetSuitableUTXOs",utxos)
	for _,u := range utxos {
		for _,i := range u.indexes{
			money += u.tx.Outputs[i].Value
			m[string(u.tx.ID)] = append(m[string(u.tx.ID)],i)
			//fmt.Printf("%v\n",m)
			if money >= amount {
				goto EXIT
			}
		}
	}
EXIT:
	return money,m
}

func (bc *BlockChain) GetBalance(address string) float64 {
	var balance float64
	uts := bc.GetUTXOs(address)
	for _,u := range uts {
		for _,i := range u.indexes {
			balance += u.tx.Outputs[i].Value
		}
	}
	return balance
}

func (bc *BlockChain) SaveTx(transx *Transaction,bucketname string){
	err := bc.db.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		var err error
		bucket,err = tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			return err
		}
		err = bucket.Put(transx.ID,transx.Serialize())
		if err != nil {
			return err
		}
		return nil
	})
	CheckError("BlockChain.SaveTx",err)
}

func (bc *BlockChain) RemoveNewTx(hash []byte) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		var err error
		bucket,err = tx.CreateBucketIfNotExists([]byte(NewTransactionBucket))
		if err != nil {
			return err
		}
		err = bucket.Delete(hash)
		if err != nil {
			return err
		}
		return nil
	})
	CheckError("BlockChain.RemoveNewTx",err)
}

func (bc *BlockChain) GetNewTxs(num int) []*Transaction {
	var txs []*Transaction
	err := bc.db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			var i int = 0
			b := tx.Bucket([]byte(NewTransactionBucket))
			c := b.Cursor()
			for k, v := c.First(); k != nil && i < num; k, v = c.Next() {
				transx := DeserializeTx(v)
				txs = append(txs,transx)
				i++
			}
			return nil
		})
	CheckError("BlockChain.SaveTx",err)
	return txs
}

func (bc *BlockChain) SaveBlock(block *Block){
	err := bc.db.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		var err error
		bucket,err = tx.CreateBucketIfNotExists([]byte(BlockBucketName))
		if err != nil {
			return err
		}
		err = bucket.Put(block.Hash,block.Serialize())
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(LastBlockHashKey),block.Hash)
		if err != nil {
			return err
		}
		return nil
	})
	CheckError("BlockChain.SaveTx",err)
}