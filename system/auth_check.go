package system

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtCustomClaims 自定义 JWT
type JwtCustomClaims struct {
	UID string `json:"uid"`
	*jwt.StandardClaims
}

// CreateAccessToken 生成 AccessToken timeout_secod 为超时秒数
func CreateAccessToken(uid, JWTSecretKey string, timeoutSecond int64) string {
	mySigningKey := []byte(JWTSecretKey)
	claims := JwtCustomClaims{
		uid,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(timeoutSecond)).Unix(),
		},
	}

	PrintInterface(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}

// CheckAccessTokenValid 判断 accessToken 与 uid 的授权是否有效，并重新刷新时间
func CheckAccessTokenValid(accessToken, JWTSecretKey string) (*JwtCustomClaims, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	}

	var claims JwtCustomClaims
	token, err := jwt.ParseWithClaims(accessToken, &claims, keyFunc)
	println(claims.UID)
	if err == nil {
		if token.Valid {
			return &claims, nil
		}
	}

	return nil, err
}

// GetClaims ...
func GetClaims(key interface{}) *JwtCustomClaims {
	token := key.(*jwt.Token)
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims
	}

	return nil
}
