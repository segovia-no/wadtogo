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
}
