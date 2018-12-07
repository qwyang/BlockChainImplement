package main

type BlockChain struct {
	blocks []*Block
	pow    *ProofOfWork
}

func NewGenesisBlock() *Block {
	block := NewBlock("Genesis Block", []byte{})
	return block
}

func NewBlockChain() *BlockChain {
	block := NewGenesisBlock()
	block.TargetBits = targetBits
	pow := NewPow(block)
	nonce, hash := pow.mining()
	block.Nonce = nonce
	block.Hash = hash[:]
	return &BlockChain{
		blocks: []*Block{NewGenesisBlock()},
		pow:    pow,
	}
}
func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.Hash
	block := NewBlock(data, prevHash)
	block.TargetBits = targetBits
	bc.pow.SetBlock(block)
	nonce, hash := bc.pow.mining()
	block.Nonce = nonce
	block.Hash = hash[:]
	bc.blocks = append(bc.blocks, block)
}
