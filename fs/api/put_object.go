package api

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/bucket"
	"youngfs/fs/fullpath"
	"youngfs/fs/server"
	"youngfs/log"
	"youngfs/vars"
)

func PutObjectHandler(c *gin.Context) {
	var file io.ReadCloser
	var filename string
	file, head, err := c.Request.FormFile("file")
	if head != nil && err == nil {
		filename = head.Filename
		if head.Header.Get("Content-Encoding") == "gzip" {
			file, err = gzip.NewReader(file)
			if err != nil {
				err := errors.ErrContentEncoding
				c.Set(vars.CodeKey, err.ErrorCode)
				c.Set(vars.ErrorKey, err.Error())
				c.JSON(
					err.HTTPStatusCode,
					gin.H{
						vars.UUIDKey:  c.Value(vars.UUIDKey),
						vars.CodeKey:  err.ErrorCode,
						vars.ErrorKey: err.Error(),
					},
				)
				return
			}
		}
	}
	if err != nil {
		file = c.Request.Body
		filename = ""
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			file, err = gzip.NewReader(file)
			if err != nil {
				err := errors.ErrContentEncoding
				c.Set(vars.CodeKey, err.ErrorCode)
				c.Set(vars.ErrorKey, err.Error())
				c.JSON(
					err.HTTPStatusCode,
					gin.H{
						vars.UUIDKey:  c.Value(vars.UUIDKey),
						vars.CodeKey:  err.ErrorCode,
						vars.ErrorKey: err.Error(),
					},
				)
				return
			}
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
		err := errors.ErrIllegalBucketName
		c.Set(vars.CodeKey, err.ErrorCode)
		c.Set(vars.ErrorKey, err.Error())
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  err.ErrorCode,
				vars.ErrorKey: err.Error(),
			},
		)
		return
	}
	if !fp.IsLegalObjectName() {
		err := errors.ErrIllegalObjectName
		c.Set(vars.CodeKey, err.ErrorCode)
		c.Set(vars.ErrorKey, err.Error())
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  err.ErrorCode,
				vars.ErrorKey: err.Error(),
			},
		)
		return
	}
	fp = fp.Clean()

	err = server.PutObject(c, bkt, fp, file)
	if err != nil {
		apiErr := &errors.APIError{}
		if !errors.As(err, &apiErr) {
			log.Errorw("a non api error is returned", vars.ErrorKey, err.Error())
			apiErr = errors.ErrNonApiErr
		}
		if apiErr.IsServerErr() {
			log.Errorf("uuid:%s\n error:%v\n", c.Value(vars.UUIDKey), err)
		}
		c.Set(vars.CodeKey, apiErr.ErrorCode)
		c.Set(vars.ErrorKey, apiErr.Error())
		c.JSON(
			apiErr.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  apiErr.ErrorCode,
				vars.ErrorKey: apiErr.Error(),
			},
		)
		return
	}

	c.Status(http.StatusCreated)
	return
}
