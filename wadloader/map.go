package wadloader

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Struct definitions for the WAD file format
type Map struct {
	Name string
	Vertexes []Vertex
	Linedefs []Linedef
	Things []Thing
	Nodes []Node
	Subsectors []SubSector
	Segs []Seg
}

type Vertex struct {
	XPos int16
	YPos int16
}

type Linedef struct {
	StartVertex uint16
	EndVertex uint16
	Flags uint16
	LineType uint16
	SectorTag uint16
	RightSidedef uint16
	LeftSidedef uint16
}

type Thing struct {
	XPos int16
	YPos int16
	Angle uint16
	Type uint16
	Flags uint16
}

type Node struct {
	XPartition int16
	YPartition int16
	ChangeXPartition int16
	ChangeYPartition int16
	RightBoxTop int16
	RightBoxBottom int16
	RightBoxLeft int16
	RightBoxRight int16
	LeftBoxTop int16
	LeftBoxBottom int16
	LeftBoxLeft int16
	LeftBoxRight int16
	RightChildIdx uint16
	LeftChildIdx uint16
}

type SubSector struct {
	SegCount uint16
	FirstSegIdx uint16
}

type Seg struct {
	StartVertex uint16
	EndVertex uint16
	Angle uint16
	LinedefIdx uint16
	Direction uint16
}

// TODO: make this function generic that is available for any struct that can be parsed from a WAD
func (wp *WADParser) parseMapThings(lump Lump) []Thing {
	wp.checkValidByteReader()

	var readThings []Thing
	lumpThingsCount := int(lump.LumpSize / 10)

	for i := 0; i < lumpThingsCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 10))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var t Thing
		err := binary.Read(wp.byteReader, binary.LittleEndian, &t)

		if err != nil {
			fmt.Println("[Warn] parseMapThings: Error while reading a thing (skipping): " + err.Error())
			continue
		}

		readThings = append(readThings, t)
	}

	return readThings
}
