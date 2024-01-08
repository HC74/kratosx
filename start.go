package kratosx

import (
	"github.com/HC74/kratosx/core"
	"github.com/HC74/kratosx/core/logger"
	"github.com/HC74/kratosx/middleware"
	"github.com/go-kratos/kratos/v2"
	"os"
)

const (
	defaultFilePath = "config/config.yaml"
)

func New(opts ...Option) *kratos.App {
	o := &options{}

	// 设置默认值
	hostname, _ := os.Hostname()
	o.defaultInit(hostname, defaultFilePath)
	ID(hostname)
	ConfigPath(defaultFilePath)

	// 配置Option项
	o.configNew()

	// 日志Option项
	o.logger()

	for _, opt := range opts {
		opt(o)
	}

	// 加载配置信息,抛出存在异常
	err := o.loadConfig()

	if err != nil {
		panic(err)
	}

	core.Init(o.config, o.loggerFields)

	middlewares := middleware.New(o.config)

	kratosOptions := o.kratosOption()

	if o.services != nil {
		hserver := httpServer(o.config.App().Server.Http, middlewares)
		gserver := grpcServer(o.config.App().Server.Grpc, middlewares)
		o.services(o.config, hserver, gserver)
		kratosOptions = append(kratosOptions, kratos.Server(hserver, gserver))
	}

	if o.config.App().Log != nil {
		// 初始化logger
		kratosOptions = append(kratosOptions, kratos.Logger(logger.Instance()))
	}

	kratosOptions = append(kratosOptions, o.kratosOptions...)

	return kratos.New(
		kratosOptions...,
	)
}
