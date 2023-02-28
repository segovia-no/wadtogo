package wadloader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)


type WADParser struct {
	byteReader *bytes.Reader
}

func (wp *WADParser) setupByteReader(b []byte) {
	wp.byteReader = bytes.NewReader(b)
}


type WADHeader struct {
	WadType [4]byte
	DirectoryEntries uint32
	DirectoryOffset uint32
}

func (wp *WADParser) readHeaderData() WADHeader {
	wp.checkValidByteReader()
	
	var wadHeader WADHeader
	err := binary.Read(wp.byteReader, binary.LittleEndian, &wadHeader)

	if err != nil {
		fmt.Println("[Error] readHeaderData: Invalid data when reading the WAD Header:", err)
		os.Exit(1)
	}

	return wadHeader
}


type DirectoryData struct {
	LumpOffset uint32
	LumpSize uint32
	LumpName [8]byte
}

func (wp *WADParser) readDirectoryData(seekAt int64) DirectoryData {
	wp.checkValidByteReader()

	wp.byteReader.Seek(seekAt, io.SeekStart)
	
	var directoryData DirectoryData
	err := binary.Read(wp.byteReader, binary.LittleEndian, &directoryData)

	if err != nil {
		fmt.Println("[Error] readDirectoryData: Invalid data when reading WADs directory data:", err)
		os.Exit(1)
	}

	return directoryData
}


func (wp *WADParser) checkValidByteReader() {
	if (wp.byteReader == nil) {
		fmt.Println("[Error] No available byte reader, invoke it using setupByteReader(). Aborting execution")
		os.Exit(1)
	}
}