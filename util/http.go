package util

import (
	"net/http"
	"strconv"
)

func GetContentLength(h http.Header) uint64 {
	contentLength := h.Get("Content-Length")
	if contentLength != "" {
		length, err := strconv.ParseUint(contentLength, 10, 64)
		if err != nil {
			return 0
		}
		return length
	}
	return 0
}
