package api

import (
	"github.com/gin-gonic/gin"
	"icesfs/errors"
	"icesfs/server"
	"icesfs/set"
	"net/http"
)

func HeadSetRulesHandler(c *gin.Context) {
	setName := set.Set(c.Param("set"))
	if len(setName) < 2 { // include /*set
		err := errors.GetAPIErr(errors.ErrIllegalSetName)
		c.Status(err.HTTPStatusCode)
		return
	}
	setName = setName[1:]
	if !setName.IsLegal() {
		err := errors.GetAPIErr(errors.ErrIllegalSetName)
		c.Status(err.HTTPStatusCode)
		return
	}

	_, err := server.GetSetRules(c, setName)
	if err != nil {
		err, ok := err.(errors.APIError)
		if ok != true {
			err = errors.GetAPIErr(errors.ErrServer)
		}
		c.Status(err.HTTPStatusCode)
		return
	}

	c.Status(http.StatusOK)
	return
}
