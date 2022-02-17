package api

import (
	"crypto/md5"
	"github.com/gin-gonic/gin"
	"icesos/directory"
	"icesos/entry"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	"icesos/storage_engine"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//type PutObjectInfo struct {
//	User       string `form:"user" json:"user" uri:"user" binding:"required"`
//	SecretKey  string `form:"secretKey" json:"secretKey" uri:"secretKey" binding:"required"`
//}

func PutObjectHandler(c *gin.Context) {
	//putObjectInfo := &PutObjectInfo{}
	//
	//err := c.Bind(putObjectInfo)
	//if err != nil {
	//	c.JSON(
	//		http.StatusBadRequest,
	//		gin.H{
	//			"status": http.StatusBadRequest,
	//			"error":  err.Error(),
	//		},
	//	)
	//	return
	//}
	//
	//set := iam.Set(putObjectInfo.Set)
	//fp := full_path.FullPath(putObjectInfo.ObjectName)
	//if !fp.IsLegal() {
	//	c.JSON(
	//		http.StatusBadRequest,
	//		gin.H{
	//			"error": errors.ErrorCodeResponse[errors.ErrIllegalObjectName].Error(),
	//		},
	//	)
	//	return
	//}
	//fp = fp.Clean()
	//
	//user := iam.User(putObjectInfo.User)
	//if !user.Identify(putObjectInfo.SecretKey) {
	//	c.JSON(
	//		http.StatusBadRequest,
	//		gin.H{
	//			"error": errors.ErrorCodeResponse[errors.ErrUserAuthenticate].Error(),
	//		},
	//	)
	//	return
	//}
	//
	//if !user.WriteSetPermission(set) {
	//	c.JSON(
	//		http.StatusBadRequest,
	//		gin.H{
	//			"error": errors.ErrorCodeResponse[errors.ErrSetWriteAuthenticate].Error(),
	//		},
	//	)
	//	return
	//}'

	ctime := time.Unix(time.Now().Unix(), 0)

	file, head, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	set, fp := iam.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
	if len(fp) == 0 || fp[len(fp)-1] == '/' {
		fp += full_path.FullPath(head.Filename)
	}
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

	fid, err := storage_engine.PutObject(uint64(head.Size), file)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	volumeId, fid := storage_engine.SplitFid(fid)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	md5Ret := md5.Sum(b)

	err = file.Close()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	err = directory.InsertInode(
		&directory.Inode{
			FullPath: fp,
			Set:      set,
			Mtime:    ctime,
			Ctime:    ctime,
			Mode:     os.ModePerm,
			FileSize: uint64(head.Size),
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

	err = entry.InsertEntry(
		&entry.Entry{
			FullPath: fp,
			Set:      set,
			Mtime:    ctime,
			Ctime:    ctime,
			Mode:     os.ModePerm,
			Md5:      md5Ret,
			FileSize: uint64(head.Size),
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
