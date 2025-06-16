package routers

import(
	"github.com/gin-gonic/gin"
	"todolist/controllers"
)

func SetupRouters(r *gin.Engine) {
	//获取所有todo
	r.GET("/todos",controllers.GetTodos)
	//创建todo
	r.POST("/todos",controllers.CreateTodo)

	//更新todo
	r.PUT("/todos/:id",controllers.UpdateTodo)

	//删除todo
	r.DELETE("/todos/:id",controllers.DeleteTodo)
}