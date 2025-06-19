package controllers

import (
	"net/http"
	"todolist/config"
	"todolist/models"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

type TagInput struct {
	Name string `json:"name" binding:"required"` // 标签名称，必填项
}

// POST/tags
func CreateTag(c *gin.Context) {
	var input TagInput
	if err := c.ShouldBindJSON(&input); err != nil { // 绑定JSON数据到input结构体中
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	userID := c.GetUint("user_id")
	tag := models.Tag{Name: input.Name, UserID: userID} // 创建Tag模型
	//写入数据库
	if err := config.DB.Create(&tag).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}
	utils.Success(
		c,
		tag, // 返回创建的Tag信息
		"创建成功", // 返回成功信息
	)
}

// GET/tags 获取所有标签
func GetTags(c *gin.Context) {
	userID := c.GetUint("user_id")
	var tags []models.Tag
	if err := config.DB.Where("user_id = ?", userID).Find(&tags).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}
	utils.Success(
		c,
		tags, // 返回标签列表
		"查询成功", // 返回成功信息
	)
}

func DeleteTag(c *gin.Context) {
	tagID := c.Param("id")
	userID := c.GetUint("user_id")

	var tag models.Tag
	if err := config.DB.First(&tag, tagID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "标签不存在")
		return
	}
	if tag.UserID != userID { // 判断标签是否属于当前用户
		utils.Error(c, http.StatusForbidden, "权限不足")
		return
	}
	config.DB.Model(&tag).Association("Todos").Clear() // 清除标签与任务的关联关系
	if err := config.DB.Delete(&tag).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.Success(
		c,
		nil,
		"删除成功", // 返回成功信息
	)
}

// GET/todos/:id/tags 获取某个任务的所有标签
func GetTodoTags(c *gin.Context) {
	todoID := c.Param("id")        // 获取Todo的ID
	userID := c.GetUint("user_id") // 获取当前用户的ID

	var todo models.Todo // 创建Todo模型
	if err := config.DB.Preload("Tags").Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {

		utils.Error(c, http.StatusNotFound, "任务不存在")
		return
	}

	utils.Success(
		c,
		gin.H{
			"tags": todo.Tags, // 返回标签列表
		},
		"查询成功", // 返回成功信息
	)
}

func RemoveTodoTag(c *gin.Context) {
	todoID := c.Param("id")        // 获取Todo的ID
	tagID := c.Param("tag_id")     // 获取Tag的ID
	userID := c.GetUint("user_id") // 获取当前用户的ID

	var todo models.Todo                                                                                               // 创建Todo模型
	if err := config.DB.Preload("Tags").Where("id =? AND user_id =?", todoID, userID).First(&todo).Error; err != nil { // 查找Todo，包括已关联的Tag和未关联的Tag
		utils.Error(c, http.StatusNotFound, "任务不存在")
		return
	}
	var tag models.Tag                                                                               // 创建Tag模型
	if err := config.DB.Where("id =? AND user_id =?", tagID, userID).First(&tag).Error; err != nil { // 查找Tag

		utils.Error(c, http.StatusNotFound, "标签不存在")
		return
	}

	if err := config.DB.Model(&todo).Association("Tags").Delete(&tag); err != nil { // 删除关联关系
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.Success(
		c,
		nil,
		"删除成功", // 返回成功信息
	)
}

type TodoTagInput struct { // 关联Todo和Tag的结构体，用于创建关联关系
	TagIDs []uint `json:"tag_ids"` // Tag的ID，必填项
}

// 给某个任务设置标签
// POST/todos/:id/tags 关联Todo和Tag
func SetTodoTags(c *gin.Context) {
	var input TodoTagInput
	if err := c.ShouldBind(&input); err != nil {
		utils.Error(
			c,
			http.StatusBadRequest,
			"参数错误", // 返回错误信息
		)
		return
	}

	todoID := c.Param("id")
	userID := c.GetUint("user_id")

	var todo models.Todo
	if err := config.DB.Preload("Tags").Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		utils.Error(
			c,
			http.StatusNotFound,
			"任务不存在", // 返回错误信息
		)
		return
	}

	//查询所有TAg
	var tags []models.Tag
	if err := config.DB.Where("id IN ? AND user_id = ?", input.TagIDs, userID).Find(&tags).Error; err != nil { // 查找所有Tag，包括已关联的Tag和未关联的Tag
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "查询Tag失败"}) // 返回错误信息
		utils.Error(
			c,
			http.StatusInternalServerError,
			"查询Tag失败", // 返回错误信息
		)
		return
	}

	if err := config.DB.Model(&todo).Association("Tags").Replace(&tags); err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "关联Tag失败"}) // 返回错误信息
		utils.Error(
			c,
			http.StatusInternalServerError,
			"关联Tag失败") // 返回错误信息)
		return
	}
	utils.Success(
		c,
		tags, // 返回关联后的Todo信息
		"已设置标签",      // 返回成功信息
	)
}
