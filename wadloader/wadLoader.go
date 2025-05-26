package wadloader

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type WADLoader struct {
	WADBuffer   []byte
	WADParser   WADParser
	WADFilename string
	WADHeader   WADHeader
	WADLumps    WADLumps

	Palettes []Palette
	Maps     []Map
	Music    []MusicLump
	Sprites  []Patch
	Patches  []Patch
	Flats    []Flat
}

func (wl *WADLoader) OpenAndLoad(wadFilename string) {
	wl.WADFilename = wadFilename

	wadBuffer, readErr := os.ReadFile(wl.WADFilename)
	if readErr != nil {
		fmt.Println("[Error] openAndLoad: Couldn't read the WAD file", readErr)
		os.Exit(1)
	}

	wl.WADBuffer = wadBuffer

	wl.WADParser.setupByteReader(wl.WADBuffer)
	wl.WADHeader = wl.WADParser.readHeaderData()
}

type WADLumps []Lump

func (wl *WADLoader) ReadWADLumps() {
	if wl.WADBuffer == nil || wl.WADHeader.LumpEntries < 1 {
		fmt.Println("[Error] readWADLumps: Insufficient data to read WAD Directories")
		os.Exit(0)
	}

	dirOffset := int64(wl.WADHeader.LumpDirectoryOffset)

	for i := 0; i < int(wl.WADHeader.LumpEntries); i++ {
		dirData := wl.WADParser.readLumpInfo(dirOffset + int64(i*16))
		wl.WADLumps = append(wl.WADLumps, dirData)
	}
}

func (wl *WADLoader) LoadMaps() {
	mapLumps := wl.DetectMaps()
	wl.LoadMapLumps(mapLumps)
}

type MapRawLumps struct {
	MapName string
	Lumps   []Lump
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

		for _, nextLump := range wl.WADLumps[idx+1:] {
			if nextLump.LumpSize == 0 {
				break
			}

			nextLumpNameStr := string(bytes.Trim(nextLump.LumpName[:], "\x00"))

			for i := 0; i < len(neededMapLumps); i++ {
				if neededMapLumps[i] == string(nextLumpNameStr) {
					currentMapLumps.Lumps = append(currentMapLumps.Lumps, nextLump)
					neededMapLumps = append(neededMapLumps[:i], neededMapLumps[i+1:]...)
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
				newMap.Things = wl.WADParser.parseMapThings(currLump)
			case "LINEDEFS":
				newMap.Linedefs = wl.WADParser.parseMapLinedefs(currLump)
			case "SIDEDEFS":
				newMap.Sidedefs = wl.WADParser.parseMapSidedefs(currLump)
			case "VERTEXES":
				newMap.Vertexes = wl.WADParser.parseMapVertexes(currLump)
			case "SEGS":
				newMap.Segs = wl.WADParser.parseMapSegs(currLump)
			case "SSECTORS":
				newMap.SSectors = wl.WADParser.parseMapSSectors(currLump)
			case "NODES":
				newMap.Nodes = wl.WADParser.parseMapNodes(currLump)
			case "SECTORS":
				newMap.Sectors = wl.WADParser.parseMapSectors(currLump)
			case "REJECT":
			case "BLOCKMAP":
				// TODO: Missing rest of implementations
			}
		}

		wl.Maps = append(wl.Maps, newMap)
	}
}

type MusicLump struct {
	name   string
	format string
	lump   Lump
}

func (wl *WADLoader) GetMusicLumps() ([]MusicLump, bool) {
	if len(wl.WADLumps) < 1 {
		fmt.Println("[Warn] getMusicLumps: No Lumps detected loaded, cannot detect music!")
		return nil, true
	}

	var musicLumps []MusicLump

	for _, lump := range wl.WADLumps {
		if lump.LumpSize == 0 {
			continue
		}

		// music lumps names start with "D_"
		lumpName := string(bytes.Trim(lump.LumpName[:], "\x00"))

		if !(strings.HasPrefix(lumpName, "D_")) {
			continue
		}

		musicFormat, err := wl.WADParser.getMusicFormatFromLump(&lump)
		if err != nil {
			errinfo := fmt.Sprintf("[Warn] getMusicLumps: Cannot detect music format for %v, omitting this lump., %v", lumpName, err)
			fmt.Println(errinfo)
		}

		curMusicLump := MusicLump{
			name:   lumpName,
			format: musicFormat,
			lump:   lump,
		}

		musicLumps = append(musicLumps, curMusicLump)
	}

	return musicLumps, false
}

func (wl *WADLoader) LoadPalettes() {
	palettes, err := wl.DetectPalettes()
	if err != nil {
		fmt.Println("[Error] LoadPalettes: Cannot detect palettes lumps, aborting.", err)
		os.Exit(1)
	}

	wl.Palettes = palettes
}

func (wl *WADLoader) LoadGraphics() {
	sprites, patches, flats, err := wl.DetectGraphics()
	if err != nil {
		fmt.Println("[Error] LoadGraphics: Cannot detect graphics lumps, aborting.", err)
		os.Exit(1)
	}

	wl.Sprites = sprites
	wl.Patches = patches
	wl.Flats = flats
}
