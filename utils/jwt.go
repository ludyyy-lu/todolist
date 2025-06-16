package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key") // 替换为你的密钥，签名时使用

// 生成token
func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID, // 存储用户ID
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 过期时间为24小时
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)  // 创建token
	return token.SignedString(jwtKey) // 给token加签名
}

// 解析token，从token中提取用户ID
func ParseToken(tokenStr string) (uint, error) {
	// 第二个参数是一个回调函数，用于验证token的签名
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token)(any, error){
		return jwtKey, nil // 使用相同的密钥解析token
	})
	//第一行 把token中的claims强制转换成jwt.MapClaims类型（键值对），然后赋值给claims变量。
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//因为 MapClaims 是 map[string]interface{}，而 JSON 的数字默认是 float64 类型，所以你要先转成 float64，再转成 uint
		id := uint(claims["user_id"].(float64)) // 从claims中获取用户ID
		return id, nil
	}
	return 0, err
}