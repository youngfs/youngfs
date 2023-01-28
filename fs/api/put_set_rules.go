package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/server"
	fs_set "youngfs/fs/set"
	"youngfs/log"
	"youngfs/vars"
)

func PutSetRulesHandler(c *gin.Context) {
	set := fs_set.Set(c.Param("set"))
	if len(set) < 1 { // include /*set
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
	set = set[1:]
	if set != "" && !set.IsLegal() {
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

	setRules := &fs_set.SetRules{}
	err := c.Bind(setRules)
	if err != nil {
		apiErr := errors.ErrRouter
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

	if setRules.Set == "" || set == "" {
		if set != "" {
			setRules.Set = set
		} else {
			set = setRules.Set
		}

		if set == "" {
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
	}

	if setRules.Set != set {
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

	setRules.MAXShardSize *= 1024 * 1024 // change MiB

	err = server.InsertSetRules(c, setRules)
	if err != nil {
		apiErr := &errors.APIError{}
		if !errors.As(err, &apiErr) {
			log.Errorw("a non api error is returned", vars.ErrorKey, err.Error())
			apiErr = errors.ErrNonApiErr
		}
		if apiErr.IsServerErr() {
			log.Errorf("uuid:%s\n error:%+v\n", c.Value(vars.UUIDKey), apiErr)
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

	c.Status(http.StatusCreated)
}
