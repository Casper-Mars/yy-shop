package biz

import (
	"context"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"yy-shop/internal/conf"
)

type EncryptService interface {
	Encrypt(ctx context.Context, target []byte) (result []byte, err error)
	// Token 签发token
	Token(ctx context.Context, user *User) (string, error)
}

type encryptServiceImpl struct {
	authConfig *conf.Auth
}

func NewEncryptService(authConfig *conf.Bootstrap) EncryptService {
	return &encryptServiceImpl{
		authConfig: authConfig.Auth,
	}
}

func (e *encryptServiceImpl) Encrypt(ctx context.Context, target []byte) (result []byte, err error) {
	encodeToString := base64.StdEncoding.EncodeToString(target)
	return []byte(encodeToString), nil
}

func (e *encryptServiceImpl) Token(ctx context.Context, user *User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(e.authConfig.GetExpireDuration().AsDuration())), // 设置token的过期时间
	})
	return claims.SignedString([]byte(e.authConfig.GetJwtSecret()))
}
