package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/full_path"
	"youngfs/fs/server"
	fs_set "youngfs/fs/set"
	"youngfs/fs/ui"
	"youngfs/log"
	"youngfs/vars"
)

func ListObjectsHandler(c *gin.Context) {
	set, fp := fs_set.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
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
	if fp != "/" && fp[len(fp)-1] == '/' {
		fp = fp[:len(fp)-1]
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

	ents, err := server.ListObjects(c, set, fp)
	if err != nil {
		apiErr := &errors.APIError{}
		if !errors.As(err, &apiErr) {
			log.Errorw("a non api error is returned", vars.ErrorKey, err.Error())
			apiErr = errors.ErrNonApiErr
		}
		if apiErr.IsServerErr() {
			log.Errorf("uuid:%s\n error:%v\n", c.Value(vars.UUIDKey), err)
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

	accepts := c.Request.Header.Values("Accept")
	for _, str := range accepts {
		if str == "application/json" {
			c.JSON(
				http.StatusOK,
				gin.H{
					vars.UUIDKey: c.Value(vars.UUIDKey),
					"Path":       fp,
					"Set":        set,
					"Entries":    ents,
				},
			)
			return
		}
	}

	c.HTML(
		http.StatusOK,
		ui.FSName,
		gin.H{
			"FullPath":  string(fp),
			"Set":       string(set),
			"PathLinks": fp.ToPathLink(),
			"Entries":   ents,
		},
	)
	return
}
