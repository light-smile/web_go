# 

### Description

用于快速web开发框架，集成配置管理、日志、gorm、mqtt、websocket等模块，包含了跨域处理、请求响应封装、请求参数校验，也可以根据自己需求自定义一些功能。

### 安装项目依赖

```
go mod tidy
```

### 安装swagger

```
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go get -u github.com/alecthomas/template
```

### 初始化swagger

```bash
swag init
```

### 代码热更新安装fresh

```bash
go get github.com/pilu/fresh
```

### 安装后终端输入

```bash
fresh
```

### 未安装fresh也可以直接启动

```
go run main.go
```

### 测试是否启动成功

访问localhost:3000/test/hello,收到json数据则启动成功。

### 编程小白一枚，可能存在诸多问题；如遇不懂欢迎交流

### 功能模块介绍

├─common         公共类型、函数封装

├─config         配置文件

├─docs         swagger 文档

├─global         全局变量

├─log         日志文件

├─logger         日志功能模块

├─middleware         中间件

├─provider         提供各种功能，如：数据库、配置功能、mqtt连接等等

├─router         路由地址管理

├─src         业务开发

│  ├─dao

│  │  ├─constant

│  │  ├─login

│  │  └─test

│  ├─models

│  │  ├─constant

│  │  ├─gateway

│  │  ├─login

│  │  ├─test

│  │  └─users

│  └─service

│      ├─constant

│      ├─gateway

│      ├─login

│      └─test

├─test         测试

└─tool         工具函数