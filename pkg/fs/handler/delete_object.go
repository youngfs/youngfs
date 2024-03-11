package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"net/http"
)

func (h *Handler) DeleteObjectHandler(c *gin.Context) {
	bkt, fp := bucket.Bucket(c.Param(QueryBucketKey)), fullpath.FullPath(c.Param(QueryPathKey))
	if !bkt.IsLegal() {
		h.errorHandler(c, errors.ErrIllegalBucketName)
		return
	}
	if !fp.IsLegal() {
		h.errorHandler(c, errors.ErrIllegalObjectName)
		return
	}
	fp = fp.Clean()

	err := h.svr.DeleteObject(c, bkt, fp)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	return
}
