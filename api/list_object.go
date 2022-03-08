package api

import (
	"github.com/gin-gonic/gin"
	"icesos/errors"
	"icesos/full_path"
	"icesos/server"
	"icesos/set"
	"net/http"
)

func ListObjectHandler(c *gin.Context) {
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

	ents, err := server.Svr.ListObejcts(c, setName, fp)
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

	c.JSON(
		http.StatusOK,
		gin.H{
			"Path":    fp,
			"Entries": ents,
		},
	)

	return
}
