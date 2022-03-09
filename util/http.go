package util

import (
	"net/http"
	"strconv"
)

func GetContentLength(r *http.Request) uint64 {
	contentLength := r.Header.Get("Content-Length")
	if contentLength != "" {
		length, err := strconv.ParseUint(contentLength, 10, 64)
		if err != nil {
			return 0
		}
		return length
	}
	return 0
}
