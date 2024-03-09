package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/server"
	"github.com/youngfs/youngfs/pkg/log"
)

type Handler struct {
	logger log.Logger
	svr    *server.Server
}

func New(logger log.Logger, svr *server.Server) *Handler {
	return &Handler{
		logger: logger,
		svr:    svr,
	}
}

func (h *Handler) errorHandler(c *gin.Context, err error) {
	e := &errors.Error{}
	if !errors.As(err, &e) {
		h.logger.Errorf("uuid:%s\n error:%v\n", c.Value(UUIDKey), err)
		e = errors.ErrNonApiErr
	}
	if e.IsServerErr() {
		h.logger.Errorf("uuid:%s\n error:%v\n", c.Value(UUIDKey), err)
	}
	c.Set(CodeKey, e.Code)
	c.Set(ErrorKey, e.Error())
	c.JSON(
		e.HTTPStatusCode,
		gin.H{
			UUIDKey:  c.Value(UUIDKey),
			CodeKey:  e.Code,
			ErrorKey: e.Error(),
		},
	)
}
