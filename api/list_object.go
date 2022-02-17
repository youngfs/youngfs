package api

import (
	"github.com/gin-gonic/gin"
	"icesos/directory"
	"icesos/errors"
	"icesos/full_path"
	"icesos/iam"
	"net/http"
	"sort"
)

//type ListObjectInfo struct {
//	User       string `form:"user" json:"user" uri:"user" binding:"required"`
//	SecretKey  string `form:"secretKey" json:"secretKey" uri:"secretKey" binding:"required"`
//	Set        string `form:"set" json:"set" uri:"set" binding:"required"`
//	ObjectName string `form:"objectName" json:"objectName" uri:"objectName" binding:"required"`
//}

func ListObjectHandler(c *gin.Context) {
	//listObjectInfo := &ListObjectInfo{}
	//
	//err := c.Bind(listObjectInfo)
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
	//set := iam.Set(listObjectInfo.Set)
	//fp := full_path.FullPath(listObjectInfo.ObjectName)
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
	//user := iam.User(listObjectInfo.User)
	//if !user.Identify(listObjectInfo.SecretKey) {
	//	c.JSON(
	//		http.StatusBadRequest,
	//		gin.H{
	//			"error": errors.ErrorCodeResponse[errors.ErrUserAuthenticate].Error(),
	//		},
	//	)
	//	return
	//}
	//
	//if !user.ReadSetPermission(set) {
	//	c.JSON(
	//		http.StatusBadRequest,
	//		gin.H{
	//			"error": errors.ErrorCodeResponse[errors.ErrSetReadAuthenticate].Error(),
	//		},
	//	)
	//	return
	//}

	set, fp := iam.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
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

	inodes, err := directory.GetInodes(set, fp)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	sort.Sort(directory.InodeSlice(inodes))
	c.JSON(
		http.StatusOK,
		gin.H{
			"Path":    fp,
			"Entries": inodes,
		},
	)

	return
}
