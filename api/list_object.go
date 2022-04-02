package api

import (
	"github.com/gin-gonic/gin"
	"icesos/errors"
	"icesos/full_path"
	"icesos/server"
	"icesos/set"
	"icesos/ui"
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
	if fp != "/" && fp[len(fp)-1] == '/' {
		fp = fp[:len(fp)-1]
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

	accepts := c.Request.Header.Values("Accept")
	for _, str := range accepts {
		if str == "application/json" {
			c.JSON(
				http.StatusOK,
				gin.H{
					"Path":    fp,
					"Set":     setName,
					"Entries": ents,
				},
			)
			return
		}
	}

	c.HTML(
		http.StatusOK,
		ui.UIName,
		gin.H{
			"FullPath":  string(fp),
			"Set":       string(setName),
			"PathLinks": fp.ToPathLink(),
			"Entries":   ents,
		},
	)
	return
}
