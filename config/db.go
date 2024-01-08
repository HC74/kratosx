package config

import "time"

// Database 数据库配置
type Database struct {
	Enable        bool          // 是否启用数据库
	TablePrefix   string        // 表的前缀
	Drive         string        // 驱动
	Dsn           string        // 数据库链接
	MaxLifetime   time.Duration // 最大生存时间 例如: 2h
	MaxOpenConn   int           // 最大连接数量
	MaxIdleConn   int           // 最大空闲数量
	LogLevel      int           // 日志级别
	PrepareStmt   bool
	DryRun        bool
	SlowThreshold time.Duration // 慢sql阈值
}
