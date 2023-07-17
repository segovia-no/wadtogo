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
	LumpEntries uint32
	LumpDirectoryOffset uint32
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


type Lump struct {
	LumpOffset uint32
	LumpSize uint32
	LumpName [8]byte
}

func (wp *WADParser) readLumpInfo(seekAt int64) Lump {
	wp.checkValidByteReader()

	wp.byteReader.Seek(seekAt, io.SeekStart)
	
	var lumpData Lump
	err := binary.Read(wp.byteReader, binary.LittleEndian, &lumpData)

	if err != nil {
		fmt.Println("[Error] readLumpInfo: Invalid data when reading WADs lump data:", err)
		os.Exit(1)
	}

	return lumpData
}

func (wp *WADParser) checkValidByteReader() {
	if (wp.byteReader == nil) {
		fmt.Println("[Error] No available byte reader, invoke it using setupByteReader(). Aborting execution")
		os.Exit(1)
	}
}
