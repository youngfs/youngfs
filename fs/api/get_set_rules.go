package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/server"
	fs_set "youngfs/fs/set"
	"youngfs/fs/ui"
	"youngfs/log"
	"youngfs/vars"
)

func GetSetRulesHandler(c *gin.Context) {
	set := fs_set.Set(c.Param("set"))
	if set == "/" {
		hosts, err := server.GetHosts(c)
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

		c.HTML(
			http.StatusOK,
			ui.SetRulesName,
			gin.H{
				"Hosts": hosts,
			},
		)
		return
	}

	if len(set) < 2 { // include /*set
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

	setRules, err := server.GetSetRules(c, set)
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

	c.JSON(
		http.StatusOK,
		gin.H{
			vars.UUIDKey: c.Value(vars.UUIDKey),
			"Set":        set,
			"SetRules":   setRules,
		},
	)
	return
}
