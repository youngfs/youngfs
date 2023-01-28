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

func HeadSetRulesHandler(c *gin.Context) {
	set := fs_set.Set(c.Param("set"))
	if len(set) < 2 { // include /*set
		err := errors.ErrIllegalSetName
		c.Status(err.HTTPStatusCode)
		return
	}
	set = set[1:]
	if !set.IsLegal() {
		err := errors.ErrIllegalSetName
		c.Status(err.HTTPStatusCode)
		return
	}

	_, err := server.GetSetRules(c, set)
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

	c.Status(http.StatusOK)
	return
}
