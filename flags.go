package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	WADFilenames      []string
	printWADMusicInfo bool
	printWADMapsInfo  bool
	dumpLumpsInfo     string
	dumpWADMusicInfo  string
	dumpWADMapsInfo   string
	exportMusic       string
	exportSprites     string
	mergeWADS         bool
}

func (f *Flags) parseFlags() {
	printWADMusicInfo := flag.Bool("musicinfo", false, "Print WAD's music info via console")
	printWADMapsInfo := flag.Bool("mapsinfo", false, "Print WAD's maps info via console")
	dumpLumpsInfo := flag.String("lumpsinfo-dump", "", "Dump WAD's lumps info to file")
	dumpWADMusicInfo := flag.String("musicinfo-dump", "", "Dump WAD's music info to file")
	dumpWADMapsInfo := flag.String("mapsinfo-dump", "", "Dump WAD's maps info to file")
	exportMusic := flag.String("music-export", "", "Export WAD's music to folder")
	exportSprites := flag.String("sprite-export", "", "Export WAD's sprites to folder")
	mergeWads := flag.Bool("mergewads", false, "Merge Multiple WADS into one")

	flag.Parse()

	f.printWADMusicInfo = *printWADMusicInfo
	f.printWADMapsInfo = *printWADMapsInfo
	f.dumpLumpsInfo = *dumpLumpsInfo
	f.dumpWADMusicInfo = *dumpWADMusicInfo
	f.dumpWADMapsInfo = *dumpWADMapsInfo
	f.exportMusic = *exportMusic
	f.exportSprites = *exportSprites
	f.mergeWADS = *mergeWads
	f.WADFilenames = flag.Args()
}

func (f *Flags) printFlags() {
	fmt.Printf("%+v\n", f)
}
