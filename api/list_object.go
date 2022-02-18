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

func ListObjectHandler(c *gin.Context) {
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
