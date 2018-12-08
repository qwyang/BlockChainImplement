package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func IntToBytes(d int64)[]byte{
	buffer := bytes.Buffer{}
	err := binary.Write(&buffer,binary.BigEndian,d)
	if err!=nil {
		fmt.Printf("error:%s",err)
		os.Exit(-1)
	}
	return buffer.Bytes()
}

func CheckError(data string,err error){
	if err!=nil {
		fmt.Printf("error:%v,%v",data,err)
		os.Exit(-1)
	}
}
