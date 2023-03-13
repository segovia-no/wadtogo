package wadloader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var MapLumpsNames = []string{"THINGS", "LINEDEFS", "SIDEDEFS", "VERTEXES", "SEGS", "SSECTORS", "NODES", "SECTORS", "REJECT", "BLOCKMAP"}

// Struct definitions for Map data
type Map struct {
	Name string
	Things []Thing
	Linedefs []Linedef
	Sidedefs []Sidedef
	Vertexes []Vertex
	Segs []Seg
	SSectors []SSector
	Nodes []Node
	Sectors []Sector

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

type Sidedef struct {
	XOffset int16
	YOffset int16
	UpperTexture [8]byte
	LowerTexture [8]byte
	MiddleTexture [8]byte
	SectorIdx uint16
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

type SSector struct {
	SegCount uint16
	FirstSegIdx uint16
}

type Seg struct {
	VertexStart uint16
	VertexEnd uint16
	Angle int16
	LinedefNumber uint16
	Direction int16
	Offset int16
}

type Sector struct {
	FloorHeight int16
	CeilingHeight int16
	FloorTexture [8]byte
	CeilingTexture [8]byte
	LightLevel uint16
	Type uint16
	Tag uint16
}

type Reject []byte

type Blockmap struct { //TODO: Broken struct
	XOrigin int16
	YOrigin int16
	ColumnCount int16
	RowCount int16
	Offsets []int16
	Blocklists []int16
}

// Map lump parsing
func (wp *WADParser) parseMapThings(lump Lump) []Thing {  // TODO: can we do this generic for all map lumps???
	wp.checkValidByteReader()

	var readThings []Thing
	lumpThingsCount := int(lump.LumpSize / 10)

	for i := 0; i < lumpThingsCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 10))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var t Thing
		err := binary.Read(wp.byteReader, binary.LittleEndian, &t)

		if err != nil {
			fmt.Println("[Warn] parseMapThings: Error while reading a THING (skipping): " + err.Error())
			continue
		}

		readThings = append(readThings, t)
	}

	return readThings
}

func (wp *WADParser) parseMapLinedefs(lump Lump) []Linedef {
	wp.checkValidByteReader()

	var readLinedef []Linedef
	lumpThingsCount := int(lump.LumpSize / 14)

	for i := 0; i < lumpThingsCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 14))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var l Linedef
		err := binary.Read(wp.byteReader, binary.LittleEndian, &l)

		if err != nil {
			fmt.Println("[Warn] parseMapLinedefs: Error while reading a LINEDEF (skipping): " + err.Error())
			continue
		}

		readLinedef = append(readLinedef, l)
	}

	return readLinedef
}

func (wp *WADParser) parseMapSidedefs(lump Lump) []Sidedef {
	wp.checkValidByteReader()

	var readSidedef []Sidedef
	lumpSidedefCount := int(lump.LumpSize / 30)

	for i := 0; i < lumpSidedefCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 30))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var s Sidedef
		err := binary.Read(wp.byteReader, binary.LittleEndian, &s)

		if err != nil {
			fmt.Println("[Warn] parseMapSidedefs: Error while reading a SIDEDEF (skipping): " + err.Error())
			continue
		}

		readSidedef = append(readSidedef, s)
	}

	return readSidedef
}

func (wp *WADParser) parseMapVertexes(lump Lump) []Vertex {
	wp.checkValidByteReader()

	var readVertex []Vertex
	lumpVertexCount := int(lump.LumpSize / 4)

	for i := 0; i < lumpVertexCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 4))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var v Vertex
		err := binary.Read(wp.byteReader, binary.LittleEndian, &v)

		if err != nil {
			fmt.Println("[Warn] parseMapVertexes: Error while reading a VERTEX (skipping): " + err.Error())
			continue
		}

		readVertex = append(readVertex, v)
	}

	return readVertex
}

func (wp *WADParser) parseMapSegs(lump Lump) []Seg {
	wp.checkValidByteReader()

	var readSeg []Seg
	lumpSegCount := int(lump.LumpSize / 12)

	for i := 0; i < lumpSegCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 12))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var s Seg
		err := binary.Read(wp.byteReader, binary.LittleEndian, &s)

		if err != nil {
			fmt.Println("[Warn] parseMapSegs: Error while reading a SEG (skipping): " + err.Error())
			continue
		}

		readSeg = append(readSeg, s)
	}

	return readSeg
}

func (wp *WADParser) parseMapSSectors(lump Lump) []SSector {
	wp.checkValidByteReader()

	var readSSector []SSector
	lumpSSectorCount := int(lump.LumpSize / 4)

	for i := 0; i < lumpSSectorCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 4))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var ss SSector
		err := binary.Read(wp.byteReader, binary.LittleEndian, &ss)

		if err != nil {
			fmt.Println("[Warn] parseMapSSectors: Error while reading a SSECTOR (skipping): " + err.Error())
			continue
		}

		readSSector = append(readSSector, ss)
	}

	return readSSector
}

func (wp *WADParser) parseMapNodes(lump Lump) []Node {
	wp.checkValidByteReader()

	var readNode []Node
	lumpNodeCount := int(lump.LumpSize / 28)

	for i := 0; i < lumpNodeCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 28))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var n Node
		err := binary.Read(wp.byteReader, binary.LittleEndian, &n)

		if err != nil {
			fmt.Println("[Warn] parseMapNodes: Error while reading a NODE (skipping): " + err.Error())
			continue
		}

		readNode = append(readNode, n)
	}

	return readNode
}

func (wp *WADParser) parseMapSectors(lump Lump) []Sector {
	wp.checkValidByteReader()

	var readSector []Sector
	lumpSectorCount := int(lump.LumpSize / 26)

	for i := 0; i < lumpSectorCount; i++ {
		byteOffset := int64(lump.LumpOffset) + int64((i * 26))
		wp.byteReader.Seek(int64(byteOffset), io.SeekStart)

		var s Sector
		err := binary.Read(wp.byteReader, binary.LittleEndian, &s)

		if err != nil {
			fmt.Println("[Warn] parseMapSectors: Error while reading a SECTOR (skipping): " + err.Error())
			continue
		}

		readSector = append(readSector, s)
	}

	return readSector
}

// helper functions
func IsLumpAMapLump(l *Lump) bool {
	for _, v := range MapLumpsNames {
		lumpName := string(bytes.Trim(l.LumpName[:], "\x00"))
		if lumpName == v {
			return true
		}
	}
	return false
}
