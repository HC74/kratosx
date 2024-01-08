package db

import (
	"github.com/HC74/kratosx/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
)

const (
	_mysql = "mysql"
	_mssql = "mssql"
)

type db struct {
	mapper map[string]*gorm.DB // 数据库
	key    string              // 键
	lock   sync.RWMutex        // 读写锁
}

type DB interface {
	// Get 获取指定名称的db实例，不指定名称则返回第一个如果实例不存在则返回nil
	Get(name ...string) *gorm.DB
}

var instance *db

func Instance() DB {
	return instance
}

func Init(confdb map[string]*config.Database) {
	if len(confdb) == 0 {
		return
	}

	instance = &db{
		lock:   sync.RWMutex{},            // 创建读写锁
		mapper: make(map[string]*gorm.DB), // 创建数据库操作体和key的关联
	}

	// 遍历配置连接数据库
	for key, conf := range confdb {
		if conf == nil {
			continue
		}

		if err := instance.createDB(key, conf); err != nil {
			panic("database init error :" + err.Error())
		}

	}

	// 如果配置了多个库，则不能启用快速获取
	if len(instance.mapper) != 1 {
		instance.key = ""
	}
}

// createDB 初始化数据库
func (d *db) createDB(name string, conf *config.Database) error {
	if !conf.Enable {
		// 未启用此db
		return nil
	}
	client, err := gorm.Open(d.open(conf.Drive, conf.Dsn), &gorm.Config{
		Logger: newLog(conf.LogLevel, conf.SlowThreshold),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	sdb, _ := client.DB()
	sdb.SetConnMaxLifetime(conf.MaxLifetime)
	sdb.SetMaxOpenConns(conf.MaxOpenConn)
	sdb.SetMaxIdleConns(conf.MaxIdleConn)

	d.lock.Lock()
	d.mapper[name] = client
	d.key = name
	d.lock.Unlock()
	return nil
}

// Get 根据配置名回去DB实例
func (d *db) Get(name ...string) *gorm.DB {
	if d.key == "" && len(name) == 0 {
		return nil
	}
	// 读锁
	d.lock.RLock()
	defer d.lock.RUnlock()

	key := d.key
	if len(name) != 0 {
		key = name[0]
	}
	value := d.mapper[key]
	return value
}

func (d *db) delete(name string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.mapper, name)
}

// open 打开数据库链接
func (d *db) open(diver, dns string) gorm.Dialector {
	switch diver {
	case _mysql:
		return mysql.Open(dns)
	case _mssql:
		return sqlserver.Open(dns)
	default:
		return nil
	}
}
