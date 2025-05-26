package main

import (
	"fmt"
	"os"

	wl "github.com/segovia-no/wadtogo/wadloader"
)

var flagReader Flags
var wads []wl.WADLoader

func main() {
	fmt.Println("WADToGo - Another WAD Tool")
	fmt.Println("--------------------------")

	flagReader.parseFlags()

	var wadcount uint8 = uint8(len(flagReader.WADFilenames))
	if wadcount < 1 {
		fmt.Println("Cannot perform actions without at least one WAD file")
	}

	wads = make([]wl.WADLoader, wadcount)
	for idx, filepath := range flagReader.WADFilenames {
		wads[idx] = loadWAD(filepath)
	}

	if flagReader.mergeWADS {
		processMergeWads()
		return
	}

	for _, wad := range wads {
		processSingleFileActions(&wad)
	}

}

func loadWAD(filePath string) wl.WADLoader {
	wad := wl.WADLoader{}
	wad.OpenAndLoad(filePath)

	fmt.Println("WAD Filename:", wad.WADFilename)
	fmt.Println("WAD Type:", string(wad.WADHeader.WadType[:]))
	fmt.Println("# Lumps:", wad.WADHeader.LumpEntries)

	fmt.Println("--------------------------")
	wad.ReadWADLumps()

	return wad
}

func processMergeWads() {
	if len(wads) < 2 {
		fmt.Println("Cannot merge wads with less than two WAD files")
		os.Exit(1)
	}

	//TODO: Merge wads
}

func processSingleFileActions(wad *wl.WADLoader) {
	if flagReader.printWADMusicInfo || flagReader.dumpWADMusicInfo != "" || flagReader.exportMusic != "" {
		musicLumps, _ := wad.GetMusicLumps()
		wad.Music = append(wad.Music, musicLumps...)
	}

	if flagReader.printWADMapsInfo || flagReader.dumpWADMapsInfo != "" {
		wad.LoadMaps()
	}

	if flagReader.exportSprites != "" {
		wad.LoadPalettes()
		wad.LoadGraphics()
	}

	// Command execution
	if flagReader.dumpLumpsInfo != "" {
		wl.DumpLumpsToTextFile(flagReader.dumpLumpsInfo, wad.WADLumps)
	}

	if flagReader.printWADMusicInfo {
		wl.PrintSongNames(wad.Music)
	}

	if flagReader.dumpWADMusicInfo != "" {
		wl.DumpSongNamesToTextFile(flagReader.dumpWADMusicInfo, wad.Music)
	}

	if flagReader.printWADMapsInfo {
		wl.PrintMapNames(wad.Maps)
	}

	if flagReader.dumpWADMapsInfo != "" {
		wl.DumpMapNamesToTextFile(flagReader.dumpWADMapsInfo, wad.Maps)
	}

	if flagReader.exportSprites != "" {
		fmt.Println("Exporting sprites...")

		err := wad.ExportAllSprites(flagReader.exportSprites)
		if err != nil {
			fmt.Println("[Error] Cannot export sprites - " + err.Error())
		}

		fmt.Println("Sprites exported successfully")
	}

	if flagReader.exportMusic != "" {
		fmt.Println("Exporting songs...")

		err := wad.ExportAllSongs(flagReader.exportMusic)
		if err != nil {
			fmt.Println("[Error] Cannot export songs - " + err.Error())
		}
		fmt.Println("Songs exported successfully")
	}
}
