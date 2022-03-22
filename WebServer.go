package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func StartRouter() *gin.Engine {
	r := gin.Default()
	//把静态资源static/js目录挂在到相对路径js
	r.Static("/js", "./static/js")
	//挂载templates目录下html资源
	r.LoadHTMLGlob("templates/*.html")
	topGroup := r.Group("/page")   //页面组
	configGroup := r.Group("/api") //api组
	handlerHtml(topGroup)          //处理html组
	handlerConfig(configGroup)     //处理API配置组
	return r
}
func handlerHtml(r *gin.RouterGroup) {
	r.GET("/login", func(context *gin.Context) {
		context.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/config", func(context *gin.Context) {
		context.HTML(http.StatusOK, "config.html", nil)
	})
	r.POST("/saveConfig", func(context *gin.Context) {
		context.String(http.StatusOK, "{OK}")
	})
}

func handlerConfig(r *gin.RouterGroup) {

	r.GET("/config", func(context *gin.Context) {
		json, err := ioutil.ReadFile("./static/config.json")
		if err == nil {
			println(err)
		}
		fmt.Println(string(json))
		context.String(http.StatusOK, string(json))
	})
	r.POST("/config/save", func(context *gin.Context) {
		context.String(http.StatusOK, "保存成功")
	})
}

func responseError(err error, resp *http.Response) bool {
	if err != nil {
		fmt.Println(resp, err)
		return true
	}
	return false
}
