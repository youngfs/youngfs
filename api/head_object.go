package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"icesos/errors"
	"icesos/full_path"
	"icesos/server"
	"icesos/set"
	"icesos/util"
	"net/http"
	"strconv"
)

const timeFormat = "2006-01-02 15:04:05"

func HeadObjectHandler(c *gin.Context) {
	ctx := context.Background()

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

	ent, err := server.Svr.GetObject(ctx, setName, fp)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.ErrorCodeResponse[errors.ErrServer]
		}
		c.Status(err.HTTPStatusCode)
		return
	}

	c.Header("Full-Path", string(ent.FullPath))
	c.Header("Set", string(ent.Set))
	c.Header("Creation-Time", ent.Ctime.Format(timeFormat))
	c.Header("Mode", strconv.FormatUint(uint64(ent.Mode), 10))
	c.Header("Mime", ent.Mime)
	c.Header("Md5", util.Md5ToStr(ent.Md5))
	c.Header("File-Size", strconv.FormatUint(ent.FileSize, 10))
	c.Header("Fid", ent.Fid)
	c.Status(http.StatusOK)
	return
}
