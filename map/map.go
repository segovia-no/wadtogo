package Map

type Map struct {
	Name string
	Vertexes []MapVertex
	Linedefs []Linedef
	Things []Thing
	Nodes []Node
	Subsectors []SubSector
	Segs []Seg
}

type MapVertex struct {
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
