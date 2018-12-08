package main

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	block := NewGenesisBlock()
	return &BlockChain{
		blocks: []*Block{block},
	}
}
func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.Hash
	block := NewBlock(data, prevHash)
	bc.blocks = append(bc.blocks, block)
}
