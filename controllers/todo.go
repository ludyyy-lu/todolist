package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"todolist/models"
)

func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo);err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if todo.Status == "" {
		todo.Status = "pending" //默认状态
	}

	todo.UserID = 1 //暂时不写用户模块

	if err := models.CreateTodo(&todo);err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"todo": todo,
	})
}