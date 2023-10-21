package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/fs/api"
	"github.com/youngfs/youngfs/pkg/fs/ui"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/vars"
	"io"
	"io/fs"
	"net/http"
)

var router *gin.Engine

func InitRouter() {
	if !vars.Debug {
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
	r.PUT("/:bucket/*path", Logger("put object"), api.PutObjectHandler)
	r.POST("/:bucket/*path", Logger("put object"), api.PutObjectHandler)
	r.GET("/:bucket/*path", Logger("get object"), api.GetObjectHandler)
	r.DELETE("/:bucket/*path", Logger("delete object"), api.DeleteObjectHandler)
	r.HEAD("/:bucket/*path", Logger("head object"), api.HeadObjectHandler)

	router = r
}

func Run() {
	err := router.Run(":" + vars.Port)
	if err != nil {
		log.Errorw("gin router init failed", "error", err.Error())
	}
}
