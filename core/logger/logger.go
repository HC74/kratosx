package logger

import (
	"github.com/HC74/kratosx/config"
	"os"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogField map[string]any

type logger struct {
	log log.Logger
}

var instance *logger

func Instance() log.Logger {
	if instance == nil {
		return log.GetLogger()
	}
	return instance.log
}

func Helper() *log.Helper {
	return log.NewHelper(instance.log)
}

// Init 初始化日志器
func Init(lc *config.Logger, fields LogField) {
	// 没配置则跳过
	if lc == nil {
		return
	}

	// log field 转换
	var fs []any
	for key, val := range fields {
		fs = append(fs, key, val)
	}

	// 初始化
	instance = &logger{}
	instance.initFactory(lc, fs)
}

func (l *logger) initFactory(conf *config.Logger, fs []any) {
	// 创建zap logger
	l.log = log.With(l.newZapLogger(conf), fs...)
	// 设置全局logger
	log.SetLogger(instance.log)
}

func (l *logger) newZapLogger(conf *config.Logger) *kzap.Logger {
	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                                                 // 设置时间字段的键名为 time
		LevelKey:       "level",                                                // 设置日志级别字段的键名为 level
		NameKey:        "log",                                                  // 设置记录器名字字段的键名为 log
		CallerKey:      "caller",                                               // 设置调用者信息字段的键名为 caller
		MessageKey:     "msg",                                                  // 设置日志消息字段的键名为 msg
		LineEnding:     zapcore.DefaultLineEnding,                              // 使用系统默认的换行
		EncodeLevel:    zapcore.LowercaseLevelEncoder,                          // 小写编码器
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"), // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,                         // 持续时间的编码，使用以秒为单位的编码器
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 输出器配置
	var output []zapcore.WriteSyncer
	for _, val := range conf.Output {
		if val == "stdout" {
			output = append(output, zapcore.AddSync(os.Stdout))
		}
		if val == "file" {
			output = append(output, zapcore.AddSync(&lumberjack.Logger{
				Filename:   conf.File.Path,
				MaxSize:    conf.File.MaxSize,
				MaxBackups: conf.File.MaxBackup,
				MaxAge:     conf.File.MaxAge,
				Compress:   conf.File.Compress,
				LocalTime:  conf.File.LocalTime,
			}))
		}
	}

	// 使用go zap
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),  // 编码器配置
		zapcore.NewMultiWriteSyncer(output...), // 输出方式
		zapcore.Level(conf.Level),              // 设置日志级别
	)

	return kzap.NewLogger(zap.New(core))
}
