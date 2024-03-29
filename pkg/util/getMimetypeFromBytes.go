package util

import (
	"net/http"
	"strings"
)

func GetFileExtFromBytes(content []byte) (string) {
	mimetype := strings.Split(http.DetectContentType(content), "/")[1]

	if len(mimetype) > 0 {
		return mimetype
	} else {
		return "jpeg"
	}
}