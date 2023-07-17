package wadloader

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
)

// Graphic structs
type Palette [256]PaletteColor

type PaletteColor struct {
	Red uint8
	Green uint8
	Blue uint8
}

type Patch struct {
	Name string
	Width uint16
	Height uint16
	LeftOffset int16
	TopOffset int16
	PostOffsets []uint32
	PatchPosts []PatchPost
}

type PatchPost []PatchPostSegment

type PatchPostSegment struct {
	TopOffset uint8
	Length uint8
	PixelData []byte
}

// Graphic indexes structs
type spriteMarkerIndexes = struct {
	S_START int
	S_END int
	SS_START int
	SS_END int
}

type flatMarkerIndexes = struct {
	F_START int
	F_END int
	F1_START int
	F1_END int
	F2_START int
	F2_END int
}

type patchesMarkerIndexes = struct {
	P_START int
	P_END int
	P1_START int
	P1_END int
	P2_START int
	P2_END int
	P3_START int
	P3_END int
}

// Graphic functions
func (wl *WADLoader) LoadPalette() {
	wp.checkValidByteReader()

	if len(wl.WADLumps) < 1 {
		fmt.Println("[Warn] LoadPalette: No Lumps loaded, cannot load palette!")
		return
	}

	for _, lump := range wl.WADLumps {
		lumpName := string(bytes.Trim(lump.LumpName[:], "\x00"))

		if string(lumpName) == "PLAYPAL" {
			wp.byteReader.Seek(int64(lump.LumpOffset), io.SeekStart)

			var palettes []Palette

			for i := 0; i < 14; i++ {
				var p Palette

				errRead := binary.Read(wp.byteReader, binary.LittleEndian, &p)
				if errRead != nil {
					fmt.Println("[Error] LoadPalette: Cannot read one of the palettes")
					return
				}

				palettes = append(palettes, p)
			}

			wl.Palettes = palettes
			return
		}
	}
}

func (wl *WADLoader) DetectGraphics() ([]Patch, []Patch, error) {

	var sprites []Patch
	// var flats []Patch // TODO: this is not a patch, but raw data
	var patches []Patch

	if len(wl.WADLumps) < 1 {
		return sprites, patches, errors.New("[Warn] DetectGraphics: No Lumps loaded, cannot detect graphic lumps")
	}

	spriteLumps, err := getSpriteLumps(wl.WADLumps)
	if err != nil {
		return sprites, patches, errors.New("[Warn] DetectGraphics: Cannot get sprite lumps")
	}

	for _, lump := range spriteLumps {
		spritePatch, err := parsePatchLump(lump)
		if err != nil {
			return sprites, patches, errors.New("[Error] DetectGraphics: Cannot parse a sprite patch lump - " + err.Error())
		}

		sprites = append(sprites, spritePatch)
	}

	patchLumps, err := getPatchLumps(wl.WADLumps)
	if err != nil {
		return sprites, patches, errors.New("[Warn] DetectGraphics: Cannot detect patch lumps")
	}

	for _, lump := range patchLumps {
		patch, err := parsePatchLump(lump)
		if err != nil {
			return sprites, patches, errors.New("[Error] DetectGraphics: Cannot parse a patch lump - " + err.Error())
		}

		patches = append(patches, patch)
	}

	return sprites, patches, nil
}

func getSpriteLumps(lumps []Lump) ([]Lump, error) {

	var spriteLumps []Lump

	sprIdx, err := getSpriteMarkerIndexes(lumps)
	if err != nil {
		return spriteLumps, errors.New("[Error] getSpriteLumps: Cannot get sprite markers indexes")
	}

	if sprIdx.S_START != 0 && sprIdx.S_END != 0 {
		for i := sprIdx.S_START + 1; i < sprIdx.S_END; i++ {
			if lumps[i].LumpSize != 0 { // ignore submarkers
				spriteLumps = append(spriteLumps, lumps[i])
			}
		}
	}

	if sprIdx.SS_START != 0 && sprIdx.SS_END != 0 {
		if sprIdx.SS_END < sprIdx.S_START || sprIdx.SS_START > sprIdx.S_END {

			var ssLumps []Lump
			for i := sprIdx.SS_START + 1; i < sprIdx.SS_END; i++ {
				if lumps[i].LumpSize != 0 { // ignore submarkers
					ssLumps = append(ssLumps, lumps[i])
				}
			}

			spriteLumps = append(spriteLumps, ssLumps...)
		}
	}

	return spriteLumps, nil
}

func getPatchLumps(lumps []Lump) ([]Lump, error) {

	var patchLumps []Lump

	patchIdx, err := getPatchMarkerIndexes(lumps)
	if err != nil {
		return patchLumps, errors.New("[Error] getPatchLumps: Cannot get patch markers indexes")
	}

	if patchIdx.P_START != 0 && patchIdx.P_END != 0 {
		for i := patchIdx.P_START + 1; i < patchIdx.P_END; i++ {
			if lumps[i].LumpSize != 0 { // ignore submarkers
				patchLumps = append(patchLumps, lumps[i])
			}
		}
	}

	return patchLumps, nil
}

func parsePatchLump(patchLump Lump) (Patch, error) {
	wp.checkValidByteReader()

	var patch Patch

	if patchLump.LumpSize < 8 {
		return patch, errors.New("[Error] parsePatchLump: Provided lump doesn't have enough bytes to parse header")
	}

	wp.byteReader.Seek(int64(patchLump.LumpOffset), io.SeekStart)

	lumpName := string(bytes.Trim(patchLump.LumpName[:], "\x00"))
	patch.Name = lumpName

	var patchHeader struct {
		Width uint16
		Height uint16
		LeftOffset int16
		TopOffset int16
	}

	errRead := binary.Read(wp.byteReader, binary.LittleEndian, &patchHeader)
	if errRead != nil {
		return patch, errors.New("[Error] parsePatchLump: Cannot parse header info of " + lumpName + " - " + errRead.Error())
	}

	patch.Width = patchHeader.Width
	patch.Height = patchHeader.Height
	patch.LeftOffset = patchHeader.LeftOffset
	patch.TopOffset = patchHeader.TopOffset

	var patchHeaderPostOffsets []uint32
	for i := 0; i < int(patch.Width) - 1; i++ {
		var postOffset uint32

		errRead = binary.Read(wp.byteReader, binary.LittleEndian, &postOffset)
		if errRead != nil {
			return patch, errors.New("[Error] parsePatchLump: Cannot parse a patch post offset of " + lumpName + " - " + errRead.Error())
		}

		patchHeaderPostOffsets = append(patchHeaderPostOffsets, postOffset)
	}
	
	patch.PostOffsets = patchHeaderPostOffsets

	for _, currOffset := range patch.PostOffsets {
		pPost, err := parsePatchPost(patchLump.LumpOffset + currOffset)
		if err != nil {
			return patch, errors.New("[Error] parsePatchLump: Cannot parse a patch post of " + lumpName + " - " + errRead.Error())
		}

		patch.PatchPosts = append(patch.PatchPosts, pPost)
	}

	return patch, nil
}

func parsePatchPost(lumpOffset uint32) (PatchPost, error) {

	var patchPost PatchPost
	var currInnerPostOffset uint32 = lumpOffset

	for {
		patchPostSeg, currOffset, err := parsePatchPostSegment(currInnerPostOffset)

		if err != nil {
			return patchPost, errors.New("[Error] parsePatchPost: Cannot parse a patch post segment - " + err.Error())
		}

		if (patchPostSeg.TopOffset == 255) {
			break
		}

		patchPost = append(patchPost, patchPostSeg)
		currInnerPostOffset = uint32(currOffset)
	}

	return patchPost, nil
}

func parsePatchPostSegment(offset uint32) (PatchPostSegment, int64, error) {
	wp.checkValidByteReader()
	wp.byteReader.Seek(int64(offset), io.SeekStart)

	var pPost PatchPostSegment

	var patchPostHeaderFields struct {
		TopOffset uint8
		Length uint8
		PaddingPre uint8 // ignore data, only use is to move seeker
	}

	errRead := binary.Read(wp.byteReader, binary.LittleEndian, &patchPostHeaderFields)
	if errRead != nil {
		return pPost, 0, errors.New("[Error] parsePatchPostSegment: Cannot parse patch post header data - " + errRead.Error())
	}

	pPost.TopOffset = patchPostHeaderFields.TopOffset
	if pPost.TopOffset == 255 {
		return pPost, 0, nil
	}

	pPost.Length = patchPostHeaderFields.Length

	pixelData := make([]byte, pPost.Length)

	errRead = binary.Read(wp.byteReader, binary.LittleEndian, &pixelData)
	if errRead != nil {
		return pPost, 0, errors.New("[Error] parsePatchPostSegment: Cannot parse patch post pixel data - " + errRead.Error())
	}

	pPost.PixelData = pixelData

	var paddingPost uint8
	errRead = binary.Read(wp.byteReader, binary.LittleEndian, &paddingPost)
	if errRead != nil {
		return pPost, 0, errors.New("[Error] parsePatchPostSegment: Cannot parse patch post segment end padding - " + errRead.Error())
	}

	// get bytereader's current offset so next segment can be read
	currOffset, seekErr := wp.byteReader.Seek(0, io.SeekCurrent)
	if seekErr != nil {
		return pPost, currOffset, errors.New("[Error] parsePatchPostSegment: Cannot get current offset of bytereader - " + seekErr.Error())
	}

	return pPost, currOffset, nil
}

func getSpriteMarkerIndexes(lumps []Lump) (spriteMarkerIndexes, error) {
	var sprIdx spriteMarkerIndexes
	
	for idx, lump := range lumps {
		lumpName := string(bytes.Trim(lump.LumpName[:], "\x00"))
		switch lumpName {
		case "S_START":
			sprIdx.S_START = idx
		case "S_END":
			sprIdx.S_END = idx
		case "SS_START":
			sprIdx.S_START = idx
		case "SS_END":
			sprIdx.SS_END = idx
		}
	}

	if sprIdx.S_START > sprIdx.S_END {
		var emptyIdx spriteMarkerIndexes
		errorStr := "[Error] getSpriteMarkerIndexes: Malformed S_START and S_END index values"
		fmt.Println(errorStr)
		return emptyIdx, errors.New(errorStr)
	}

	if sprIdx.SS_START != 0 && sprIdx.SS_END != 0 {
		if sprIdx.SS_START > sprIdx.SS_END {
			sprIdx.SS_START = 0
			sprIdx.SS_END = 0
			fmt.Println("[Warn] getSpriteMarkerIndexes: Malformed SS_START and SS_END index values, invalidating these indexes")
		}
	}

	return sprIdx, nil
}

func getFlatMarkerIndexes(lumps []Lump) (flatMarkerIndexes, error) {
	var flatIdx flatMarkerIndexes
	
	for idx, lump := range lumps {
		lumpName := string(bytes.Trim(lump.LumpName[:], "\x00"))
		switch lumpName {
		case "F_START":
			flatIdx.F_START = idx
		case "F_END":
			flatIdx.F_END = idx
		case "F1_START":
			flatIdx.F1_START = idx
		case "F1_END":
			flatIdx.F1_END = idx
		case "F2_START":
			flatIdx.F2_START = idx
		case "F2_END":
			flatIdx.F2_END = idx
		}
	}

	if flatIdx.F_START > flatIdx.F_END {
		var emptyIdx flatMarkerIndexes
		errorStr := "[Error] getFlatMarkerIndexes: Malformed F_START and F_END index values"
		fmt.Println(errorStr)
		return emptyIdx, errors.New(errorStr)
	}

	return flatIdx, nil
}

func getPatchMarkerIndexes(lumps []Lump) (patchesMarkerIndexes, error) {
	var patchIdx patchesMarkerIndexes
	
	for idx, lump := range lumps {
		lumpName := string(bytes.Trim(lump.LumpName[:], "\x00"))
		switch lumpName {
		case "P_START":
			patchIdx.P_START = idx
		case "P_END":
			patchIdx.P_END = idx
		case "P1_START":
			patchIdx.P1_START = idx
		case "P1_END":
			patchIdx.P1_END = idx
		case "P2_START":
			patchIdx.P2_START = idx
		case "P2_END":
			patchIdx.P2_END = idx
		case "P3_START":
			patchIdx.P3_START = idx
		case "P3_END":
			patchIdx.P3_END = idx
		}
	}

	if patchIdx.P_START > patchIdx.P_END {
		var emptyIdx patchesMarkerIndexes
		errorStr := "[Error] getPatchMarkerIndexes: Malformed P_START and P_END index values"
		fmt.Println(errorStr)
		return emptyIdx, errors.New(errorStr)
	}

	return patchIdx, nil
}

func (wl *WADLoader) ExportAllSprites(outputFolder string) error {

	if outputFolder == "" {
		return errors.New("[Error] createFolder: No folder name specified")
	}

	lastChar := outputFolder[len(outputFolder)-1:]
	if lastChar == "/" {
		outputFolder = outputFolder[:len(outputFolder)-1]
	}

	_, err := os.Stat(outputFolder)
	if (err != nil) {
		err = os.Mkdir(outputFolder, 0755)
		if (err != nil) {
			return errors.New("[Error] createFolder: Cannot create the target folder")
		}
	}

	for _, sprite := range wl.Sprites {
		exportErr := ExportSprite(sprite, wl.Palettes[0], outputFolder)
		if exportErr != nil {
			return errors.New("[Error] ExportAllSprites: Cannot export sprite - " + sprite.Name + " - " + exportErr.Error())
		}
	}

	for _, patchSprites := range wl.Patches {
		exportErr := ExportSprite(patchSprites, wl.Palettes[0], outputFolder)
		if exportErr != nil {
			return errors.New("[Error] ExportAllSprites: Cannot export patch sprite - " + patchSprites.Name + " - " + exportErr.Error())
		}
	}

	return nil
}

func ExportSprite(sprite Patch, palette Palette, outputFolder string) error {
	spriteImg := image.NewRGBA(image.Rect(0, 0, int(sprite.Width), int(sprite.Height)))

	for idx, post := range sprite.PatchPosts {
		for _, postSegment := range post {
			for i := 0; i < int(postSegment.Length); i++ {
				x := idx
				y := int(postSegment.TopOffset) + i
				pixelColor := palette[postSegment.PixelData[i]]

				spriteImg.Set(x, y, color.RGBA{pixelColor.Red, pixelColor.Green, pixelColor.Blue, 255})
			}
		}
	}

	spriteFile, err := os.Create(outputFolder + "/" + sprite.Name + ".png")
	if err != nil {
		return errors.New("[Error] ExportSprite: Cannot create the target file for the sprite - " + sprite.Name + " - " + err.Error())
	}

	defer spriteFile.Close()
	png.Encode(spriteFile, spriteImg)

	return nil
}
