package biz

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLoginRequest(t *testing.T) {
	data := []struct {
		name     string
		username string
		password string
		wantErr  error
		wantData *LoginRequest
	}{
		{
			name:     "缺少用户名",
			password: "123456",
			wantErr:  ErrMissingUsername,
		},
		{
			name:     "缺少密码",
			username: "admin",
			wantErr:  ErrMissingPassword,
		},
		{
			name:     "正常",
			username: "admin",
			password: "123456",
			wantData: &LoginRequest{
				username: "admin",
				password: "123456",
			},
			wantErr: nil,
		},
	}
	for _, item := range data {
		t.Run(item.name, func(t *testing.T) {
			got, err := NewLoginRequest(item.username, item.password)
			assert.Equal(t, item.wantErr, err)
			if item.wantErr == nil {
				assert.Equal(t, item.wantData.username, got.username)
				assert.Equal(t, item.wantData.password, got.password)
			}
		})
	}
}
