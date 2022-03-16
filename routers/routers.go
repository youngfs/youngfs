package routers

import (
	"github.com/gin-gonic/gin"
	"icesos/api"
	"icesos/ui"
	"io/fs"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
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

	//api handler
	r.PUT("/:set/*fp", api.PutObjectHandler)
	r.POST("/:set/*fp", api.PutObjectHandler)
	r.GET("/:set/*fp", api.GetObjectHandler)
	r.DELETE("/:set/*fp", api.DeleteObjectHandler)
	r.HEAD("/:set/*fp", api.HeadObjectHandler)
	return r
}
