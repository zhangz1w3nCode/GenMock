package main

import (
	"net/http"
)

func main() {
	// 注册 mock 接口
	http.HandleFunc("/mock", MockHandler) // 使用 server.MockHandler

	// 启动服务，监听 8080 端口
	http.ListenAndServe(":8080", nil)
}
