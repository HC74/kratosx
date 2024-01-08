package middleware

import (
	"context"
	kc "github.com/HC74/kratosx/config"
	"github.com/HC74/kratosx/core/logger"
	lg "github.com/HC74/kratosx/core/logging"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func Logging(conf *kc.Logging) middleware.Middleware {
	if conf == nil || !conf.Enable {
		return nil
	}

	return selector.Server(logging.Server(logger.Instance())).Match(func(ctx context.Context, operation string) bool {
		path := ""
		if h, is := http.RequestFromServerContext(ctx); is {
			path = h.Method + ":" + h.URL.Path
		}
		lgIns := lg.Instance()
		return !(lgIns.IsWhitelist(operation) || lgIns.IsWhitelist(path))
	}).Build()
}
