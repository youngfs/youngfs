package api

import (
	"compress/flate"
	"compress/gzip"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/bucket"
	"youngfs/fs/fullpath"
	"youngfs/fs/server"
)

func PutObjectHandler(c *gin.Context) {
	var file io.ReadCloser
	var filename string
	file, head, err := c.Request.FormFile("file")
	if head != nil && err == nil {
		filename = head.Filename
		switch head.Header.Get("Content-Encoding") {
		case "gzip":
			file, err = gzip.NewReader(file)
			if err != nil {
				errorHandler(c, errors.ErrContentEncoding.WrapErrNoStack(err))
				return
			}
		case "deflate":
			file = flate.NewReader(file)
		case "br":
			file = io.NopCloser(brotli.NewReader(file))
		}
	}
	if err != nil {
		file = c.Request.Body
		filename = ""
		switch head.Header.Get("Content-Encoding") {
		case "gzip":
			file, err = gzip.NewReader(file)
			if err != nil {
				errorHandler(c, errors.ErrContentEncoding.WrapErrNoStack(err))
				return
			}
		case "deflate":
			file = flate.NewReader(file)
		case "br":
			file = io.NopCloser(brotli.NewReader(file))
		}
	}
	defer func() {
		_ = file.Close()
	}()

	bkt, fp := bucket.Bucket(c.Param("bucket")), fullpath.FullPath(c.Param("path"))
	if len(fp) == 0 || fp[len(fp)-1] == '/' {
		fp += fullpath.FullPath(filename)
	}
	if !bkt.IsLegal() {
		errorHandler(c, errors.ErrIllegalBucketName)
		return
	}
	if !fp.IsLegalObjectName() {
		errorHandler(c, errors.ErrIllegalObjectName)
		return
	}
	fp = fp.Clean()

	err = server.PutObject(c, bkt, fp, file)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.Status(http.StatusCreated)
	return
}
