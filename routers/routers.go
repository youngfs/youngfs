package routers

import (
	"github.com/gin-gonic/gin"
	"icesos/api"
	"icesos/command/vars"
	"icesos/log"
	"icesos/ui"
	"io/fs"
	"io/ioutil"
	"net/http"
)

var router *gin.Engine

func InitRouter() {
	if !vars.Debug {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
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
	r.PUT("/:set/*fp", Logger("put object"), api.PutObjectHandler)
	r.POST("/:set/*fp", Logger("put object"), api.PutObjectHandler)
	r.GET("/:set/*fp", Logger("get object"), api.GetObjectHandler)
	r.DELETE("/:set/*fp", Logger("delete object"), api.DeleteObjectHandler)
	r.HEAD("/:set/*fp", Logger("head object"), api.HeadObjectHandler)

	// set rules
	r.PUT("/SetRules/*set", Logger("put set rules"), api.PutSetRulesHandler)
	r.POST("/SetRules/*set", Logger("put set rules"), api.PutSetRulesHandler)
	r.GET("/SetRules/*set", Logger("get set rules"), api.GetSetRulesHandler)
	r.DELETE("/SetRules/*set", Logger("delete set rules"), api.DeleteSetRulesHandler)
	r.HEAD("/SetRules/*set", Logger("head set rules"), api.HeadSetRulesHandler)
	router = r
}

func Run() {
	err := router.Run(":" + vars.Port)
	if err != nil {
		log.Errorw("gin router init failed", "error", err.Error())
	}
}
