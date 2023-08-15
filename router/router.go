package router

import (
	"GFV/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(e *gin.Engine) {
	oc := controller.OnlineController{}
	e.GET("/onlinePreview", oc.OnlinePreview)
	e.GET("/js/lazyload.js", oc.Static)
	e.GET("/images/loading.gif", oc.Static)
	e.StaticFS("/office_asset", http.Dir("cache/convert"))
	e.GET("/offlinePreview", oc.OfflinePreview)
}
