package common

import (
	"memo/controller"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwt加密密钥
var jwtKey = []byte("a_common_key")

// token的claim
type Claim struct {
	jwt.StandardClaims
	UserId uint
}

// 生成token
func ReleaseToken(user controller.UserInfo) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claim := Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",
			Subject:   "user token",
		},
		UserId: user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	// 错误处理
	if err != nil {
		return "", err
	}
	// 返回token
	return tokenString, err
}

// 从token中解析claims
func ParseToken(tokenString string) (*jwt.Token, *Claim, error) {
	claim := &Claim{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claim, err
}
