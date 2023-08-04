package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = []byte("secret_key")

type Claims struct {
	UserID uint `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(id uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func AuthenticateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// 解析Token错误
	if err != nil {
		return 0, errors.New("invalid token")
	}

	// 验证Token是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	return 0, errors.New("invalid token")
}
