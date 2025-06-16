package controllers

import (
	"net/http"
	"todolist/config"
	"todolist/models"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	userIDVal, exists := c.Get("user_id") // 从上下文中获取用户ID
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	userID := userIDVal.(uint) // 转换为uint类型

	todo.UserID = userID

	if todo.Status == "" {
		todo.Status = "pending" //默认状态
	}
	if err := config.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"todo":    todo,
	})
}

// 获取所有的todo
func GetTodos(c *gin.Context) {
	userID := c.GetUint("user_id") // 从上下文中获取用户ID
	var todos []models.Todo

	if err := config.DB.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": todos})
}

func getTodoByID(c *gin.Context) (*models.Todo, bool) {
	//todo的id
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var todo models.Todo
	if err := config.DB.First("id = ? AND user_id = ?", id, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限"})
		return nil, false
	}
	return &todo, true
}

// 更新todo
func UpdateTodo(c *gin.Context) {
	todo, ok := getTodoByID(c)
	if !ok {
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
	todo, ok := getTodoByID(c)
	if !ok {
		return
	}
	if err := config.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}
