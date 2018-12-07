package main

type BlockChain struct {
	blocks []*Block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block",[]byte{})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		blocks: []*Block{NewGenesisBlock()},
	}
}
func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.blocks[len(bc.blocks) - 1]
	prevHash := lastBlock.Hash
	block := NewBlock(data,prevHash)
	bc.blocks = append(bc.blocks,block)
}