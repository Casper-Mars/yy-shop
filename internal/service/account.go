package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	v1 "yy-shop/api/v1"
	"yy-shop/internal/biz"
)

type accountService struct {
	log *log.Helper
	auc *biz.AccountUseCase
}

func NewAccountService(logger log.Logger, auc *biz.AccountUseCase) v1.AccountServer {
	return &accountService{
		log: log.NewHelper(logger),
		auc: auc,
	}
}

func (a *accountService) Login(ctx context.Context, request *v1.LoginReq) (*v1.LoginResp, error) {
	loginRequest, err := biz.NewLoginRequest(request.GetPhone(), request.GetPassword())
	if err != nil {
		return nil, errors.New(500, "登录失败", err.Error())
	}
	token, err := a.auc.Login(ctx, loginRequest)
	if err != nil {
		return nil, errors.New(500, "登录失败", err.Error())
	}
	return &v1.LoginResp{
		Token: token,
	}, nil
}

func (a *accountService) Register(ctx context.Context, request *v1.RegisterReq) (*v1.RegisterResp, error) {
	err := a.auc.Register(ctx, request.GetPhone(), request.GetPassword())
	if err != nil {
		return nil, errors.New(500, "注册失败", err.Error())
	}
	return &v1.RegisterResp{}, nil
}

func (a *accountService) Info(ctx context.Context, request *v1.InfoReq) (*v1.InfoResp, error) {
	claims, _ := jwt.FromContext(ctx)
	userInfo, err := a.auc.UserInfo(ctx, claims.(*biz.MyJwtClaims).Uid)
	if err != nil {
		return nil, errors.New(500, "获取用户信息失败", err.Error())
	}
	return &v1.InfoResp{
		Id:       userInfo.ID,
		Username: userInfo.Username,
		Avatar:   userInfo.Avatar,
	}, nil
}
