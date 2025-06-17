package controllers

import (
	"net/http"
	"todolist/config"
	"todolist/models"

	"github.com/gin-gonic/gin"
)

type TagCreateInput struct {
	Name string `json:"name" binding:"required"` // 标签名称，必填项
}
//POST/tags
func CreateTag(c *gin.Context) {
	var input TagCreateInput
	if err := c.ShouldBindJSON(&input); err != nil { // 绑定JSON数据到input结构体中
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"}) // 返回错误信息
		return
	}
	userID := c.GetUint("user_id")
	tag := models.Tag{Name: input.Name, UserID: userID} // 创建Tag模型
	if err := config.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	//这是什么意思
	c.JSON(http.StatusOK, tag)
}

type TodoTagInput struct { // 关联Todo和Tag的结构体，用于创建关联关系
	TagIDs  uint `json:"tag_ids"`  // Tag的ID，必填项
}

//POST/todos/:id/tags 关联Todo和Tag
func SetTodoTags(c *gin.Context) {
	var input TodoTagInput
	if err := c.ShouldBind(&input);err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "参数错误"})
		return
	}

	todoID := c.Param("id")
	var todo models.Todo
	if err := config.DB.Preload("Tags").Where("id = ?", todoID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	//查询所有TAg
	var tags []models.Tag
	if err := config.DB.Where("id IN ?",input.TagIDs).Find(&tags).Error; err!= nil { // 查找所有Tag，包括已关联的Tag和未关联的Tag
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询Tag失败"}) // 返回错误信息
		return
	}

	//将Tag添加到Todo的Tags字段中
	config.DB.Model(&todo).Association("Tags").Replace(&tags)
	c.JSON(
		http.StatusOK, gin.H{
			"message": "已设置标签", // 返回成功信息
			"tags":    tags,        // 返回关联后的Todo信息
		})
}