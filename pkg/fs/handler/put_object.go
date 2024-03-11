package handler

import (
	"compress/flate"
	"compress/gzip"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"io"
	"net/http"
	"strings"
)

func (h *Handler) PutObjectHandler(c *gin.Context) {
	var file io.ReadCloser
	var filename string
	multipart, err := c.Request.MultipartReader()
	if err == nil {
		for {
			part, err := multipart.NextPart()
			if err != nil {
				h.errorHandler(c, errors.ErrRouter.WarpErr(err))
				return
			}
			if strings.EqualFold(part.FormName(), "file") {
				file = part
				switch part.Header.Get("Content-Encoding") {
				case "gzip":
					file, err = gzip.NewReader(file)
					if err != nil {
						h.errorHandler(c, errors.ErrContentEncoding.WarpErr(err))
						return
					}
				case "deflate":
					file = flate.NewReader(file)
				case "br":
					file = io.NopCloser(brotli.NewReader(file))
				}
				filename = part.FileName()
				break
			}
		}
	}
	if err != nil {
		file = c.Request.Body
		filename = ""
		switch c.Request.Header.Get("Content-Encoding") {
		case "gzip":
			file, err = gzip.NewReader(file)
			if err != nil {
				h.errorHandler(c, errors.ErrContentEncoding.WarpErr(err))
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

	bkt, fp := bucket.Bucket(c.Param(QueryBucketKey)), fullpath.FullPath(c.Param(QueryPathKey))
	if len(fp) == 0 || fp[len(fp)-1] == '/' {
		fp += fullpath.FullPath(filename)
	}
	if !bkt.IsLegal() {
		h.errorHandler(c, errors.ErrIllegalBucketName)
		return
	}
	if !fp.IsLegalObjectName() {
		h.errorHandler(c, errors.ErrIllegalObjectName)
		return
	}
	fp = fp.Clean()

	err = h.svr.PutObject(c, bkt, fp, file)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.Status(http.StatusCreated)
	return
}
