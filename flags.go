package main

import (
	"fmt"
	"os"
)

type FlagReader struct {}

func checkValidArgsLength() {
	if len(os.Args) < 2 {
		fmt.Println("[Error] Not enough arguments. Aborting extraction")
		os.Exit(0)
	}
}

func (FlagReader) getWADFilenameFromFlag() string {
	checkValidArgsLength()

	filename := os.Args[1] // In the future, change this to use the flag module.

	if filename == "" {
		fmt.Println("[Error] getWADFilenameFromFlag: No WAD filename specified as argument. Aborting extraction")
		os.Exit(0)
	}

	return filename
}
