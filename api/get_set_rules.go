package api

import (
	"github.com/gin-gonic/gin"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/server"
	"icesos/set"
	"icesos/ui"
	"net/http"
)

func GetSetRulesHandler(c *gin.Context) {
	setName := set.Set(c.Param("set"))
	if setName == "/" {
		hosts, err := server.GetHosts(c)
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

		c.HTML(
			http.StatusOK,
			ui.SetRulesName,
			gin.H{
				"Hosts": hosts,
			},
		)
		return
	}

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

	setRules, err := server.GetSetRules(c, setName)
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

	c.JSON(
		http.StatusOK,
		gin.H{
			vars.UUIDKey: c.Value(vars.UUIDKey),
			"Set":        setName,
			"SetRules":   setRules,
		},
	)
	return
}
