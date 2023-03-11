package wadloader

import (
	"fmt"
)

func PrintMapNames(maps []Map) {
	fmt.Println("Map List")
	for _, m := range maps {
		fmt.Println(m.Name)
	}
}

func PrintSongNames(songs []MusicLump) {
	fmt.Println("Song list | Format")
	for _, m := range songs {
		fmt.Println(m.name + " | " + m.format)
	}
}
