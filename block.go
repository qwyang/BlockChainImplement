package main

import (
	"time"
)

type Block struct {
	Version       int64
	PrevBlockHash []byte
	Hash          []byte
	TimeStamp     int64
	TargetBits    int64
	Nonce         int64
	MerkleRoot    []byte
	Data          []byte
}

func NewGenesisBlock() *Block {
	block := NewBlock("Genesis Block", []byte{})
	return block
}

func NewBlock(data string, prevHash []byte) *Block {
	block := Block{
		Version:       1,
		PrevBlockHash: prevHash,
		TimeStamp:     time.Now().Unix(),
		TargetBits:    targetBits,
		Nonce:         0,
		Data:          []byte(data),
	}
	pow := NewPow(&block)
	nonce,hash := pow.mining()
	block.Nonce = nonce
	block.Hash = hash[:]
	return &block
}
