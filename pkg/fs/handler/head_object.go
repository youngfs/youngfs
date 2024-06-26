package handler

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (h *Handler) HeadObjectHandler(c *gin.Context) {
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

	ent, err := h.svr.GetEntry(c, bkt, fp)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	// url encode but not code /
	c.Header("Full-Path", (&url.URL{Path: string(ent.FullPath)}).String())
	// url encode and code /
	c.Header("Bucket", url.PathEscape(string(ent.Bucket)))
	c.Header("Last-Modified-Time", ent.Mtime.Format(time.RFC3339))
	c.Header("Creation-Time", ent.Ctime.Format(time.RFC3339))
	c.Header("Mode", strconv.FormatUint(uint64(ent.Mode), 10))
	c.Header("Mime", ent.Mime)
	c.Header("Md5", hex.EncodeToString(ent.Md5))
	c.Header("File-Size", strconv.FormatUint(ent.FileSize, 10))
	c.Status(http.StatusOK)
	return
}
