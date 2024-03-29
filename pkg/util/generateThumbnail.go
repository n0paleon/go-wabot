package util

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func GenerateStickerThumbnail(image []byte, mimetype string, outputMimetype string) ([]byte, error) {
    // Buat folder temp jika belum ada
    tempFolder := "temp/thumb-" + strconv.FormatInt(time.Now().Unix(), 10) + "/"
    err := os.MkdirAll(tempFolder, os.ModePerm)
    if err != nil {
        return nil, err
    }

    // Simpan file sementara dengan mimetype yang sesuai
    inputPath := tempFolder + "input." + mimetype
    err = ioutil.WriteFile(inputPath, image, 0644)
    if err != nil {
        return nil, err
    }

    // Tentukan output path dengan mimetype yang diinginkan
    outputPath := tempFolder + "output." + outputMimetype

    // Buat command ffmpeg
    cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "thumbnail,scale=72:72", "-frames:v", "1", outputPath)

    // Jalankan command
    err = cmd.Run()
    if err != nil {
        return nil, err
    }

    // Baca file thumbnail
    thumbnail, err := ioutil.ReadFile(outputPath)
    if err != nil {
        return nil, err
    }

    // Hapus folder temp dan isinya setelah thumbnail dibuat
    os.RemoveAll(tempFolder)

    return thumbnail, nil
}

func GenerateMediaThumbnail(image []byte, mimetype string, outputMimetype string) ([]byte, error) {
    // Buat folder temp jika belum ada
    tempFolder := "temp/thumb-" + strconv.FormatInt(time.Now().Unix(), 10) + "/"
    err := os.MkdirAll(tempFolder, os.ModePerm)
    if err != nil {
        return nil, err
    }

    // Simpan file sementara dengan mimetype yang sesuai
    inputPath := tempFolder + "input." + mimetype
    err = ioutil.WriteFile(inputPath, image, 0644)
    if err != nil {
        return nil, err
    }

    // Tentukan output path dengan mimetype yang diinginkan
    outputPath := tempFolder + "output." + outputMimetype

    // Buat command ffmpeg
    cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "thumbnail,scale=72:72", "-frames:v", "1", outputPath)

    // Jalankan command
    err = cmd.Run()
    if err != nil {
        return nil, err
    }

    // Baca file thumbnail
    thumbnail, err := ioutil.ReadFile(outputPath)
    if err != nil {
        return nil, err
    }

    // Hapus folder temp dan isinya setelah thumbnail dibuat
    os.RemoveAll(tempFolder)

    return thumbnail, nil
}