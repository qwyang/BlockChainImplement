package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

type CLI struct {}

func NewCLI()*CLI{
	return &CLI{}
}

func (cli *CLI) usage()string{
	var buffer bytes.Buffer
	_,err := fmt.Fprintf(&buffer,"Usage:\n")
	CheckError("CLI.usage #0",err)
	_,err = fmt.Fprintf(&buffer,"%s newchain --address ADDRESS\n",os.Args[0])
	CheckError("CLI.usage #1",err)
	_,err = fmt.Fprintf(&buffer,"%s addblock\n",os.Args[0])
	CheckError("CLI.usage #2",err)
	_,err = fmt.Fprintf(&buffer,"%s listblock\n",os.Args[0])
	CheckError("CLI.usage #3",err)
	_,err = fmt.Fprintf(&buffer,"%s getbalance --Address ADDRESS\n",os.Args[0])
	CheckError("CLI.usage #4",err)
	_,err = fmt.Fprintf(&buffer,"%s newtx --from FromAddress --to ToAddress --amount Amount\n",os.Args[0])
	CheckError("CLI.usage #4",err)
	return string(buffer.Bytes())
}

func (cli *CLI) run() {
	cmdCreateChain := flag.NewFlagSet("newchain",flag.ExitOnError)
	cmdNewTx := flag.NewFlagSet("newtx",flag.ExitOnError)
	cmdMine := flag.NewFlagSet("mine",flag.ExitOnError)
	cmdList := flag.NewFlagSet("listblock",flag.ExitOnError)
	cmdGetUTXO := flag.NewFlagSet("getutxo",flag.ExitOnError)
	cmdBalance := flag.NewFlagSet("getbalance",flag.ExitOnError)
	cmdMineAddr := cmdMine.String("address","","Address Of Receiver")
	cmdBalanceAddr := cmdBalance.String("address","","Address Of Receiver")
	cmdCreateChainAddress := cmdCreateChain.String("address","","Address Of Receiver")
	cmdGetUTXOAddress := cmdGetUTXO.String("address","","Get UTXO by Address")
	cmdNewTxFrom := cmdNewTx.String("from","","From Address")
	cmdNewTxTo := cmdNewTx.String("to","","To Address")
	cmdNewTxAmount := cmdNewTx.Float64("amount",0.0,"Amount Money To Transfer")
	if len(os.Args) < 2 {
		fmt.Printf("%s\n",cli.usage())
		os.Exit(-1)
	}
	switch os.Args[1] {
	case "newchain":
		err := cmdCreateChain.Parse(os.Args[2:])
		CheckError("CLI.run #1",err)
		if cmdCreateChain.Parsed(){
			if *cmdCreateChainAddress != ""{
				cli.NewChain(*cmdCreateChainAddress)
			}else{
				fmt.Printf("%s\n",cli.usage())
				os.Exit(-1)
			}
		}
	case "newtx":
		err := cmdNewTx.Parse(os.Args[2:])
		CheckError("CLI.run #2",err)
		if cmdNewTx.Parsed(){
			cli.NewTransaction(*cmdNewTxFrom,*cmdNewTxTo,*cmdNewTxAmount)
		}
	case "mine":
		err := cmdMine.Parse(os.Args[2:])
		CheckError("CLI.run #3",err)
		if cmdMine.Parsed(){
			if *cmdMineAddr != ""{
				cli.Mine(*cmdMineAddr)
			}
		}
	case "listblock":
		err := cmdList.Parse(os.Args[2:])
		CheckError("CLI.run #4",err)
		if cmdList.Parsed(){
			cli.ListBlock()
		}
	case "getutxo":
		err := cmdGetUTXO.Parse(os.Args[2:])
		CheckError("CLI.run #5",err)
		if cmdGetUTXO.Parsed(){
			if *cmdGetUTXOAddress != ""{
				cli.GetUTXO(*cmdGetUTXOAddress)
			}
		}
	case "getbalance":
		err := cmdBalance.Parse(os.Args[2:])
		CheckError("CLI.run #6",err)
		if cmdBalance.Parsed(){
			if *cmdBalanceAddr != ""{
				cli.cmdBalance(*cmdBalanceAddr)
			}
		}
	default:
		fmt.Printf("%s\n",cli.usage())
		os.Exit(1)
	}
}
