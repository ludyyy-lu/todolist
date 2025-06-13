package main

import (
	"todolist/config"
	"todolist/models"
	"todolist/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()
	// 自动迁移数据库
	models.AutoMigrate()
	// 初始化路由
	r := gin.Default()
	// 注册路由
	routers.SetupRouters(r)

	// 启动服务器
	r.Run(":8080")
}