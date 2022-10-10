package config

// 以结构体类型定义配置信息，每次在yaml文件中添加内容，需在结构体也响应添加
type Config struct {
	Server
	Database
	Logger
	Mqtt
}

type Server struct {
	Name string
	Mode string
	Port string
}

type Database struct {
	FileName string // 数据库文件名称
}
type Mqtt struct {
	Qos      byte
	Ip       string // mqtt ip
	Port     string
	ClientId string // clientId
	UserName string
	Password string
}

type Logger struct {
	ShowReqAndRes     bool   // 是否显示请求数据和响应数据
	LogLevel          string // 日志打印级别 debug  info  warning  error
	LogFormat         string // 输出日志格式	logfmt, json
	LogPath           string // 输出日志文件路径
	LogFileName       string // 输出日志文件名称
	LogFileMaxSize    int    // 【日志分割】单个日志文件最多存储量 单位(mb)
	LogFileMaxBackups int    // 【日志分割】日志备份文件最多数量
	LogMaxAge         int    // 日志保留时间，单位: 天 (day)
	LogCompress       bool   // 是否压缩日志
	LogStdout         bool   // 是否输出到控制台
}
