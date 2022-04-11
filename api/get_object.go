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
	if fp == "/" || fp[len(fp)-1] == '/' {
		//redirect to list object
		ListObjectHandler(c)
		return
	}

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
		ListObjectHandler(c)
		return
	}

	volumeId, _, err := storage_engine.ParseFid(ent.Fid)
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
	}

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
