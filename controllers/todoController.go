package controllers

import (
	"net/http"
	"todolist/config"
	"todolist/models"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	var todos []models.Todo
	result := config.DB.Find(&todos)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{"data": todos})
}

