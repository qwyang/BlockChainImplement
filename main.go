package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("A transfer 5 coin to B")
	bc.AddBlock("A transfer 5 coin to C")
	for i, block := range bc.blocks {
		fmt.Printf("=================block:%d===================\n", i)
		fmt.Printf("Hash:0x%x\n", block.Hash)
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("version:%d\n", block.Version)
		fmt.Printf("PrevBlockHash:0x%x\n", block.PrevBlockHash)
		fmt.Printf("TimeStamp:%d\n", block.TimeStamp)
		fmt.Printf("TargetBits:0x%x\n", block.TargetBits)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Printf("MerkleRoot:0x%x\n", block.MerkleRoot)
	}
}
