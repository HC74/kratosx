# Kratosx 基于kratos的二次开发
### Installing
##### go install 安装：
```
go install github.com/HC74/kratosx/cmd/kratosx/v2@latest
```

### Create a service
```
# 创建项目模板
kratos new project

cd project
# 拉取项目依赖
go mod download / go mod tidy

# 运行程序
kratos run
```

## 项目结构
the generated files look like:

   ```Plain Text
   ├── project
   │   ├── api                       // 编写grpc的proto文件
   │   │── cmd                       // 程序的执行入口
   │   │   └── project
   │   │       └── main.go           // 程序入口
   │   │── config                    // 配置文件目录
   │   │   │── config.go             // 自定义配置文件go
   │   │   └── config.yaml           // 配置文件,包含默认配置
   │   └── internal
   │       ├── model
   │       │   └── test.go           // 数据库db文件
   │       ├── logic
   │       │   ├── logic.go          // 业务逻辑的抽象大类
   │       │   └── greeter.go        // 业务逻辑具体的实现，可多个
   │       ├── service
   │       │   ├── service.go        // 调用底层业务逻辑代码
   │       │   └── greeter.go        // 整合业务逻辑抽象
   └────────────────────────────────────────────────────End
   ```