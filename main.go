package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
)

func main() {
	//配置GIN
	r := startRouter()
	//启动UDP，注意先开线程避免无限阻塞
	startUDPServer("127.0.0.1:8088")
	//启动GIN HTTP服务
	r.Run("127.0.0.1:8088")
}

func startRouter() *gin.Engine {
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

func startUDPServer(address string) {
	// 创建 服务器 UDP 地址结构。指定 IP + port
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("ResolveUDPAddr err:", err)
		return
	}
	// 监听 客户端连接
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("net.ListenUDP err:", err)
		return
	}

	go func() {
		defer conn.Close()
		for {
			handelUDP(conn)
		}
	}()
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

func handelUDP(conn *net.UDPConn) {
	buf := make([]byte, 1024)
	len, clientAddress, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := string(buf[:len])
	fmt.Println(msg)
	bytes := []byte(checkTicket())
	//var str =[]byte(" $F12345678F$")
	conn.WriteToUDP(bytes, clientAddress) // 简单回写数据给客户端
}

func responseError(err error, resp *http.Response) bool {
	if err != nil {
		fmt.Println(resp, err)
		return true
	}
	return false
}

//模拟发起HTTP请求
func checkTicket() string {
	clt := http.Client{}
	resp, err := clt.Get("https://dev.hao88.cloud/log/get")
	if responseError(err, resp) {
		return "check fail"
	}
	content, err := ioutil.ReadAll(resp.Body)
	respBody := string(content)
	fmt.Println(respBody, err)
	return "check ticket success"
}
