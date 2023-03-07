package wadloader

import (
	"fmt"
)

func printMapName(maps []MapLump) {
	fmt.Println("Map List")
	for _, m := range maps {
		fmt.Println(m.MapName)
	}
}

func printSongNames(songs []MusicLump) {
	fmt.Println("Song list | Format")
	for _, m := range songs {
		fmt.Println(m.name + " | " + m.format)
	}
}
