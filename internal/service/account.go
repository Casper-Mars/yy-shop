package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "yy-shop/api/account/v1"
	"yy-shop/internal/biz"
)

type accountService struct {
	v1.UnimplementedAccountServer
	log *log.Helper
	auc *biz.AccountUseCase
}

func NewAccountService(logger log.Logger, auc *biz.AccountUseCase) v1.AccountServer {
	return &accountService{
		log: log.NewHelper(logger),
		auc: auc,
	}
}

func (a *accountService) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	loginRequest, err := biz.NewLoginRequest(request.GetPhone(), request.GetPassword())
	if err != nil {
		return nil, errors.New(500, "登录失败", err.Error())
	}
	token, err := a.auc.Login(ctx, loginRequest)
	if err != nil {
		return nil, errors.New(500, "登录失败", err.Error())
	}
	return &v1.LoginResponse{
		Token: token,
	}, nil
}

func (a *accountService) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	err := a.auc.Register(ctx, request.GetPhone(), request.GetPassword())
	if err != nil {
		return nil, errors.New(500, "注册失败", err.Error())
	}
	return &v1.RegisterResponse{}, nil
}
