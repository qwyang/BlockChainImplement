package main

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct{
	Version int64
	PrevBlockHash []byte
	Hash []byte
	TimeStamp int64
	TargetBits int64
	Nonce int64
	MerkleRoot []byte
	Data []byte
}
func (block *Block) SetHash(){
	data := bytes.Join(
		[][]byte{
			IntToBytes(block.Version),
			block.PrevBlockHash,
			IntToBytes(block.TimeStamp),
			IntToBytes(block.TargetBits),
			IntToBytes(block.Nonce),
			block.MerkleRoot,
			block.Data,
		},[]byte{})
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}

func NewBlock(data string,prevHash []byte) *Block {
	block := Block{
		Version:1,
		PrevBlockHash:prevHash,
		TimeStamp:time.Now().Unix(),
		TargetBits:10,
		Nonce:0,
		Data:[]byte(data),
	}
	block.SetHash()
	return &block
}
