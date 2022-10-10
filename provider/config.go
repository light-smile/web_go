package provider

import (
	"dnds_go/global"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局使用配置

// InitConfig Config 配置
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Read config failed: ", err)
	}

	// 监听配置文件，当文件内容发生改变时，自动更新
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed: ", e.Name)
	})

	viper.Unmarshal(&global.Conf)
}
