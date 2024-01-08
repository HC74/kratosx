package config

import (
	kc "github.com/go-kratos/kratos/v2/config"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/proto"
)

var mars *config

type Config interface {
	Load() error
	Scan(v interface{}) error
	Value(key string) Value
	Close() error
	App() *App
}

type config struct {
	app *App
	mar kc.Config
}

// Instance 回去配置实例
func Instance() Config {
	return mars
}

func New(source kc.Source) Config {
	mars = &config{
		mar: kc.New(
			kc.WithSource(source),
		),
	}
	return mars
}

func (c *config) Load() error {
	if err := c.mar.Load(); err != nil {
		return err
	}
	c.app = new(App)
	return c.Scan(c.app)
}

func (c *config) App() *App {
	return c.app
}

func (c *config) Scan(dst any) error {
	if _, ok := dst.(proto.Message); ok {
		return c.mar.Scan(&dst)
	}

	// 序列化json
	res := map[string]any{}
	if err := c.mar.Scan(&res); err != nil {
		return err
	}

	dc := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           dst,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	decoder, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return err
	}
	return decoder.Decode(res)
}

func (c *config) transformValue(val kc.Value) Value {
	return &value{Value: val}
}

func (c *config) Value(key string) Value {
	return c.transformValue(c.mar.Value(key))
}

func (c *config) Close() error {
	return c.mar.Close()
}
