package api

import (
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/errors"
	"github.com/youngfs/youngfs/fs/bucket"
	"github.com/youngfs/youngfs/fs/fullpath"
	"github.com/youngfs/youngfs/fs/server"
	"net/http"
)

func GetObjectHandler(c *gin.Context) {
	bkt, fp := bucket.Bucket(c.Param("bucket")), fullpath.FullPath(c.Param("path"))
	if fp == "/" || fp[len(fp)-1] == '/' {
		//redirect to list object
		ListObjectsHandler(c)
		return
	}

	if !bkt.IsLegal() {
		errorHandler(c, errors.ErrIllegalBucketName)
		return
	}
	if !fp.IsLegal() {
		errorHandler(c, errors.ErrIllegalObjectName)
		return
	}
	fp = fp.Clean()

	ent, err := server.GetEntry(c, bkt, fp)
	if err != nil {
		errorHandler(c, err)
		return
	}

	if ent.IsDirectory() {
		//redirect to list object
		ListObjectsHandler(c)
		return
	}

	err = server.GetObject(c, ent, c.Writer)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.Status(http.StatusOK)
	return
}
