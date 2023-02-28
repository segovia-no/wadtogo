package main

import (
	"fmt"

	wl "github.com/segovia-no/wadtogo/wadloader"
)

var flagReader FlagReader
var wadLoader wl.WADLoader

func main() {
	fmt.Println("WADToGo - Another WAD Tool")
	wadFilename := flagReader.getWADFilenameFromFlag()
	wadLoader.OpenAndLoad(wadFilename)
}
