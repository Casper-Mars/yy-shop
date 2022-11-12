package server

import (
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwt2 "github.com/golang-jwt/jwt/v4"
	v1 "yy-shop/api/v1"
	"yy-shop/internal/biz"
	"yy-shop/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Bootstrap,
	logger log.Logger,
	as v1.AccountServer,
	ps v1.ProductServer,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			selector.Server(
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(c.GetAuth().GetJwtSecret()), nil
				}, jwt.WithClaims(func() jwt2.Claims {
					return &biz.MyJwtClaims{}
				})),
			).Match(whiteList(c)).Build(),
		),
	}
	if c.GetServer().Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.GetServer().Grpc.Network))
	}
	if c.GetServer().Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.GetServer().Grpc.Addr))
	}
	if c.GetServer().Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.GetServer().Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterAccountServer(srv, as)
	v1.RegisterProductServer(srv, ps)
	return srv
}
