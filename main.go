package main

import (
	"fmt"

	wl "github.com/segovia-no/wadtogo/wadloader"
)

var flagReader Flags
var wadLoader wl.WADLoader

func main() {
	fmt.Println("WADToGo - Another WAD Tool")
	fmt.Println("--------------------------")

	flagReader.parseFlags()
	wadLoader.OpenAndLoad(flagReader.WADFilename)

	fmt.Println("WAD Filename:", wadLoader.WADFilename)
	fmt.Println("WAD Type:", string(wadLoader.WADHeader.WadType[:]))
	fmt.Println("# Lumps:", wadLoader.WADHeader.LumpEntries)

	fmt.Println("--------------------------")

	wadLoader.ReadWADLumps()

	// Setup corresponding data depending on flags
	if flagReader.printWADMusicInfo || flagReader.dumpWADMusicInfo != "" || flagReader.exportMusic != "" {
		musicLumps, _ := wl.GetMusicLumps(wadLoader.WADLumps)
		wadLoader.Music = append(wadLoader.Music, musicLumps...)
	}

	if flagReader.printWADMapsInfo || flagReader.dumpWADMapsInfo != "" {
		wadLoader.LoadMaps()
	}

	if flagReader.exportSprites != "" {
		wadLoader.LoadPalettes()
		wadLoader.LoadGraphics()
	}

	// Command execution
	if flagReader.dumpLumpsInfo != "" {
		wl.DumpLumpsToTextFile(flagReader.dumpLumpsInfo, wadLoader.WADLumps)
	}

	if flagReader.printWADMusicInfo {
		wl.PrintSongNames(wadLoader.Music)
	}

	if flagReader.dumpWADMusicInfo != "" {
		wl.DumpSongNamesToTextFile(flagReader.dumpWADMusicInfo, wadLoader.Music)
	}

	if flagReader.printWADMapsInfo {
		wl.PrintMapNames(wadLoader.Maps)
	}

	if flagReader.dumpWADMapsInfo != "" {
		wl.DumpMapNamesToTextFile(flagReader.dumpWADMapsInfo, wadLoader.Maps)
	}

	if flagReader.exportSprites != "" {
		fmt.Println("Exporting sprites...")

		err := wadLoader.ExportAllSprites(flagReader.exportSprites)
		if err != nil {
			fmt.Println("[Error] Cannot export sprites - " + err.Error())
		}

		fmt.Println("Sprites exported successfully")
	}

	if flagReader.exportMusic != "" {
		fmt.Println("Exporting songs...")

		err := wadLoader.ExportAllSongs(flagReader.exportMusic)
		if err != nil {
			fmt.Println("[Error] Cannot export songs - " + err.Error())
		}
		fmt.Println("Songs exported successfully")
	}
}
