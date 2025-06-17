package controllers

import (
	"net/http"
	"strconv"
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

	//解析分页参数
	//从 URL 参数里拿出用户传的分页参数
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	// string -> int
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size

	//模糊搜索
	keyword := c.Query("keyword")
	sort := c.DefaultQuery("sort", "created_at_desc") // 默认按创建时间降序排序

	orderStr := "created_at DESC" // 默认排序方式
	switch sort {
	case "created_at_asc":
		orderStr = "created_at ASC"
	case "updated_at_desc":
		orderStr = "updated_at DESC"
	case "updated_at_asc":
		orderStr = "updated_at ASC"
	}

	var todos []models.Todo
	query := config.DB.Where("user_id = ?", userID)
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	//执行查询语句
	err := query.Order(orderStr).Limit(size).Offset(offset).Find(&todos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"size":    size,
		"data":    todos,
		"totals":  len(todos),
		"message": "获取成功",
	})
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
