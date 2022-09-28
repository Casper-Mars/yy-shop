package biz

import (
	"context"
	"errors"
	"fmt"
	"yy-shop/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrRegisterParamEmpty = errors.New("用户名或者密码不能为空")
	ErrMissingUsername    = errors.New("用户名不能为空")
	ErrMissingPassword    = errors.New("密码不能为空")
	ErrUserNotExist       = errors.New("用户不存在")
	ErrPasswordWrong      = errors.New("密码错误")
	ErrLoginFail          = errors.New("登陆失败")
)

type LoginRequest struct {
	username string
	password string
}

func NewLoginRequest(username, password string) (*LoginRequest, error) {
	// 校验参数
	if username == "" {
		return nil, ErrMissingUsername
	}
	if password == "" {
		return nil, ErrMissingPassword
	}
	return &LoginRequest{
		username: username,
		password: password,
	}, nil
}

type User struct {
	ID       int64  // 用户ID
	Username string // 用户名
	Password string // 密码
	Nickname string // 昵称
	Avatar   string // 头像
}

func (u *User) CheckAuth(ctx context.Context, password string, encryptService EncryptService) error {
	// 校验密码
	encrypt, err := encryptService.Encrypt(ctx, []byte(password))
	if err != nil {
		return err
	}
	if u.Password != string(encrypt) {
		return ErrPasswordWrong
	}
	return nil
}

type UserRepo interface {
	// FetchByUsername 获取指定用户名的用户的信息，如果用户不存在，则返回 ErrUserNotExist。
	FetchByUsername(ctx context.Context, username string) (user *User, err error)
	// FetchByUid 获取指定用户名的用户的信息，如果用户不存在，则返回 ErrUserNotExist。
	FetchByUid(ctx context.Context, uid int64) (user *User, err error)
	// FetchByUidList 获取指定用户名的用户的信息，如果用户不存在，则返回空。
	FetchByUidList(ctx context.Context, uidList []int64) (user map[int64]*User, err error)
	// Save 保存用户信息并返回用户的id。
	Save(ctx context.Context, user *User) (id int64, err error)
}

type AccountUseCase struct {
	authConfig     *conf.Auth
	encryptService EncryptService
	userRepo       UserRepo
	logger         *log.Helper
}

//NewAccountUseCase 创建一个AccountUseCase，依赖作为参数传入
func NewAccountUseCase(logger log.Logger, authConfig *conf.Bootstrap, userRepo UserRepo, encryptService EncryptService) *AccountUseCase {
	return &AccountUseCase{
		encryptService: encryptService,
		userRepo:       userRepo,
		logger:         log.NewHelper(logger),
		authConfig:     authConfig.Auth,
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

//Login 登录，认证成功返回token，认证失败返回错误
func (a *AccountUseCase) Login(ctx context.Context, loginReq *LoginRequest) (token string, err error) {
	// 获取用户信息
	user, err := a.userRepo.FetchByUsername(ctx, loginReq.username)
	if err != nil {
		return "", fmt.Errorf("登录失败：%w", err)
	}
	// 校验密码
	err = user.CheckAuth(ctx, loginReq.password, a.encryptService)
	if err != nil {
		return "", fmt.Errorf("登录失败:%w", err)
	}
	// 生成token
	token, err = a.encryptService.Token(ctx, user)
	if err != nil {
		a.logger.Errorf("登录失败，生成token失败：%v", err)
		return "", ErrLoginFail
	}
	return token, nil
}
