package api

import (
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/vars"
)

func errorHandler(c *gin.Context, err error) {
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
}
