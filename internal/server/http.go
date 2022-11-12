package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwt2 "github.com/golang-jwt/jwt/v4"
	v1 "yy-shop/api/v1"
	"yy-shop/internal/biz"
	"yy-shop/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var whiteList = func(c *conf.Bootstrap) selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		for _, v := range c.GetAuth().GetWhiteList() {
			if v == operation {
				return false
			}
		}
		return true
	}
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *conf.Bootstrap,
	as v1.AccountServer,
	ps v1.ProductServer,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
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
	if c.GetServer().Http.Network != "" {
		opts = append(opts, http.Network(c.GetServer().Http.Network))
	}
	if c.GetServer().Http.Addr != "" {
		opts = append(opts, http.Address(c.GetServer().Http.Addr))
	}
	if c.GetServer().Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.GetServer().Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterAccountHTTPServer(srv, as)
	v1.RegisterProductHTTPServer(srv, ps)
	return srv
}
