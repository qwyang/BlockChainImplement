package main

import "fmt"

func main(){
	bc := NewBlockChain()
	bc.AddBlock("A transfer 5 coin to B")
	bc.AddBlock("A transfer 5 coin to C")
	for i,block := range bc.blocks {
		fmt.Printf("=================block:%d===================\n",i)
		fmt.Printf("Data:%s\n",block.Data)
		fmt.Printf("version:%d\n",block.Version)
		fmt.Printf("PrevBlockHash:%x\n",block.PrevBlockHash)
		fmt.Printf("TimeStamp:%d\n",block.TimeStamp)
		fmt.Printf("TargetBits:%x\n",block.TargetBits)
		fmt.Printf("Nonce:%d\n",block.Nonce)
		fmt.Printf("MerkleRoot:%x\n",block.MerkleRoot)
	}
}