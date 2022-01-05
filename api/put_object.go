package api

import (
	"github.com/gin-gonic/gin"
	"icesos/directory"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	"icesos/storage_engine"
	"net/http"
	"os"
	"time"
)

type PutObjectInfo struct {
	User       string `form:"user" json:"user" uri:"user" binding:"required"`
	SecretKey  string `form:"secretKey" json:"secretKey" uri:"secretKey" binding:"required"`
	Set        string `form:"set" json:"set" uri:"set" binding:"required"`
	ObjectName string `form:"objectName" json:"objectName" uri:"objectName" binding:"required"`
	FileSize   uint64 `form:"fileSize" json:"fileSize" uri:"fileSize" binding:"required"`
	Fid        string `form:"fid" json:"fid" uri:"fid" binding:"required"`
}

func PutObjectHandler(c *gin.Context) {
	putObjectInfo := &PutObjectInfo{}

	err := c.Bind(putObjectInfo)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status": http.StatusBadRequest,
				"error":  err.Error(),
			},
		)
		return
	}

	set := iam.Set(putObjectInfo.Set)
	fp := full_path.FullPath(putObjectInfo.ObjectName)
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

	user := iam.User(putObjectInfo.User)
	if !user.Identify(putObjectInfo.SecretKey) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrUserAuthenticate].Error(),
			},
		)
		return
	}

	if !user.WriteSetPermission(set) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": errors.ErrorCodeResponse[errors.ErrSetWriteAuthenticate].Error(),
			},
		)
		return
	}

	err = directory.InsertInode(
		&directory.Inode{
			FullPath: fp,
			Set:      set,
			Mtime:    time.Unix(time.Now().Unix(), 0),
			Ctime:    time.Unix(time.Now().Unix(), 0),
			Mode:     os.ModePerm,
			FileSize: putObjectInfo.FileSize,
		}, true)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	volumeId, fid := storage_engine.SplitFid(putObjectInfo.Fid)
	err = entry.InsertEntry(
		&entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    time.Unix(time.Now().Unix(), 0),
			Ctime:    time.Unix(time.Now().Unix(), 0),
			Mode:     os.ModePerm,
			FileSize: putObjectInfo.FileSize,
			VolumeId: volumeId,
			Fid:      fid,
		})
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.Status(http.StatusCreated)
	return
}
