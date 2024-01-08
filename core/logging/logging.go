package logging

import (
	"github.com/HC74/kratosx/config"
	"sync"
)

type logging struct {
	mu  sync.RWMutex
	set map[string]bool
}

type Logging interface {
	IsWhitelist(path string) bool
}

var instance *logging

func Instance() Logging {
	return instance
}

func Init(ec *config.Logging) {
	if ec == nil {
		return
	}

	instance = &logging{
		mu:  sync.RWMutex{},
		set: ec.Whitelist,
	}

}

func (c *logging) IsWhitelist(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.set[name]
}
