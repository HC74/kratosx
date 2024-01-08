package server

import (
	"github.com/HC74/kratosx/config"
	thttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"time"
)

// Cors 跨域配置
func Cors(conf *config.Cors) thttp.ServerOption {
	maxAge := time.Minute * 8
	if conf.MaxAge != 0 {
		maxAge = conf.MaxAge
	}
	opts := []handlers.CORSOption{
		handlers.AllowedOrigins(conf.AllowOrigins),
		handlers.AllowedMethods(conf.AllowMethods),
		handlers.AllowedHeaders(conf.AllowHeaders),
		handlers.ExposedHeaders(conf.ExposeHeaders),
		handlers.MaxAge(int(maxAge.Seconds())),
	}

	if conf.AllowCredentials {
		opts = append(opts, handlers.AllowCredentials())
	}

	// 加入kratos拦截器，解决跨域
	return thttp.Filter(handlers.CORS(opts...))
}
