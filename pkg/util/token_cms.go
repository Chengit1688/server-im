package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CmsClaims struct {
	UserID  string
	RoleKey string
	jwt.StandardClaims
}

func CmsCreateToken(userID string, roleKey string, duration time.Duration) (tokenString string, err error) {
	expireTime := time.Now().Add(duration)
	claims := &CmsClaims{
		UserID:  userID,
		RoleKey: roleKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func CmsParseToken(tokenString string) (userID, roleKey string, expiresAt int64, err error) {
	var (
		claims CmsClaims
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
	roleKey = claims.RoleKey
	expiresAt = claims.StandardClaims.ExpiresAt
	return
}
