package utils

import (
	"github.com/dgrijalva/jwt-go"
	"mxshop_web_api/middlewares"
	"mxshop_web_api/models"
	"time"
)

func GenerateToken(Id uint, NickName string, Role uint) (string, error) {
	j := middlewares.JWT{}
	claims := models.CustomClaims{
		ID:          Id,
		NickName:    NickName,
		AuthorityId: Role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 设置30天过期
			Issuer:    "pluto",
		},
	}
	token, err := j.CreateToken(claims)
	return token, err
}
