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

func ImageToWEBP(inputData []byte, updateId int64) ([]byte, []byte, bool, error) {
	var (
		currUpdateId 	= strconv.FormatInt(updateId, 10)
		currPath     	= path.Join("downloads", currUpdateId)
		inputPath    	= path.Join(currPath, "input.jpg")
		outputPath   	= path.Join(currPath, "output.webp")
		thumbPath    	= path.Join(currPath, "thumbnail.png")
	)
 
	os.MkdirAll(currPath, os.ModePerm)
	if err := os.WriteFile(inputPath, inputData, os.ModePerm); err != nil {
		errorHandling.LogErr(err)
		return nil, nil, false, err
	}
 
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-quality", "60", "-loop", "0", "-fs", "300k", "-vcodec", "libwebp", "-vf", "scale='min(512,iw)':'min(512,ih)':force_original_aspect_ratio=decrease,fps=23,pad=512:512:-1:-1:color=white@0.0,split[a][b];[a]palettegen=reserve_transparent=on:transparency_color=ffffff[p];[b][p]paletteuse", "-f", "webp", outputPath)
	if err := cmd.Run(); err != nil {
		errorHandling.LogErr(err)
		return nil, nil, false, err
	}
 
	imgBytes, err := os.ReadFile(outputPath)
	if err != nil {
		errorHandling.LogErr(err)
		return nil, nil, false, err
	}
 
	var IsAnimated bool = false
	var thumbnailData []byte
	if http.DetectContentType(inputData) == "image/GIF" || http.DetectContentType(inputData) == "image/gif" {
	    	IsAnimated = true
	}

	thumbnailCmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "thumbnail,scale=150:150", "-frames:v", "1", thumbPath)
	if err := thumbnailCmd.Run(); err != nil {
		errorHandling.LogErr(err)
		return nil, nil, false, err
	}

	thumbnailData, err = os.ReadFile(thumbPath)
	if err != nil {
		errorHandling.LogErr(err)
		return nil, nil, false, err
	}
 
	return imgBytes, thumbnailData, IsAnimated, nil
} 


func VideoToWEBP(inputData []byte, updateId int64, fileLength *uint64, duration uint64) ([]byte, []byte, error) {
	var (
		currUpdateId = strconv.FormatInt(updateId, 10)
		currPath     = path.Join("downloads", currUpdateId)
		inputPath    = path.Join(currPath, "input.mp4")
		outputPath   = path.Join(currPath, "output.webp")
		thumbnailPath = path.Join(currPath, "thumbnail.jpg")
	)

	os.MkdirAll(currPath, os.ModePerm)
	if err := os.WriteFile(inputPath, inputData, os.ModePerm); err != nil {
		errorHandling.LogErr(err)
		return nil, nil, err
	}

	// Menghitung bitrate
	bitrate := uint64(50)// CalculateBitrate(fileLength, 300000, duration)

	// Convert video ke webp
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-loop", "0", "-vf", "scale='min(512,iw)':'min(512,ih)':force_original_aspect_ratio=decrease,fps=10,pad=512:512:-1:-1:color=white@0.0,split[a][b];[a]palettegen=reserve_transparent=on:transparency_color=ffffff[p];[b][p]paletteuse", "-lossless", "1", "-preset", "default", "-quality", "50", "-b:v", strconv.FormatUint(bitrate, 10) + "k", outputPath)
	if err := cmd.Run(); err != nil {
		errorHandling.LogErr(err)
		return nil, nil, err
	}

	// Mendapatkan thumbnail
	thumbnailCmd := exec.Command("ffmpeg", "-i", inputPath, "-ss", "00:00:01.000", "-vf", "thumbnail,scale=150:150", "-vframes", "1", thumbnailPath)
	if err := thumbnailCmd.Run(); err != nil {
		errorHandling.LogErr(err)
		return nil, nil, err
	}

	thumbnailData, err := os.ReadFile(thumbnailPath)
	if err != nil {
		errorHandling.LogErr(err)
		return nil, nil, err
	}

	videoBytes, err := os.ReadFile(outputPath)
	if err != nil {
		errorHandling.LogErr(err)
		return nil, nil, err
	}

	return videoBytes, thumbnailData, nil
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