package util

import (
	"bytes"
	"io"
	"net/http"
)

func MimeDetect(file io.Reader) (string, io.Reader) {
	mimeBuffer := make([]byte, 512)
	size, _ := file.Read(mimeBuffer)
	if size > 0 {
		mime := http.DetectContentType(mimeBuffer[:size])
		return mime, io.MultiReader(bytes.NewReader(mimeBuffer[:size]), file)
	}
	return "", file
}
