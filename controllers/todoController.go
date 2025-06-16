package controllers

import (
	"net/http"
	"todolist/config"
	"todolist/models"

	"github.com/gin-gonic/gin"
)

// 获取所有的todo
func GetTodos(c *gin.Context) {
	var todos []models.Todo
	result := config.DB.Find(&todos)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": todos})
}

// 更新todo
func UpdateTodo(c *gin.Context) {
	var todo models.Todo
	id := c.Param("id")
	//查询任务是否存在
	if err := config.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	// 用JSON数据绑定要更新的字段
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	//更新字段
	todo.Title = input.Title
	todo.Description = input.Description
	todo.Status = input.Status

	if err := config.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todo})
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	var todo models.Todo
	//查询任务是否存在
	if err := config.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if err := config.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}
