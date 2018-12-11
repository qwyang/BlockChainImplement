package main

import (
	"bytes"
	"fmt"
)

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
