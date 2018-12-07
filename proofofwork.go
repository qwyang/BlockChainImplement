package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 24

type ProofOfWork struct {
	block      *Block
	targetHash *big.Int
}

func NewPow(b *Block) *ProofOfWork {
	tmp := big.NewInt(1)
	hashValue := tmp.Lsh(tmp, 256-targetBits)
	pow := ProofOfWork{
		b,
		hashValue,
	}
	return &pow
}

func (pow *ProofOfWork) SetBlock(b *Block) {
	pow.block = b
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	block := pow.block
	data := bytes.Join(
		[][]byte{
			IntToBytes(block.Version),
			block.PrevBlockHash,
			IntToBytes(block.TimeStamp),
			IntToBytes(block.TargetBits),
			IntToBytes(nonce),
			block.MerkleRoot,
			block.Data,
		}, []byte{})
	return data
}

func (pow *ProofOfWork) mining() (int64, [32]byte) {
	var nonce int64
	var hash [32]byte
	var hashInt big.Int
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.targetHash) == -1 {
			fmt.Printf("find hash:%x,target hash:%x\n", hash, pow.targetHash)
			return nonce, hash
		}
		//fmt.Printf("curent hash:%x,target hash:%x\n",hash,pow.targetHash)
		nonce++
	}
	return nonce, hash
}
