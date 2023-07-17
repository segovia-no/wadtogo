package wadloader

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

func (wp *WADParser) ExportSong(song *MusicLump, outputFolder string) error {
	wp.checkValidByteReader()
	wp.byteReader.Seek(int64(song.lump.LumpOffset), io.SeekStart)

	filename := song.name + "." + song.format
	finalPath := outputFolder + "/" + filename

	os.Remove(finalPath)
	file, err := os.Create(finalPath)
	
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

func (wl *WADLoader) ExportAllSongs(folderName string) error {
	wp.checkValidByteReader()

	if len(wl.Music) < 1 {
		return errors.New("[Error] ExportAllSongs: No music data inside WAD Loader")
	}

	if folderName == "" {
		return errors.New("[Error] ExportAllSongs: No folder name specified")
	}

	lastChar := folderName[len(folderName)-1:]
	if lastChar == "/" {
		folderName = folderName[:len(folderName)-1]
	}

	_, err := os.Stat(folderName)
	if (err != nil) {
		err = os.Mkdir(folderName, 0755)
		if (err != nil) {
			return errors.New("[Error] ExportAllSongs: Cannot create the target folder")
		}
	}

	for _, song := range wl.Music {
		wp.ExportSong(&song, folderName)
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
