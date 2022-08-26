package api

import (
	"github.com/gin-gonic/gin"
	"icesfs/command/vars"
	"icesfs/errors"
	"icesfs/full_path"
	"icesfs/server"
	"icesfs/set"
	"net/http"
)

func RecoverObjectHandler(c *gin.Context) {
	setName, fp := set.Set(c.Param("set")), full_path.FullPath(c.Param("fp"))
	if !setName.IsLegal() {
		err := errors.GetAPIErr(errors.ErrIllegalSetName)
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
	if !fp.IsLegalObjectName() {
		err := errors.GetAPIErr(errors.ErrIllegalObjectName)
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

	err := server.RecoverObject(c, setName, fp)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.GetAPIErr(errors.ErrServer)
		}
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

	c.Status(http.StatusAccepted)
	return
}
