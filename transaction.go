package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const reward = 50
const GenesisBlockInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"


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
	g_input := Input{nil,-1,data,}
	g_output := Output{reward,toAddress,	}
	tx := Transaction{
		nil,
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

func NewTransaction(fromAddress string,toAddress string, amount float64, bc *BlockChain) *Transaction{
	total,m := bc.GetSuitableUTXOs(fromAddress,amount)
	//fmt.Printf("NewTransaction SuitableUTXOS:%f,%v",total,m)
	if total < amount {
		return nil
	}
	inputs := []Input{}
	outputs := []Output{}
	for hash,indexes := range m {
		for _,index := range indexes{
			input := Input {[]byte(hash),index,fromAddress}
			inputs = append(inputs,input)
		}
	}
	output := Output{
		Value: amount,
		LockScript:toAddress,
	}
	outputs = append(outputs,output)
	if amount < total {//找零
		output = Output{
			Value: total - amount,
			LockScript:fromAddress,
		}
		outputs = append(outputs,output)
	}
	tx := Transaction{
		nil,
		 inputs,
		 outputs,
	}
	tx.SetId()
	return &tx
}

func (tx *Transaction) IsValid() bool {
	return true
}

func DeserializeTx(data []byte) *Transaction {
	var b Transaction
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&b)
	CheckError("Deserialize",err)
	return &b
}

