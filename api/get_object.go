package api

import (
	"github.com/gin-gonic/gin"
	"icesos/errors"
	"icesos/full_path"
	"icesos/server"
	"icesos/set"
	"icesos/storage_engine"
	"net/http"
)

func GetObjectHandler(c *gin.Context) {
	setName, fp := set.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
	if !setName.IsLegal() {
		err := errors.ErrorCodeResponse[errors.ErrIllegalSetName]
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				"code":  err.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}
	if !fp.IsLegal() {
		err := errors.ErrorCodeResponse[errors.ErrIllegalObjectName]
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				"code":  err.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}
	fp = fp.Clean()

	ent, err := server.Svr.GetObject(c, setName, fp)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.ErrorCodeResponse[errors.ErrServer]
		}
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				"code":  err.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}

	if ent.IsDirectory() {
		//redirect to list object
		accepts := c.Request.Header.Values("Accept")
		for _, str := range accepts {
			if str == "application/json" {
				ListObjectHandler(c)
				return
			}
		}

		err := errors.ErrorCodeResponse[errors.ErrInvalidPath]
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				"code":  err.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}

	volumeId, _ := storage_engine.SplitFid(ent.Fid)
	url, err := server.Svr.StorageEngine.GetVolumeIp(c, volumeId)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.ErrorCodeResponse[errors.ErrServer]
		}
		c.JSON(
			err.HTTPStatusCode,
			gin.H{
				"code":  err.ErrorCode,
				"error": err.Error(),
			},
		)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "http://"+url+"/"+ent.Fid)
	return
}
