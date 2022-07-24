package biz

import (
	"context"
	"encoding/base64"
)

type EncryptService interface {
	Encrypt(ctx context.Context, target []byte) (result []byte, err error)
}

type encryptServiceImpl struct {
}

func NewEncryptService() EncryptService {
	return &encryptServiceImpl{}
}

func (e *encryptServiceImpl) Encrypt(ctx context.Context, target []byte) (result []byte, err error) {
	encodeToString := base64.StdEncoding.EncodeToString(target)
	return []byte(encodeToString), nil
}
