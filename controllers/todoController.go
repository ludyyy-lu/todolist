package controllers

import (
	"net/http"
	"strconv"
	"todolist/config"
	"todolist/models"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		utils.Error(
			c,
			http.StatusBadRequest,
			"获取失败",
		)
		return
	}

	userIDVal, exists := c.Get("user_id") // 从上下文中获取用户ID
	if !exists {
		utils.Error(
			c,
			http.StatusUnauthorized,
			"未认证用户",
		)
		return
	}
	userID := userIDVal.(uint) // 转换为uint类型

	todo.UserID = userID

	if todo.Status == "" {
		todo.Status = "pending" //默认状态
	}
	if err := config.DB.Create(&todo).Error; err != nil {
		utils.Error(
			c,
			http.StatusInternalServerError,
			"创建失败",
		)
		return
	}

	utils.Success(
		c,
		gin.H{
			"todo":    todo,
		},
		"创建成功",
	)
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
	var total int64
	query := config.DB.Where("user_id = ?", userID)
	query.Count(&total)

	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	//执行查询语句
	err := query.Order(orderStr).Limit(size).Offset(offset).Find(&todos).Error
	if err != nil {
		utils.Error(
			c,
			http.StatusInternalServerError,
			"获取失败",
		)
		return
	}

	responseData := gin.H{
		"page":    page,
		"size":    size,
		"data":    todos,
		"totals":  total,
	}

	utils.Success(
		c,
		responseData,
		"获取成功",
	)
}

func getTodoByID(c *gin.Context) (*models.Todo, bool) {
	//todo的id
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "任务不存在或无权限")
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
		utils.Error(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	//更新字段
	todo.Title = input.Title
	todo.Description = input.Description
	todo.Status = input.Status

	if err := config.DB.Save(&todo).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.Success(c, gin.H{"data": todo}, "更新成功")
}

// 因为是Gorm.Model，所以会自动生成ID、CreatedAt、UpdatedAt等字段
// 所以这是软删除
func DeleteTodo(c *gin.Context) {
	todo, ok := getTodoByID(c)
	if !ok {
		return
	}
	if err := config.DB.Delete(&todo).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}
	utils.Success(c, nil, "删除成功")
}

func GetTodoStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")

	var result []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	// 分组统计每种状态
	if err := config.DB.Model(&models.Todo{}).	//指定表名是todos
	    Where("user_id = ?", userID).			//只查询当前用户的任务
		Select("status, COUNT(*) as count").	//选择状态和计数
		Group("status").						//按状态分组
		Scan(&result).Error; err != nil {		//把每行查到的数据填到result结构体中
			utils.Error(c, http.StatusInternalServerError, "统计失败")
			return
		}

		//把结果转换成 map 数据结构
		stats := map[string] int64 {
			"pending": 0,
			"in_progress": 0,
			"done": 0,
			"expired": 0,
		}
		for _, r := range result {
			stats[r.Status] = r.Count
		}

		utils.Success(c, gin.H{"statistics": stats}, "统计成功")
}

func RecoverTodo(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")

	var todo models.Todo
	//查找软删除的数据，需要用 Unscoped
	if err := config.DB.Unscoped().Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "任务不存在")
		return
	}

	// 只有软删除的数据才允许恢复
	if !todo.DeletedAt.Valid {
		utils.Error(c, http.StatusBadRequest, "任务未被软删除, 无需恢复")
		return
	}

	// 通过更新deletedAt为null来恢复
	if err := config.DB.Model(&todo).Update("deleted_at",nil).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "恢复失败")
		return
	}
}
