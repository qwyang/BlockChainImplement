package main

func main() {
	/*TestCase:
	newchain(Alice)
	mine(Alice)
	mine(Bob)
	newTx(from:Alice,to:Bob,10)
	newTx(from:Alice,to:Bob,10)
	mine(Alice)
	getbalance(Alice)=> 80
	getbalance(Bob)=> 70
	 */
	cli := NewCLI()
	cli.run()
}