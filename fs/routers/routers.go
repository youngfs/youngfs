package routers

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"net/http"
	api2 "youngfs/fs/api"
	"youngfs/fs/ui"
	"youngfs/log"
	"youngfs/vars"
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

	//api handler
	// object
	r.PUT("/:set/*fp", Logger("put object"), api2.PutObjectHandler)
	r.POST("/:set/*fp", Logger("put object"), api2.PutObjectHandler)
	r.GET("/:set/*fp", Logger("get object"), api2.GetObjectHandler)
	r.DELETE("/:set/*fp", Logger("delete object"), api2.DeleteObjectHandler)
	r.HEAD("/:set/*fp", Logger("head object"), api2.HeadObjectHandler)

	// set rules
	r.PUT("/SetRules/*set", Logger("put set rules"), api2.PutSetRulesHandler)
	r.POST("/SetRules/*set", Logger("put set rules"), api2.PutSetRulesHandler)
	r.GET("/SetRules/*set", Logger("get set rules"), api2.GetSetRulesHandler)
	r.DELETE("/SetRules/*set", Logger("delete set rules"), api2.DeleteSetRulesHandler)
	r.HEAD("/SetRules/*set", Logger("head set rules"), api2.HeadSetRulesHandler)
	router = r
}

func Run() {
	err := router.Run(":" + vars.Port)
	if err != nil {
		log.Errorw("gin router init failed", "error", err.Error())
	}
}
