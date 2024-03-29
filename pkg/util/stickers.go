package util

import (
	"TuruBot/pkg/errorHandling"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
)

func ImageToWEBP(inputData []byte, mimetype string, updateId int64) ([]byte, bool, error) {
	var (
		currUpdateId 	= strconv.FormatInt(updateId, 10)
		currPath     	= path.Join("temp", currUpdateId)
		inputPath    	= path.Join(currPath, "input.jpg")
		outputPath   	= path.Join(currPath, "output.webp")
	)
 
	os.MkdirAll(currPath, os.ModePerm)
	if err := os.WriteFile(inputPath, inputData, os.ModePerm); err != nil {
		errorHandling.LogErr(err)
		return nil, false, err
	}
 
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-quality", "60", "-loop", "0", "-fs", "300k", "-vcodec", "libwebp", "-vf", "scale='min(512,iw)':'min(512,ih)':force_original_aspect_ratio=decrease,fps=23,pad=512:512:-1:-1:color=white@0.0,split[a][b];[a]palettegen=reserve_transparent=on:transparency_color=ffffff[p];[b][p]paletteuse", "-f", "webp", outputPath)
	if err := cmd.Run(); err != nil {
		errorHandling.LogErr(err)
		return nil, false, err
	}
 
	imgBytes, err := os.ReadFile(outputPath)
	if err != nil {
		errorHandling.LogErr(err)
		return nil, false, err
	}
 
	var IsAnimated bool = false
	if http.DetectContentType(inputData) == "image/GIF" || http.DetectContentType(inputData) == "image/gif" {
	    	IsAnimated = true
	}

	// delete temp files
	os.RemoveAll(currPath)
 
	return imgBytes, IsAnimated, nil
}


func VideoToWEBP(inputData []byte, updateId int64) ([]byte, error) {
	var (
		currUpdateId = strconv.FormatInt(updateId, 10)
		currPath     = path.Join("temp", currUpdateId)
		inputPath    = path.Join(currPath, "input.mp4")
		outputPath   = path.Join(currPath, "output.webp")
	)

	os.MkdirAll(currPath, os.ModePerm)
	if err := os.WriteFile(inputPath, inputData, os.ModePerm); err != nil {
		errorHandling.LogErr(err)
		return nil, err
	}

	// Menghitung bitrate
	bitrate := uint64(50)// CalculateBitrate(fileLength, 300000, duration)

	// Convert video ke webp
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-loop", "0", "-vf", "scale='min(512,iw)':'min(512,ih)':force_original_aspect_ratio=decrease,fps=10,pad=512:512:-1:-1:color=white@0.0,split[a][b];[a]palettegen=reserve_transparent=on:transparency_color=ffffff[p];[b][p]paletteuse", "-lossless", "1", "-preset", "default", "-quality", "50", "-b:v", strconv.FormatUint(bitrate, 10) + "k", outputPath)
	if err := cmd.Run(); err != nil {
		errorHandling.LogErr(err)
		return nil, err
	}

	videoBytes, err := os.ReadFile(outputPath)
	if err != nil {
		errorHandling.LogErr(err)
		return nil, err
	}

	// delete temp files
	os.RemoveAll(currPath)

	return videoBytes, nil
}

func CalculateBitrate(fileLength *uint64, targetSize uint64, durationSeconds uint64) uint64 {
	// Target size in bits
	targetSizeBits := targetSize * 8

	// Calculate total bits needed for target size
	totalBits := targetSizeBits * uint64(durationSeconds)

	// Calculate bitrate in bits per second
	bitrate := uint64(math.Round(float64(totalBits) / float64(durationSeconds)))

	return bitrate
}