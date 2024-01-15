package core

import (
	"github.com/HC74/kratosx/config"
	"github.com/HC74/kratosx/core/db"
	"github.com/HC74/kratosx/core/logger"
	"github.com/HC74/kratosx/core/logging"
	rds "github.com/HC74/kratosx/core/redis"
)

func Init(conf config.Config, fs logger.LogField) {
	// 初始化全局日志
	logger.Init(conf.App().Log, fs)

	// 初始化数据库
	db.Init(conf.App().Database)

	// 初始化redis数据库
	rds.Init(conf.App().Redis)

	// logging 初始化
	logging.Init(conf.App().Logging)
}
