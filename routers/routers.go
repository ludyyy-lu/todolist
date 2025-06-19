package routers

import (
	"todolist/controllers"
	"todolist/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouters(r *gin.Engine) {
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

		auth.GET("/todos/statistics", controllers.GetTodoStatistics)
		auth.PATCH("/todos/:id/recover", controllers.RecoverTodo)
		auth.GET("/todos/deleted", controllers.GetDeletedTodos)
	}
}
