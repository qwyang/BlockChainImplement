package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

const reward = 50
const GenesisBlockInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Input struct {
	TxID []byte
	Index int64
	UnlockScripts string
}

type Output struct {
	Value float64
	LockScript string
}

type Transaction struct {
	ID []byte
	Inputs []Input
	Outputs []Output
}

func (tx *Transaction) Serialize() []byte {
	buffer := bytes.Buffer{}
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(tx)
	CheckError("Transaction.SetId",err)
	return buffer.Bytes()
}

func (tx *Transaction) SetId(){
	data := sha256.Sum256(tx.Serialize())
	tx.ID = data[:]
}

func NewCoinbaseTx(toAddress string,data string) *Transaction {
	g_input := Input{
		UnlockScripts:data,
	}
	g_output := Output{
		Value:reward,
	}
	tx := Transaction{
		[]byte{},
		[]Input{g_input},
		[]Output{g_output},
	}
	tx.SetId()
	return &tx
}

func NewTransaction(fromAddress string,toAddress string, value float64) *Transaction{
	g_input := Input{

	}
	g_output := Output{
		Value:value,
	}
	tx := Transaction{
		[]byte{},
		[]Input{g_input},
		[]Output{g_output},
	}
	tx.SetId()
	return &tx
}


