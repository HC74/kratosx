package kratosx

import (
	"github.com/HC74/kratosx/config"
	"github.com/HC74/kratosx/core/logger"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Option is an application option.
type Option func(o *options)
type RegisterServerFn func(config config.Config, hs *http.Server, gs *grpc.Server)

type options struct {
	id   string
	name string
	// 版本
	version string
	// 配置文件地址
	configPath string

	// 配置
	config config.Config
	// 日志字段
	loggerFields  logger.LogField
	services      RegisterServerFn
	kratosOptions []kratos.Option
}

// ID with service id.
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Options kratos option
func Options(opts ...kratos.Option) Option {
	return func(o *options) {
		o.kratosOptions = opts
	}
}

// RegisterServer 注册服务
func RegisterServer(service RegisterServerFn) Option {
	return func(o *options) {
		o.services = service
	}
}

// ConfigPath 文件配置路径，默认为config/config.yaml
func ConfigPath(path string) Option {
	return func(o *options) {
		o.configPath = path
	}
}

// configNew 加入config
func (o *options) configNew() {
	o.config = config.New(file.NewSource(o.configPath))
}

// defaultInit 默认配置
func (o *options) defaultInit(name, confPath string) {
	o.name = name
	o.configPath = confPath
}

// loadConfig 加载配置文件
func (o *options) loadConfig() error {
	err := o.config.Load()
	return err
}

// logger 装配logger日志
func (o *options) logger() {
	o.loggerFields = logger.LogField{
		"id":      o.id,
		"name":    o.name,
		"version": o.version,
		"trace":   tracing.TraceID(),
		"span":    tracing.SpanID(),
	}
}

// kratosOption kratos默认的option
func (o *options) kratosOption() []kratos.Option {
	defaultOption := []kratos.Option{
		kratos.ID(o.id),
		kratos.Name(o.name),
		kratos.Version(o.version),
		kratos.Metadata(map[string]string{}),
	}
	return defaultOption
}
