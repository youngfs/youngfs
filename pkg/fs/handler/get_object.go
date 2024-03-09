package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"net/http"
	"net/url"
)

func (h *Handler) GetObjectHandler(c *gin.Context) {
	bkt, fp := bucket.Bucket(c.Param(QueryBucketKey)), fullpath.FullPath(c.Param(QueryPathKey))
	if fp == "/" || fp[len(fp)-1] == '/' {
		//redirect to list object
		h.ListObjectsHandler(c)
		return
	}

	if !bkt.IsLegal() {
		h.errorHandler(c, errors.ErrIllegalBucketName)
		return
	}
	if !fp.IsLegal() {
		h.errorHandler(c, errors.ErrIllegalObjectName)
		return
	}
	fp = fp.Clean()

	ent, err := h.svr.GetEntry(c, bkt, fp)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	if ent.IsDirectory() {
		//redirect to list object
		h.ListObjectsHandler(c)
		return
	}

	err = h.svr.GetObject(c, ent, c.Writer)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(ent.Name())))
	c.Status(http.StatusOK)
	return
}
