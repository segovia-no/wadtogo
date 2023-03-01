package wadloader

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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


func (wp *WADParser) getMusicFormatFromLump(lump *Lump) (string, error) {
	wp.checkValidByteReader()

	wp.byteReader.Seek(int64(lump.LumpOffset), io.SeekStart)

	// Look for identification header
	// ASCII = MTHD -> MIDI format
	// ASCII = MUS -> MUS format

	var header [4]byte
	err := binary.Read(wp.byteReader, binary.LittleEndian, &header)

	if err != nil {
		return "Unknown", errors.New("[Error] getMusicFormatFromLump: Couldn't read music lump header")
	}

	musicFormat := string(header[:])

	if strings.Contains(musicFormat, "MTHD") {
		return "MIDI", nil
	} else if strings.Contains(musicFormat, "MUS") {
		return "MUS", nil
	}

	return "Invalid", errors.New("[Error] getMusicFormatFromLump: Invalid music format detected")
}


func (wp *WADParser) checkValidByteReader() {
	if (wp.byteReader == nil) {
		fmt.Println("[Error] No available byte reader, invoke it using setupByteReader(). Aborting execution")
		os.Exit(1)
	}
}