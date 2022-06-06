package util

import (
	"douyin-proj/src/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte(config.SecretKey)

type Claims struct {
	UserId uint `json:"userId"`
	jwt.StandardClaims
}

// ReleaseToken 生成token
func ReleaseToken(id uint) (string, error) {
	expireTime := time.Now().Add(1 * time.Hour) // token过期时间
	claims := &Claims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "douyin",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// VerifyToken 校验token
func VerifyToken(tokenString string) (uint, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	// token 合法并且该用户在数据库中存在
	if token != nil && token.Valid {
		return Claims.UserId, nil
	}
	return 0, err
}
