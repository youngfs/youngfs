package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/youngfs/youngfs/pkg/fs/handler"
	"github.com/youngfs/youngfs/pkg/log"
)

func Logger(logger log.Logger, reqs ...string) gin.HandlerFunc {
	req := ""
	for _, u := range reqs {
		req += u
	}

	return func(c *gin.Context) {
		c.Set(handler.UUIDKey, uuid.New().String())
		logger.Infow("request start",
			handler.UUIDKey, c.Value(handler.UUIDKey),
			requestNameKey, req,
			requestUrlKey, c.Request.URL,
		)
		c.Next()
		logger.Infow("request finish",
			handler.UUIDKey, c.Value(handler.UUIDKey),
			requestNameKey, req,
			responseHTTPCodeKey, c.Writer.Status(),
			handler.CodeKey, c.Value(handler.CodeKey),
			handler.ErrorKey, c.Value(handler.ErrorKey),
		)
	}
}
