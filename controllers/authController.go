package controllers

import (
	"net/http"
	"todolist/config"
	"todolist/models"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input models.User

	// 检查请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查数据库 查看用户名是否已经存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		//查到用户信息，说明已经存在
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
	}
	//加密
	hashed, _ := utils.HashPassword(input.Password)
	input.Password = hashed

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败，用户名可能重复"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未注册"})
		return
	}

	if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
