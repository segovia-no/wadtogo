package main

import (
	"bytes"
	"fmt"
	"os"
)

var wp WADParser

type WADLoader struct {
	WADBuffer []byte
	WADFilename string
	WADHeader WADHeader
	WADDirectories WADDirectories

	Maps []MapDirectory
}

func (wl *WADLoader) openAndLoad(wadFilename string) {
	wl.WADFilename = wadFilename

	wadBuffer, readErr := os.ReadFile(wl.WADFilename)
	if readErr != nil {
		fmt.Println("[Error] openAndLoad: Couldn't read the WAD file", readErr)
		os.Exit(1)
	}

	wl.WADBuffer = wadBuffer

	wp.setupByteReader(wl.WADBuffer)
	wl.WADHeader = wp.readHeaderData()

	fmt.Println("WAD Filename:", wl.WADFilename)
	fmt.Println("WAD Type:", string(wl.WADHeader.WadType[:]))
	fmt.Println("Lumps:", wl.WADHeader.DirectoryEntries)

	wl.readWADDirectories()
	wl.detectMaps()
}

type WADDirectories []DirectoryData

func (wl *WADLoader) readWADDirectories() {
	if wl.WADBuffer == nil || wl.WADHeader.DirectoryEntries < 1 {
		fmt.Println("[Error] readWADDirectories: Insufficient data to read WAD Directories")
		os.Exit(0)
	}

	dirOffset := int64(wl.WADHeader.DirectoryOffset)

	for i := 0; i < int(wl.WADHeader.DirectoryEntries); i++ {
		dirData := wp.readDirectoryData(dirOffset + int64(i*16))
		wl.WADDirectories = append(wl.WADDirectories, dirData)
	}
}

type MapDirectory struct {
	MapName string
	Lumps []*DirectoryData
}

func (wl *WADLoader) detectMaps() {
	if len(wl.WADDirectories) < 1 {
		fmt.Println("[Warn] detectMaps: No Lumps detected loaded, cannot detect maps!")
		return
	}
	
	// map lumps use a 0 byte marker and complies with the minimum types of lumps
	for idx, lump := range wl.WADDirectories {
		if lump.LumpSize != 0 {
			continue
		}

		var currentMapDirectory MapDirectory

		neededMapLumps := []string{"THINGS", "LINEDEFS", "SIDEDEFS", "VERTEXES", "SEGS", "SSECTORS", "NODES", "SECTORS", "REJECT", "BLOCKMAP"}

		for _, nextLump := range wl.WADDirectories[idx + 1:] {
			if nextLump.LumpSize == 0 {
				break
			}

			for i := 0; i < len(neededMapLumps); i++ {
				nextLumpNameStr := string(bytes.Trim(nextLump.LumpName[:], "\x00"))

				if neededMapLumps[i] == string(nextLumpNameStr) {
					currentMapDirectory.Lumps = append(currentMapDirectory.Lumps, &nextLump)
					neededMapLumps = append(neededMapLumps[:i], neededMapLumps[i+1:]... )
				}

				if len(neededMapLumps) < 1 {
					currentMapDirectory.MapName = string(bytes.Trim(lump.LumpName[:], "\x00"))
					wl.Maps = append(wl.Maps, currentMapDirectory)
					break
				}
			}
		}
	}
}