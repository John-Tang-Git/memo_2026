package common

import (
	"fmt"
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

type UserInfo struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Password string
}

// 生成token
func ReleaseToken(user UserInfo) (string, error) {
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

	fmt.Printf("签发时间: %s\n", time.Unix(claim.IssuedAt, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("过期时间: %s\n", time.Unix(claim.ExpiresAt, 0).Format("2006-01-02 15:04:05"))

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
