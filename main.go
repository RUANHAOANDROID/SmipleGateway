package main

func main() {
	//配置GIN
	r := StartRouter()
	//启动UDP，注意先开线程避免无限阻塞
	StartUDPServer("127.0.0.1:8088")
	//启动GIN HTTP服务
	r.Run("127.0.0.1:8088")
}
