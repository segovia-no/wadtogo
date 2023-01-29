package main

import (
	"fmt"
	"os"
)

var wp WADParser

type WADLoader struct {
	WADBuffer []byte
	WADFilename string
	WADHeader WADHeader
	WADDirectories WADDirectories
}

func (wl *WADLoader) openAndLoad(wadFilename string) {

	wl.WADFilename = wadFilename

	wadBuffer, readErr := os.ReadFile(wl.WADFilename)
	if readErr != nil {
		fmt.Println("[Error]: Couldn't read the WAD file", readErr)
		os.Exit(1)
	}

	wl.WADBuffer = wadBuffer

	wp.setupByteReader(wl.WADBuffer)
	wl.WADHeader = wp.readHeaderData()

	fmt.Println("WAD Type:", string(wl.WADHeader.WadType[:]))
	fmt.Println("WAD Dirs:", wl.WADHeader.DirectoryEntries)
	fmt.Println("WAD Dir offset:", wl.WADHeader.DirectoryOffset)

	wl.readWADDirectories()
	fmt.Println(wl.WADDirectories[0])

}

type WADDirectories []DirectoryData

func (wl *WADLoader) readWADDirectories() {

	if wl.WADBuffer == nil || wl.WADHeader.DirectoryEntries < 1 {
		fmt.Println("[Error]: Insufficient data to read WAD Directories")
		os.Exit(0)
	}

	dirOffset := int64(wl.WADHeader.DirectoryOffset)

	for i := 0; i < int(wl.WADHeader.DirectoryEntries); i++ {
		dirData := wp.readDirectoryData(dirOffset + int64(i*16))
		wl.WADDirectories = append(wl.WADDirectories, dirData)
	}
	
}
