package global

import (
	"dnds_go/config"
	"sync"

	// "dnds_go/provider"

	ut "github.com/go-playground/universal-translator"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
	全局变量可以保存在当前文件
*/
var (
	// DB GormDefaultDb mysql默认连接
	DB     *gorm.DB       // 全局数据库
	L      sync.Mutex     // 全局读写锁
	Logger *zap.Logger    // 全局日志
	Conf   *config.Config // 全局配置
	Trans  ut.Translator  // 全局翻译器
	// TimeFormatYmdhis MqClent *provider.Agent
	// TimeFormatYmdhis 时间格式 yyyy-mm-dd hh:ii:ss
	TimeFormatYmdhis = "2006-01-02 15:04:05"

	// TimeFormatYmd 时间格式 yyyy-mm-dd
	TimeFormatYmd = "2006-01-02"
)

// var WL sync.Mutex
