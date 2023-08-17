package router

import (
	"GFV/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(e *gin.Engine) {
	oc := controller.OnlineController{}
	e.GET("/onlineConvert", oc.OnlineConvert) //转换通过url获取的文件
	e.POST("/localConvert", oc.LocalConvert)  //转换本地文件
	e.GET("/js/lazyload.js", oc.Static)
	e.GET("/js/pdfobject.js", oc.Static)
	e.GET("/images/loading.gif", oc.Static)
	e.StaticFS("/office_asset", http.Dir("cache/convert"))
	e.StaticFS("/pdf_asset", http.Dir("cache/pdf"))
	e.GET("/preview", oc.Preview) //预览转换后的文件
}
