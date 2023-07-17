package wadloader

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

var wp WADParser

type WADLoader struct {
	WADBuffer []byte
	WADFilename string
	WADHeader WADHeader
	WADLumps WADLumps

	Palette []Palette
	Maps []Map
	Music []MusicLump
	Graphics []Patch
}

func (wl *WADLoader) OpenAndLoad(wadFilename string) {
	wl.WADFilename = wadFilename

	wadBuffer, readErr := os.ReadFile(wl.WADFilename)
	if readErr != nil {
		fmt.Println("[Error] openAndLoad: Couldn't read the WAD file", readErr)
		os.Exit(1)
	}

	wl.WADBuffer = wadBuffer

	wp.setupByteReader(wl.WADBuffer)
	wl.WADHeader = wp.readHeaderData()
}

type WADLumps []Lump

func (wl *WADLoader) ReadWADLumps() {
	if wl.WADBuffer == nil || wl.WADHeader.LumpEntries < 1 {
		fmt.Println("[Error] readWADLumps: Insufficient data to read WAD Directories")
		os.Exit(0)
	}

	dirOffset := int64(wl.WADHeader.LumpDirectoryOffset)

	for i := 0; i < int(wl.WADHeader.LumpEntries); i++ {
		dirData := wp.readLumpInfo(dirOffset + int64(i*16))
		wl.WADLumps = append(wl.WADLumps, dirData)
	}
}

func (wl *WADLoader) LoadMaps() {
	mapLumps := wl.DetectMaps()
	wl.LoadMapLumps(mapLumps)
}

type MapRawLumps struct {
	MapName string
	Lumps []Lump
}

func (wl *WADLoader) DetectMaps() []MapRawLumps {

	var rawMaps []MapRawLumps

	if len(wl.WADLumps) < 1 {
		fmt.Println("[Warn] DetectMaps: No Lumps detected loaded, cannot detect maps!")
		return rawMaps
	}
	
	// map lumps use a 0 byte marker and complies with the minimum types of lumps
	for idx, lump := range wl.WADLumps {
		if lump.LumpSize != 0 || !strings.HasPrefix(string(lump.LumpName[:]), "E") {
			continue
		}

		var currentMapLumps MapRawLumps

		var neededMapLumps []string
		neededMapLumps = append(neededMapLumps, MapLumpsNames...)

		for _, nextLump := range wl.WADLumps[idx + 1:] {
			if nextLump.LumpSize == 0 {
				break
			}

			nextLumpNameStr := string(bytes.Trim(nextLump.LumpName[:], "\x00"))

			for i := 0; i < len(neededMapLumps); i++ {
				if neededMapLumps[i] == string(nextLumpNameStr) {
					currentMapLumps.Lumps = append(currentMapLumps.Lumps, nextLump)
					neededMapLumps = append(neededMapLumps[:i], neededMapLumps[i+1:]... )
				}
			}

			if len(neededMapLumps) < 1 {
				currentMapLumps.MapName = string(bytes.Trim(lump.LumpName[:], "\x00"))
				rawMaps = append(rawMaps, currentMapLumps)
				break
			}
		}
	}

	return rawMaps
}

func (wl *WADLoader) LoadMapLumps(allMapsRaw []MapRawLumps) {
	if len(allMapsRaw) < 1 {
		fmt.Println("[Warn] LoadMapLumps: No maps inside the provided slice")
		return
	}

	for _, currMap := range allMapsRaw {
		var newMap Map
		newMap.Name = currMap.MapName

		for _, currLump := range currMap.Lumps {
			lumpNameStr := string(bytes.Trim(currLump.LumpName[:], "\x00"))

			switch lumpNameStr {
			case "THINGS":
				newMap.Things = wp.parseMapThings(currLump)
			case "LINEDEFS":
				newMap.Linedefs = wp.parseMapLinedefs(currLump)
			case "SIDEDEFS":
				newMap.Sidedefs = wp.parseMapSidedefs(currLump)
			case "VERTEXES":
				newMap.Vertexes = wp.parseMapVertexes(currLump)
			case "SEGS":
				newMap.Segs = wp.parseMapSegs(currLump)
			case "SSECTORS":
				newMap.SSectors = wp.parseMapSSectors(currLump)
			case "NODES":
				newMap.Nodes = wp.parseMapNodes(currLump)
			case "SECTORS":
				newMap.Sectors = wp.parseMapSectors(currLump)
			case "REJECT":
			case "BLOCKMAP":
				// TODO: Missing rest of implementations
			}
		}

		wl.Maps = append(wl.Maps, newMap)
	}
}


type MusicLump struct {
	name string
	format string
	lump Lump
}

func GetMusicLumps(wl WADLumps) ([]MusicLump, bool) {
	if len(wl) < 1 {
		fmt.Println("[Warn] getMusicLumps: No Lumps detected loaded, cannot detect music!")
		return nil, true
	}

	var musicLumps []MusicLump

	for _, lump := range wl {
		if lump.LumpSize == 0 {
			continue
		}
		
		// music lumps names start with "D_"
		lumpName := string(bytes.Trim(lump.LumpName[:], "\x00"))

		if !(strings.HasPrefix(lumpName, "D_")) {
			continue
		}

		musicFormat, err := wp.getMusicFormatFromLump(&lump)
		if err != nil {
			errinfo := fmt.Sprintf("[Warn] getMusicLumps: Cannot detect music format for %v, omitting this lump., %v", lumpName, err)
			fmt.Println(errinfo)
		}

		curMusicLump := MusicLump {
			name: lumpName,
			format: musicFormat,
			lump: lump,
		}

		musicLumps = append(musicLumps, curMusicLump)
	}

	return musicLumps, false
}
