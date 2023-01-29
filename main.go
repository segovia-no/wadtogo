package main

var flagReader FlagReader
var wadLoader WADLoader

func main() {
	wadFilename := flagReader.getWADFilenameFromFlag()
	wadLoader.openAndLoad(wadFilename)
}