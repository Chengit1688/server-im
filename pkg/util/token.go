package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte(`hello world`)

type Claims struct {
	UserID string
	jwt.StandardClaims
}

func CreateToken(userID string, duration time.Duration) (tokenString string, err error) {
	expireTime := time.Now().Add(duration)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ParseToken(tokenString string) (userID string, err error) {
	var (
		claims Claims
		token  *jwt.Token
	)
	if token, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	}); err != nil {
		return
	}

	if !token.Valid {
		err = errors.New("token invalid")
		return
	}
	userID = claims.UserID
	return
}
