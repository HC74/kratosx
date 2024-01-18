package middleware

import (
	"context"
	"fmt"
	"github.com/HC74/kratosx/config"
	"github.com/HC74/kratosx/core/jwt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	kratosJwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"strings"
)

func JwtWhite(conf *config.JWT) middleware.Middleware {
	if conf == nil {
		return nil
	}

	keyFunc := func(token *jwtv4.Token) (any, error) {
		return []byte(conf.Secret), nil
	}

	whiteList := func(ctx context.Context) bool {
		operation, path := "", ""
		if tr, isok := transport.FromServerContext(ctx); isok {
			operation = tr.Operation()
		}
		if req, isok := http.RequestFromServerContext(ctx); isok {
			path = fmt.Sprintf("%s:%s", req.Method, req.URL.Path)
			// path 对应config中的jwt配置中的白名单
		}

		jwtIns := jwt.Instance()
		return jwtIns.IsWhitelist(path) || jwtIns.IsWhitelist(operation)
	}

	return selector.Server(kratosJwt.Server(keyFunc)).Match(func(ctx context.Context, operation string) bool {
		return !whiteList(ctx)
	}).Build()
}

func JwtBlack(conf *config.JWT) middleware.Middleware {
	if conf == nil {
		return nil
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			header, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			auths := strings.SplitN(header.RequestHeader().Get("Authorization"), " ", 2)
			if len(auths) != 2 || !strings.EqualFold(auths[0], "Bearer") {
				return handler(ctx, req)
			}

			token := auths[1]

			// 判断token是否在黑名单内
			jwtIns := jwt.Instance()
			if jwtIns.IsBlacklist(token) {
				return nil, errors.Unauthorized("UNAUTHORIZED", "JWT token is lose efficacy")
			}

			ctx = jwtIns.SetToken(ctx, token)
			return handler(ctx, req)
		}
	}
}
