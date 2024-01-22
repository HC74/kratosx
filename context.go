package kratosx

import (
	"context"
	"github.com/HC74/kratosx/config"
	"github.com/HC74/kratosx/core/db"
	"github.com/HC74/kratosx/core/jwt"
	"github.com/HC74/kratosx/core/loader"
	"github.com/HC74/kratosx/core/logger"
	rds "github.com/HC74/kratosx/core/redis"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type Context interface {
	Logger() *log.Helper

	Ctx() context.Context
	//GetMetadata(string) string
	//SetMetadata(key, value string)

	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Config() config.Config
	Endpoint() []string
	DB(name ...string) *gorm.DB
	Redis(name ...string) *redis.Client
	Loader(name string) []byte
	JWT() jwt.Jwt

	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}

type ctx struct {
	context.Context // 上下文
	kratos.AppInfo  // kratos实例
}

// MustContext returns the Transport value stored in ctx, if any.
func MustContext(c context.Context) Context {
	app, _ := kratos.FromContext(c)
	return &ctx{
		Context: c,
		AppInfo: app,
	}
}

// Ctx 获取上下文
func (c *ctx) Ctx() context.Context {
	return c.Context
}

func (c *ctx) Loader(name string) []byte {
	return loader.Instance().Get(name)
}

// Logger 获取日志处理器
func (c *ctx) Logger() *log.Helper {
	helper := logger.Helper()
	return helper.WithContext(c)
}

// Redis 获取redis
func (c *ctx) Redis(name ...string) *redis.Client {
	return rds.Instance().Get(name...)
}

// DB 获取数据库
func (c *ctx) DB(name ...string) *gorm.DB {
	return db.Instance().Get(name...).WithContext(c.Ctx())
}

// JWT 获取令牌验证器
func (c *ctx) JWT() jwt.Jwt {
	return jwt.Instance()
}

// Config 获取配置对象
func (c *ctx) Config() config.Config {
	return config.Instance()
}

func (c *ctx) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Deadline()
}

func (c *ctx) Done() <-chan struct{} {
	return c.Context.Done()
}

func (c *ctx) Err() error {
	return c.Context.Err()
}

func (c *ctx) Value(key any) any {
	return c.Context.Value(key)
}
