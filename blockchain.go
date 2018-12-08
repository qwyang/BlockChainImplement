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

func NewBlockChain() *BlockChain {
	var lastHash []byte
	db,err := bolt.Open(DataBaseFile,0600,nil)
	CheckError("NewBlockChain #1",err)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil{
			lastHash = bucket.Get([]byte(lastHashKey))
		} else {
			block := NewGenesisBlock()
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
	CheckError("NewBlockChain #1",err)
	return &BlockChain{
		db,
		lastHash,
	}
}
func (bc *BlockChain) AddBlock(data string) {
	block := NewBlock(data,bc.lastHash)
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