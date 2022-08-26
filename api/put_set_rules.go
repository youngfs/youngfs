package api

import (
	"github.com/gin-gonic/gin"
	"icesfs/command/vars"
	"icesfs/errors"
	"icesfs/server"
	"icesfs/set"
	"net/http"
)

func PutSetRulesHandler(c *gin.Context) {
	setName := set.Set(c.Param("set"))
	if len(setName) < 1 { // include /*set
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
	if setName != "" && !setName.IsLegal() {
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

	setRules := &set.SetRules{}
	err := c.Bind(setRules)
	if err != nil {
		apiErr := errors.GetAPIErr(errors.ErrRouter)
		c.Set(vars.CodeKey, apiErr.ErrorCode)
		c.Set(vars.ErrorKey, err.Error())
		c.JSON(
			apiErr.HTTPStatusCode,
			gin.H{
				vars.UUIDKey:  c.Value(vars.UUIDKey),
				vars.CodeKey:  apiErr.ErrorCode,
				vars.ErrorKey: err.Error(),
			},
		)
		return
	}

	if setRules.Set == "" || setName == "" {
		if setName != "" {
			setRules.Set = setName
		} else {
			setName = setRules.Set
		}

		if setName == "" {
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
	}

	if setRules.Set != setName {
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

	setRules.MAXShardSize *= 1024 * 1024 // change MiB

	err = server.InsertSetRules(c, setRules)
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

	c.Status(http.StatusCreated)
}
