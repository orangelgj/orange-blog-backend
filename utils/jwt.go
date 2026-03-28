package utils

import (
	"errors"
	"gblog/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID             uint32 `json:"user_id"`
	Role               int8   `json:"role"`
	PasswordUpdateTime int64  `json:"password_update_time"`
	jwt.RegisteredClaims
}

// GenerateToken 生成一个 Token
func GenerateToken(userID uint32, role int8, passwordUpdateTime int64) (string, error) {
	var jwtKey = []byte(config.AppConfig.JWT.Secret)   // 实际开发中请放在环境变量
	expirationTime := time.Now().Add(2400 * time.Hour) // 有效期24小时
	claims := &Claims{
		UserID:             userID,
		Role:               role,
		PasswordUpdateTime: passwordUpdateTime,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   "user_auth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (Claims, error) {
	var jwtKey = []byte(config.AppConfig.JWT.Secret)

	// 1. 调用 jwt.ParseWithClaims
	// 第三个参数是一个回调函数，用来返回你签名时用的密钥 (Key)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return Claims{}, err
	}

	// 2. 校验并断言
	// token.Valid 会自动检查 ExpiresAt (过期时间) 是否合法
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return *claims, nil
	}

	return Claims{}, errors.New("invalid token")
}
