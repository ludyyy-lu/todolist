package main

import (
	"todolist/config"
	"todolist/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()
	// 初始化路由
	r := gin.Default()
	// 注册路由
	routers.SetRoutes(r)

	// 启动服务器
	r.Run(":8080")
}