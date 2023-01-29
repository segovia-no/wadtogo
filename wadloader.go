package main

import (
	"fmt"
	"os"
)

var wadParser WADParser

type WADLoader struct {}

func (WADLoader) openAndLoad(wadFilename string) {

	wadData, readErr := os.ReadFile(wadFilename)
	if readErr != nil {
		fmt.Println("[Error]: Couldn't read the WAD file")
		os.Exit(1)
	}

	wadHeader := wadParser.readHeaderData(wadData)
	wadType := string(wadHeader.WadType[:])

	fmt.Println(wadType)

}
