package wadloader

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

func (wp *WADParser) ExportSong(song *MusicLump) error {
	wp.checkValidByteReader()
	wp.byteReader.Seek(int64(song.lump.LumpOffset), io.SeekStart)

	filename := song.name + "." + song.format

	os.Remove(filename)
	file, err := os.Create(filename)
	
	if err != nil {
		return errors.New("[Error] ExportSong: Cannot create the target file")
	}
	
	defer file.Close()

	lumpData := make([]byte, song.lump.LumpSize)
	errRead := binary.Read(wp.byteReader, binary.LittleEndian, &lumpData)

	if errRead != nil {
		return errors.New("[Error] ExportSong: Cannot read the song lump data")
	}

	binary.Write(file, binary.LittleEndian, lumpData)

	return nil
}

func (wl *WADLoader) ExportAllSongs(wp *WADParser) error {
	wp.checkValidByteReader()

	if len(wl.Music) < 1 {
		return errors.New("[Error] ExportAllSongs: No music data inside WAD Loader")
	}

	for _, song := range wl.Music {
		wp.ExportSong(&song)
	}

	return nil
}

func (wl *WADLoader) GetMusicLumpFromSongName(songName string) (*MusicLump, error) {

	var musicLump MusicLump

	if len(wl.Music) < 1 {
		return &musicLump, errors.New("[Error] GetMusicLumpFromSongName: No music data inside WAD Loader")
	}

	for _, ml := range wl.Music {
		if ml.name == songName {
			return &ml, nil
		}
	}

	return &musicLump, errors.New("[Error] GetMusicLumpFromSongName: Song name not found")
}