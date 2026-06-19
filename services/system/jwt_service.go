package services

import (
	"errors"
	"south-admin-gin/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtClaims 自定义 JWT Claims
type JwtClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 签发 JWT
func GenerateToken(userID int64, username string) (string, error) {
	cfg := config.GetConfig()
	secret := cfg.JWT.Secret
	if secret == "" {
		secret = "default-secret"
	}

	expireHour := cfg.JWT.ExpireHour
	if expireHour <= 0 {
		expireHour = 72
	}

	claims := JwtClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "south-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并校验 JWT
func ParseToken(tokenStr string) (*JwtClaims, error) {
	cfg := config.GetConfig()
	secret := cfg.JWT.Secret
	if secret == "" {
		secret = "default-secret"
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的 token")
	}

	return claims, nil
}
