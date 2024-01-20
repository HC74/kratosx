package loader

import (
	"io"
	"os"
	"sync"
)

type loader struct {
	mu  sync.RWMutex
	set map[string][]byte
}

type Loader interface {
	Get(name string) []byte
}

var instance *loader

func Instance() Loader {
	return instance
}

func Init(conf map[string]string) {
	// 不存在跳过初始化
	if len(conf) == 0 {
		return
	}

	instance = &loader{
		mu:  sync.RWMutex{},
		set: make(map[string][]byte),
	}

	// 连接数据库
	for key, path := range conf {
		if err := instance.initFactory(key, path); err != nil {
			panic("加载器初始化失败:" + err.Error())
		}
	}
}

func (c *loader) Get(name string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.set[name]
}

func (c *loader) initFactory(name string, path string) error {
	// 获取文件内容
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	all, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	defer file.Close()

	c.mu.Lock()
	c.set[name] = all
	c.mu.Unlock()
	return nil
}

func (c *loader) delete(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.set, name)
}
