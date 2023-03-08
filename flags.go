package main

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	WADFilename string
	printWADMusicInfo bool
	printWADMapsInfo bool
	dumpLumpsInfo string
	dumpWADMusicInfo string
	dumpWADMapsInfo string
	exportMusic string
}

func (f *Flags) parseFlags() {
	printWADMusicInfo  := flag.Bool("musicinfo", false, "Print WAD's music info via console")
	printWADMapsInfo   := flag.Bool("mapsinfo", false, "Print WAD's maps info via console")
	dumpLumpsInfo      := flag.String("lumpsinfo-dump", "", "Dump WAD's lumps info to file")
	dumpWADMusicInfo   := flag.String("musicinfo-dump", "", "Dump WAD's music info to file")
	dumpWADMapsInfo    := flag.String("mapsinfo-dump", "", "Dump WAD's maps info to file")
	exportMusic        := flag.String("music-export", "", "Export WAD's music to file")

	flag.Parse()

	f.printWADMusicInfo = *printWADMusicInfo
	f.printWADMapsInfo  = *printWADMapsInfo
	f.dumpLumpsInfo     = *dumpLumpsInfo
	f.dumpWADMusicInfo  = *dumpWADMusicInfo
	f.dumpWADMapsInfo   = *dumpWADMapsInfo
	f.exportMusic       = *exportMusic

	flagTail := flag.Args()
	f.parseWADFilename(flagTail)
}

func (f *Flags) parseWADFilename(flagTail []string) {
	if len(flagTail) < 1 {
		fmt.Println("[Error] getWADFilenameFromFlag: No WAD filename argument specified. Aborting extraction")
		os.Exit(0)
	}

	filename := flag.Args()[0]

	if filename == "" {
		fmt.Println("[Error] getWADFilenameFromFlag: No WAD filename argument specified. Aborting extraction")
		os.Exit(0)
	}

	f.WADFilename = filename
}

func (f *Flags) printFlags() {
	fmt.Printf("%+v\n", f)
}
