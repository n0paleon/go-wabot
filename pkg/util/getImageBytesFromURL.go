package util

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// fungsi untuk mendapatkan bytes images
// return 3 tipe data

func GetImageBytes(url string) ([]byte, string, error) {
	// Mengirim GET request ke URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	// Membaca response body ke dalam byte array
	imageBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Mendeteksi MIME type dari response body
	mimeType := http.DetectContentType(imageBytes)

	// Memeriksa apakah MIME type termasuk kategori gambar
	if !strings.HasPrefix(mimeType, "image/") {
		return nil, "", errors.New("response bukan berupa gambar")
	}

	return imageBytes, mimeType, nil
}