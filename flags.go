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
	dumpWADInfo string
	dumpWADMusicInfo string
	dumpWADMapsInfo string
}

func (f *Flags) parseFlags() {
	printWADMusicInfo := flag.Bool("musicinfo", false, "Print WAD's music info via console")
	printWADMapsInfo  := flag.Bool("mapsinfo", false, "Print WAD's maps info via console")
	dumpWADInfo       := flag.String("wadinfo-dump", "", "Dump WAD's info to file")
	dumpWADMusicInfo  := flag.String("musicinfo-dump", "", "Dump WAD's music info to file")
	dumpWADMapsInfo   := flag.String("mapsinfo-dump", "", "Dump WAD's maps info to file")

	flag.Parse()

	f.printWADMusicInfo = *printWADMusicInfo
	f.printWADMapsInfo = *printWADMapsInfo
	f.dumpWADInfo = *dumpWADInfo
	f.dumpWADMusicInfo = *dumpWADMusicInfo
	f.dumpWADMapsInfo = *dumpWADMapsInfo

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
