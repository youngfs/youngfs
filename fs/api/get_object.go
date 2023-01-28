package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/full_path"
	"youngfs/fs/server"
	fs_set "youngfs/fs/set"
	"youngfs/log"
	"youngfs/vars"
)

func GetObjectHandler(c *gin.Context) {
	set, fp := fs_set.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
	if fp == "/" || fp[len(fp)-1] == '/' {
		//redirect to list object
		ListObjectHandler(c)
		return
	}

	if !set.IsLegal() {
		err := errors.ErrIllegalSetName
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
	if !fp.IsLegal() {
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

	ent, err := server.GetObject(c, set, fp)
	if err != nil {
		apiErr := &errors.APIError{}
		if !errors.As(err, &apiErr) {
			log.Errorw("a non api error is returned", vars.ErrorKey, err.Error())
			apiErr = errors.ErrNonApiErr
		}
		if apiErr.IsServerErr() {
			log.Errorf("uuid:%s\n error:%+v\n", c.Value(vars.UUIDKey), apiErr)
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

	if ent.IsDirectory() {
		//redirect to list object
		ListObjectHandler(c)
		return
	}

	url, err := server.GetFidUrl(c, ent.Fid)
	if err != nil {
		apiErr := &errors.APIError{}
		if !errors.As(err, &apiErr) {
			log.Errorw("a non api error is returned", vars.ErrorKey, err.Error())
			apiErr = errors.ErrNonApiErr
		}
		if apiErr.IsServerErr() {
			log.Errorf("uuid:%s\n error:%+v\n", c.Value(vars.UUIDKey), apiErr)
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

	c.Redirect(http.StatusFound, url)
	// http.StatusMovedPermanently 301: The URL of the requested resource has been changed permanently. The new URL is given in the response.
	// http.StatusFound            302: This response code means that the URI of requested resource has been changed temporarily. Further changes in the URI might be made in the future. Therefore, this same URI should be used by the client in future requests.
	return
}
