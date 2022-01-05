package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	"icesos/storage_engine"
	"net/http"
	"strconv"
)

type GetObjectInfo struct {
	User       string `form:"user" json:"user" uri:"user" binding:"required"`
	SecretKey  string `form:"secretKey" json:"secretKey" uri:"secretKey" binding:"required"`
	Set        string `form:"set" json:"set" uri:"set" binding:"required"`
	ObjectName string `form:"objectName" json:"objectName" uri:"objectName" binding:"required"`
}

func GetObjectHandler(c *gin.Context) {
	getObjectInfo := &GetObjectInfo{}

	err := c.Bind(getObjectInfo)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	set := iam.Set(getObjectInfo.Set)
	fp := full_path.FullPath(getObjectInfo.ObjectName)
	if !fp.IsLegal() {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrIllegalObjectName].Error(),
			},
		)
		return
	}
	fp = fp.Clean()

	user := iam.User(getObjectInfo.User)
	if !user.Identify(getObjectInfo.SecretKey) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrUserAuthenticate].Error(),
			},
		)
		return
	}

	if !user.ReadSetPermission(set) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrSetReadAuthenticate].Error(),
			},
		)
		return
	}

	nowEntry, err := entry.GetEntry(set, fp)
	if err != nil {
		if err == redis.Nil {
			err = errors.ErrorCodeResponse[errors.ErrInvalidPath]
		}
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	url, err := storage_engine.GetVolumeIp(nowEntry.VolumeId)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "http://"+url+"/"+strconv.FormatUint(nowEntry.VolumeId, 10)+","+nowEntry.Fid)
	return
}
