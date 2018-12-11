package main

import (
	"bytes"
	"fmt"
)

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