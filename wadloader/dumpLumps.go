package wadloader

import (
	"bytes"
	"fmt"
	"os"
)

func DumpLumpsToTextFile(filename string, lumps WADLumps) {
	os.Remove(filename)
	file, err := os.Create(filename)

	if err != nil {
		fmt.Println("[Error] Cannot create a file called", filename, err)
		return
	}

	defer file.Close()

	var errWrite error

	_, errWrite = file.WriteString("Lump name | Size (bytes)\n")

	for _, lump := range lumps {

		lumpName := bytes.Trim(lump.LumpName[:], "\x00")

		outStr := fmt.Sprintf("%s %v \n", lumpName, lump.LumpSize)
		_, errWrite = file.WriteString(outStr)
	}

	if errWrite != nil {
		fmt.Println("[Error] Cannot add lump data to dump file", filename, err)
		return
	}

	fmt.Println("[Info] Lumps dumped into", filename)
}

func DumpMapNamesToTextFile(filename string, maps []MapRawLumps) {
	os.Remove(filename)
	file, err := os.Create(filename)

	if err != nil {
		fmt.Println("[Error] Cannot create a file called", filename, err)
		return
	}

	defer file.Close()

	var errWrite error

	_, errWrite = file.WriteString("Map list\n")

	for _, m := range maps {
		_, errWrite = file.WriteString(m.MapName + "\n")
	}

	if errWrite != nil {
		fmt.Println("[Error] Cannot add map data to dump file", filename, err)
		return
	}

	fmt.Println("[Info] Map list dumped into", filename)
}

func DumpSongNamesToTextFile(filename string, songs []MusicLump) {
	os.Remove(filename)
	file, err := os.Create(filename)

	if err != nil {
		fmt.Println("[Error] Cannot create a file called", filename, err)
		return
	}

	defer file.Close()

	var errWrite error

	_, errWrite = file.WriteString("Song list | Format \n")

	for _, m := range songs {
		_, errWrite = file.WriteString(m.name + " | " + m.format + "\n")
	}

	if errWrite != nil {
		fmt.Println("[Error] Cannot add song data to dump file", filename, err)
		return
	}

	fmt.Println("[Info] Song list dumped into", filename)
}
