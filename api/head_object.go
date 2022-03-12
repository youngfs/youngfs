package api

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"icesos/errors"
	"icesos/full_path"
	"icesos/server"
	"icesos/set"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func HeadObjectHandler(c *gin.Context) {
	setName, fp := set.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
	if !setName.IsLegal() {
		err := errors.ErrorCodeResponse[errors.ErrIllegalSetName]
		c.Status(err.HTTPStatusCode)
		return
	}
	if !fp.IsLegal() {
		err := errors.ErrorCodeResponse[errors.ErrIllegalObjectName]
		c.Status(err.HTTPStatusCode)
		return
	}
	fp = fp.Clean()

	ent, err := server.Svr.GetObject(c, setName, fp)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.ErrorCodeResponse[errors.ErrServer]
		}
		c.Status(err.HTTPStatusCode)
		return
	}

	// url encode but not code /
	c.Header("Full-Path", (&url.URL{Path: string(ent.FullPath)}).String())
	// url encode and code /
	c.Header("Set", url.PathEscape(string(ent.Set)))
	c.Header("Last-Modified-Time", ent.Mtime.Format(time.RFC3339))
	c.Header("Creation-Time", ent.Ctime.Format(time.RFC3339))
	c.Header("Mode", strconv.FormatUint(uint64(ent.Mode), 10))
	c.Header("Mime", ent.Mime)
	c.Header("Md5", hex.EncodeToString(ent.Md5))
	c.Header("File-Size", strconv.FormatUint(ent.FileSize, 10))
	c.Header("Fid", ent.Fid)
	c.Status(http.StatusOK)
	return
}
