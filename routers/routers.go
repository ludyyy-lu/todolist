package routers

import(
	"github.com/gin-gonic/gin"
	"todolist/controllers"
)

func SetupRouters(r *gin.Engine) {
	r.POST("/todos",controllers.CreateTodo)
}