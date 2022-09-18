package biz

import (
	"context"
	errors "errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"yy-shop/internal/conf"
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

func TestAccountUseCase_Login(t *testing.T) {

	controller := gomock.NewController(t)
	repo := NewMockUserRepo(controller)
	encryptService := NewMockEncryptService(controller)
	accountUseCase := NewAccountUseCase(log.DefaultLogger, &conf.Bootstrap{
		Auth: &conf.Auth{
			JwtSecret: "123",
		},
	}, repo, encryptService)

	data := []struct {
		name      string
		mockFunc  func()
		wantErr   assert.ErrorAssertionFunc
		wantToken string
		ctx       context.Context
		req       *LoginRequest
	}{
		{
			name: "正常登陆",
			mockFunc: func() {
				encryptService.EXPECT().Encrypt(gomock.Any(), gomock.Any()).Return([]byte("123"), nil).Times(1)
				repo.EXPECT().FetchByUsername(gomock.Any(), gomock.Any()).Return(&User{
					Password: "123",
				}, nil).Times(1)
				encryptService.EXPECT().Token(gomock.Any(), gomock.Any()).Return("123", nil).Times(1)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return true
			},
			wantToken: "123",
			ctx:       context.Background(),
			req: &LoginRequest{
				username: "123",
				password: "123",
			},
		},
		{
			name: "用户不存在",
			mockFunc: func() {
				repo.EXPECT().FetchByUsername(gomock.Any(), gomock.Any()).Return(nil, ErrUserNotExist).Times(1)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorAs(t, err, &ErrUserNotExist)
				return false
			},
			wantToken: "123",
			ctx:       context.Background(),
			req: &LoginRequest{
				username: "123",
				password: "123",
			},
		},
		{
			name: "密码校验不过",
			mockFunc: func() {
				encryptService.EXPECT().Encrypt(gomock.Any(), gomock.Any()).Return([]byte("1233"), nil).Times(1)
				repo.EXPECT().FetchByUsername(gomock.Any(), gomock.Any()).Return(&User{
					Password: "123",
				}, nil).Times(1)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorAs(t, err, &ErrPasswordWrong)
				return false
			},
			wantToken: "123",
			ctx:       context.Background(),
			req: &LoginRequest{
				username: "123",
				password: "123",
			},
		},
		{
			name: "token生成失败",
			mockFunc: func() {
				encryptService.EXPECT().Encrypt(gomock.Any(), gomock.Any()).Return([]byte("123"), nil).Times(1)
				encryptService.EXPECT().Token(gomock.Any(), gomock.Any()).Return("", errors.New("123")).Times(1)
				repo.EXPECT().FetchByUsername(gomock.Any(), gomock.Any()).Return(&User{
					Password: "123",
				}, nil).Times(1)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorAs(t, err, &ErrLoginFail)
				return false
			},
			wantToken: "123",
			ctx:       context.Background(),
			req: &LoginRequest{
				username: "123",
				password: "123",
			},
		},
	}
	for _, item := range data {
		t.Run(item.name, func(t *testing.T) {
			item.mockFunc()
			got, err := accountUseCase.Login(item.ctx, item.req)
			if !item.wantErr(t, err) {
				return
			}
			assert.Equal(t, item.wantToken, got)
		})
	}
}
