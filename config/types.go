package config

import "time"

type App struct {
	Server   *Server              // 服务配置
	Log      *Logger              // 日志配置
	Logging  *Logging             // 日志白名单配置
	Database map[string]*Database // 数据库配置
}

// Server 服务配置
type Server struct {
	// Http服务配置
	Http *HttpService
	// Grpc服务配置
	Grpc *GrpcService
}

// GrpcService grpc服务
type GrpcService struct {
	Network string        // 配置服务端的 network 协议，如 tcp
	Addr    string        // 配置服务端监听的地址 例如: 0.0.0.0:8080
	Timeout time.Duration // 配置服务端的超时设置 例如: 10s
}

// HttpService http服务
type HttpService struct {
	Network        string
	Addr           string        // 配置服务端的 network 协议，如 tcp
	Timeout        time.Duration // 配置服务端的超时设置 例如: 10s
	FormatResponse bool          // 默认格式化返回值
	Cors           *Cors         // 跨域
	Marshal        *Marshal
}

// Cors 跨域设置
type Cors struct {
	AllowCredentials    bool          // 跨域请求中允许传递身份验证凭据
	AllowOrigins        []string      // 允许的跨域请求的来源
	AllowMethods        []string      // 允许的请求方式
	AllowHeaders        []string      // 请求头
	ExposeHeaders       []string      // 在跨域请求中允许暴露哪些头信息
	MaxAge              time.Duration // 超时
	AllowPrivateNetwork bool
}

// Marshal 序列化设置
type Marshal struct {
	ForceUseJson    bool // 是否强制使用 JSON 格式进行序列化
	EmitUnpopulated bool // 是否在编组时包含未赋值的字段
	UseProtoNames   bool // 否使用 Protocol Buffers 中定义的字段名进行序列化
}

type Logging struct {
	Enable    bool            // 是否开启白名单
	Whitelist map[string]bool // 路由 -> 是否为白名单 白名单不记录
}

// Logger 日志配置项
type Logger struct {
	Level  int8      // 日志级别
	Output []string  // 输出形式,stdout:输出到控制台，file:输出到文件
	File   *struct { // 如果输出形式为文件，所对应的文件配置
		Path      string // 文件输出的目录文件，例: ./logs/output.log
		MaxSize   int    // 日志文件最大保留(M)
		MaxBackup int    // 日志文件最多保存个数
		MaxAge    int    // 日志最多保存的天数
		Compress  bool   // 是否要进行日志压缩，using gzip
		LocalTime bool   // 是否使用本地的时间
	}
}
