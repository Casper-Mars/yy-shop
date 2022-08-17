package biz

import (
	"context"
	"github.com/golang/mock/gomock"
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

func TestUser_CheckAuth(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	encryptService := NewMockEncryptService(controller)
	data := []struct {
		name     string
		initFunc func()
		wantErr  error
		password string
		user     User
	}{
		{
			name: "密码正确",
			initFunc: func() {
				encryptService.EXPECT().Encrypt(gomock.Any(), gomock.Any()).Return([]byte("123"), nil).Times(1)
			},
			wantErr:  nil,
			password: "123",
			user: User{
				Password: "123",
			},
		},
		{
			name: "密码错误",
			initFunc: func() {
				encryptService.EXPECT().Encrypt(gomock.Any(), gomock.Any()).Return([]byte("123"), nil).Times(1)
			},
			wantErr:  ErrPasswordWrong,
			password: "123",
			user: User{
				Password: "1233",
			},
		},
	}

	for _, item := range data {
		t.Run(item.name, func(t *testing.T) {
			item.initFunc()
			err := item.user.CheckAuth(context.Background(), item.password, encryptService)
			assert.Equal(t, item.wantErr, err)
		})
	}

}
