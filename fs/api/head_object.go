package api

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"youngfs/errors"
	"youngfs/fs/bucket"
	"youngfs/fs/fullpath"
	"youngfs/fs/server"
)

func HeadObjectHandler(c *gin.Context) {
	bkt, fp := bucket.Bucket(c.Param("bucket")), fullpath.FullPath(c.Param("path"))
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
