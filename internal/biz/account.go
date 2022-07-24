package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrRegisterParamEmpty = errors.New("用户名或者密码不能为空")
	ErrUserNotExist       = errors.New("用户不存在")
)

type User struct {
	ID       int64
	Username string
	Password string
	Nickname string
	Avatar   string
}

//CheckPassword 校验密码是否正确
func (u *User) CheckPassword(ctx context.Context, pwd string, encryptService EncryptService) (valid bool) {
	encrypt, err := encryptService.Encrypt(ctx, []byte(pwd))
	if err != nil {
		return false
	}
	if u.Password != string(encrypt) {
		return false
	}
	return true
}

type UserRepo interface {
	Fetch(ctx context.Context, id int64) (user *User, err error)
	Save(ctx context.Context, user *User) (id int64, err error)
}

type AccountUseCase struct {
	encryptService EncryptService
	userRepo       UserRepo
	logger         *log.Helper
}

func NewAccountUseCase(logger log.Logger, userRepo UserRepo, encryptService EncryptService) *AccountUseCase {
	return &AccountUseCase{
		encryptService: encryptService,
		userRepo:       userRepo,
		logger:         log.NewHelper(logger),
	}
}

//Register 注册
func (a *AccountUseCase) Register(ctx context.Context, username, pwd string) (err error) {
	// 校验参数
	if username == "" || pwd == "" {
		return fmt.Errorf("注册失败：%w", ErrRegisterParamEmpty)
	}
	// 加密密码
	encrypt, err := a.encryptService.Encrypt(ctx, []byte(pwd))
	if err != nil {
		log.Errorf("注册失败，参数[username: %s，pwd:%s]，err:%v", username, pwd, err)
		return fmt.Errorf("注册失败")
	}
	_, err = a.userRepo.Save(ctx, &User{
		Username: username,
		Password: string(encrypt),
	})
	if err != nil {
		return fmt.Errorf("注册失败：%w", err)
	}
	return nil
}
