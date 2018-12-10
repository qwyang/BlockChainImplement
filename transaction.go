package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const reward = 50
const GenesisBlockInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Input struct {
	TxID []byte
	Index int64
	UnlockScripts string
}

func (in *Input) String() string {
	var buffer bytes.Buffer
	_,err := fmt.Fprintf(&buffer,"{PreTxHash:%x,Index:%d,Scripts:%s}",in.TxID,in.Index,in.UnlockScripts)
	CheckError("Input.String #1",err)
	return string(buffer.Bytes())
}

func (in *Input) Unlock(unlockdata string) bool {
	return in.UnlockScripts == unlockdata
}

type Output struct {
	Value float64
	LockScript string
}

func (out *Output) String() string {
	var buffer bytes.Buffer
	_,err := fmt.Fprintf(&buffer,"{Value:%x,Scripts:%s}",out.Value,out.LockScript)
	CheckError("Output.String #1",err)
	return string(buffer.Bytes())
}

func (out *Output) Unlock(unlockdata string) bool {
	return out.LockScript == unlockdata
}

type Transaction struct {
	ID []byte
	Inputs []Input
	Outputs []Output
}

func (tx *Transaction) String() string {
	var buffer bytes.Buffer
	_,err := fmt.Fprintf(&buffer,"{TxHash:%x,Inputs:%v,Outputs:%v}",tx.ID,tx.Inputs,tx.Outputs)
	CheckError("Transaction.String #1",err)
	return string(buffer.Bytes())
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
		LockScript:toAddress,
	}
	tx := Transaction{
		[]byte{},
		[]Input{g_input},
		[]Output{g_output},
	}
	tx.SetId()
	//fmt.Printf("transx:%v\n",tx)
	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].TxID == nil
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


