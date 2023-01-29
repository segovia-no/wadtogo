package main

type WADHeader struct {
	WadType [4]byte
	DirectoryEntries uint32
	DirectroryOffset uint32
}

type lumpInDirectory struct {
	lumpOffset uint32
	lumpSize uint32
	lumpName [8]byte
}
