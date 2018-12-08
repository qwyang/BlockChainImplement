package main

func main() {
	bc := NewBlockChain()
	cli := NewCLI(bc)
	cli.run()
}