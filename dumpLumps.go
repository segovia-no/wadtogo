package main

import (
	"bytes"
	"fmt"
	"os"
)

func dumpLumpsToTextFile(filename string, lumps WADDirectories) {
	os.Remove(filename)
	file, err := os.Create(filename)

	if err != nil {
		fmt.Println("[Error] Cannot create a file called", filename, err)
		return
	}

	defer file.Close()

	var errWrite error

	_, errWrite = file.WriteString("Lump name | Size\n")

	for _, lump := range lumps {

		lumpName := bytes.Trim(lump.LumpName[:], "\x00")

		outStr := fmt.Sprintf("%s %v \n", lumpName, lump.LumpSize)
		_, errWrite = file.WriteString(outStr)
	}

	if errWrite != nil {
		fmt.Println("[Error] Cannot add lump data to dump file", filename, err)
		return
	}

	fmt.Println("[Info] Lumps dumped into", filename)
}