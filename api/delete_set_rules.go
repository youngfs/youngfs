package api

import (
	"github.com/gin-gonic/gin"
	"icesfs/command/vars"
	"icesfs/errors"
	"icesfs/server"
	"icesfs/set"
	"net/http"
)

func DeleteSetRulesHandler(c *gin.Context) {
	setName := set.Set(c.Param("set"))
	if len(setName) < 2 { // include /*set
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
	setName = setName[1:]
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

	err := server.DeleteSetRules(c, setName)
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
