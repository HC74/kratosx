package middleware

import (
	"github.com/HC74/kratosx/config"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
)

func New(conf config.Config) []middleware.Middleware {
	app := conf.App()
	mds := []middleware.Middleware{
		recovery.Recovery(),
		Logging(app.Logging),
		JwtWhite(app.Jwt), // jwt白名单
		JwtBlack(app.Jwt), // jwt校验
		validate.Validator(),
		metadata.Server(),
	}
	// 原地删除不启用的中间件
	return removeDisableMiddleware(mds)
}

func removeDisableMiddleware(slice []middleware.Middleware) []middleware.Middleware {
	fast, slow := 0, 0
	for fast < len(slice) {
		if slice[fast] != nil {
			slice[slow] = slice[fast]
			slow++
		}
		fast++
	}
	return slice[:slow]
}
