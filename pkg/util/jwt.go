package util

import (
	"errors"
	"github.com/HuangQiang97/douyin-proj/pkg/constant"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtkey = []byte(constant.SecretKey)

type Claims struct {
	UserId  int64
	jwt.StandardClaims
}

func ReleaseToken(id int64) (string , error) {
	expireTime := time.Now().Add(24*time.Hour)
	claims := &Claims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "douyin",
			Subject: "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString , err := token.SignedString(jwtkey)
	if err != nil{
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (int64, error){
	Claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token)(i interface{}, err error){
		return jwtkey, nil
	})
	if err != nil{
		return 0, err
	}
	sub := time.Now().Unix() - Claims.ExpiresAt
	if sub > 0 {
		return 0, errors.New("the token is expired")
	}

	return Claims.UserId,nil
}