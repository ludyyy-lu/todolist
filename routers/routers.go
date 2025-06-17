package routers

import (
	"todolist/controllers"
	"todolist/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouters(r *gin.Engine) {
	// //获取所有todo
	// r.GET("/todos", controllers.GetTodos)
	// //创建todo
	// r.POST("/todos", controllers.CreateTodo)

	// //更新todo
	// r.PUT("/todos/:id", controllers.UpdateTodo)

	// //删除todo
	// r.DELETE("/todos/:id", controllers.DeleteTodo)

	//注册
	r.POST("/register", controllers.Register)

	//登录
	r.POST("/login", controllers.Login)

	// 需要鉴权的路由组
	auth := r.Group("/")
	auth.Use(middlewares.JWTAuth())
	{
		auth.POST("/todos", controllers.CreateTodo)
		auth.PUT("/todos/:id", controllers.UpdateTodo)
		auth.DELETE("/todos/:id", controllers.DeleteTodo)
		auth.GET("/todos", controllers.GetTodos)

		auth.GET("/tags", controllers.GetTags)
		auth.POST("/tags", controllers.CreateTag)
		auth.DELETE("/tags/:id", controllers.DeleteTag)
		
		auth.POST("/todos/:id/tags", controllers.SetTodoTags)
		auth.GET("/todos/:id/tags", controllers.GetTodoTags)
		auth.DELETE("/todos/:id/tags/:tag_id", controllers.RemoveTodoTag)
	}
}
