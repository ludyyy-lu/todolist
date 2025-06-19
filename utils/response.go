package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func Success(c *gin.Context, data any, msg string) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
		"data": data,
	})
}

func Error(c *gin.Context,code int, err string) {
	c.JSON(code, gin.H {
		"code": code,
		"msg": err,
	})
}