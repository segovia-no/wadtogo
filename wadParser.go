package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)


type WADParser struct {}

func (WADParser) readHeaderData(b []byte) WADHeader {

	byteReader := bytes.NewReader(b)
	
	var wadHeader WADHeader
	err := binary.Read(byteReader, binary.LittleEndian, &wadHeader)

	if err != nil {
		fmt.Println("[Error]: Invalid data when reading the WAD Header ", err)
		os.Exit(1)
	}

	return wadHeader

}
