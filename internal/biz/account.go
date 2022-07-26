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
	ID       int64  // 用户ID
	Username string // 用户名
	Password string // 密码
	Nickname string // 昵称
	Avatar   string // 头像
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
	// FetchByUsername 获取指定用户名的用户的信息，如果用户不存在，则返回 ErrUserNotExist。
	FetchByUsername(ctx context.Context, username string) (user *User, err error)
	// Save 保存用户信息并返回用户的id。
	Save(ctx context.Context, user *User) (id int64, err error)
}

type AccountUseCase struct {
	encryptService EncryptService
	userRepo       UserRepo
	logger         *log.Helper
}

//NewAccountUseCase 创建一个AccountUseCase，依赖作为参数传入
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
	// 判断用户是否已经注册一次了
	user, err := a.userRepo.FetchByUsername(ctx, username)
	if err != nil && !errors.Is(err, ErrUserNotExist) {
		log.Errorf("注册失败，参数[username: %s，pwd:%s]，err:%v", username, pwd, err)
		return fmt.Errorf("注册失败")
	}
	if user != nil {
		return fmt.Errorf("用户已经存在")
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
