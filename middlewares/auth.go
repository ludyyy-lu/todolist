package middlewares

import (
	"net/http"
	"strings"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败"})
			c.Abort() // 终止后续处理
			return
		}
		token := strings.TrimPrefix(auth, "Bearer")
		userID, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 无效"})
			c.Abort() // 终止后续处理
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
