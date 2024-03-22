package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"strconv"
)


func WebpWriteExifData(inputData []byte, updateId int64) ([]byte, error) {
	var (
		startingBytes = []byte{0x49, 0x49, 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x41, 0x57, 0x07, 0x00}
		endingBytes   = []byte{0x16, 0x00, 0x00, 0x00}
		b             bytes.Buffer

		currUpdateId = strconv.FormatInt(updateId, 10)
		currPath     = path.Join("downloads", currUpdateId)
		inputPath    = path.Join(currPath, "input_exif.jpg")
		outputPath   = path.Join(currPath, "output_exif.webp")
		exifDataPath = path.Join(currPath, "raw.exif")
	)

	b.Write(startingBytes)

	jsonData := map[string]interface{}{
		"sticker-pack-id":        "com.turudev.my.id",
		"sticker-pack-name":      "nopaleon",
		"sticker-pack-publisher": "nopaleon",
		"emojis":                 []string{"😀"},
	}
	jsonBytes, _ := json.Marshal(jsonData)

	jsonLength := uint32(len(jsonBytes))
	lenBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuffer, jsonLength)
	b.Write(lenBuffer)
	b.Write(endingBytes)
	b.Write(jsonBytes)

	os.MkdirAll(currPath, os.ModePerm)
	os.WriteFile(inputPath, inputData, os.ModePerm)
	os.WriteFile(exifDataPath, b.Bytes(), os.ModePerm)

	cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libwebp", "-vf", "scale='min(512,iw)':'min(512,ih)':force_original_aspect_ratio=decrease,fps=23,pad=512:512:-1:-1:color=white@0.0,split[a][b];[a]palettegen=reserve_transparent=on:transparency_color=ffffff[p];[b][p]paletteuse", "-f", "webp", outputPath)
	cmd.Run()

	imgBytes, _ := os.ReadFile(outputPath)

	return imgBytes, nil
}
