package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/fs/handler"
	"github.com/youngfs/youngfs/pkg/fs/ui"
	"io"
	"io/fs"
	"net/http"
)

func New(h *handler.Handler, options ...Option) http.Handler {
	cfg := &config{}
	for _, opt := range options {
		opt.apply(cfg)
	}

	if !cfg.debug {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	r := gin.Default()
	r.MaxMultipartMemory = 1 << 30
	//html template
	r.SetFuncMap(ui.FuncMap)
	r.SetHTMLTemplate(ui.StatusTpl)

	//static
	staticFS, _ := fs.Sub(ui.Static, "static")
	r.StaticFS("/static", http.FS(staticFS))

	//favicon.ico
	faviconHandler := func(c *gin.Context) {
		c.Data(
			http.StatusOK,
			"image/x-icon",
			ui.Favicon,
		)
	}
	r.GET("/favicon.ico", faviconHandler)
	r.HEAD("/favicon.ico", faviconHandler)

	// api handler
	// object
	r.Use(cfg.middlewares...)
	r.PUT(fmt.Sprintf("/:%s/*%s", handler.QueryBucketKey, handler.QueryPathKey), h.PutObjectHandler)
	r.POST(fmt.Sprintf("/:%s/*%s", handler.QueryBucketKey, handler.QueryPathKey), h.PutObjectHandler)
	r.GET(fmt.Sprintf("/:%s/*%s", handler.QueryBucketKey, handler.QueryPathKey), h.GetObjectHandler)
	r.DELETE(fmt.Sprintf("/:%s/*%s", handler.QueryBucketKey, handler.QueryPathKey), h.DeleteObjectHandler)
	r.HEAD(fmt.Sprintf("/:%s/*%s", handler.QueryBucketKey, handler.QueryPathKey), h.HeadObjectHandler)

	return r
}
