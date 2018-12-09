package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

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

func NewTransaction(from string,to string, value uint64) *Transaction{
	g_input := Input{

	}
	g_output := Output{

	}
	tx := Transaction{
		[]byte{},
		[]Input{g_input},
		[]Output{g_output},
	}
	tx.SetId()
	return &tx
}


