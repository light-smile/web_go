package main

import (
	"context"
	"dnds_go/global"
	"dnds_go/logger"
	"dnds_go/provider"
	"dnds_go/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

type Error struct {
	Name string
	Msg  string
}

// @title dnds_go
// @version 1.0
// @description dnds_go 测试版框架
// @license.name go
// @contact.name go-swagger帮助文档
// @contact.url http://localhost:3000/swagger/index.html
// @host localhost:3000
// @BasePath /
func main() {
	// 初始化配置文件
	provider.InitConfig()
	// 初始化数据库
	provider.InitDatabase()

	// 初始化数据库表
	provider.InitTable()

	// 初始化 validator 翻译
	provider.InitTrans()

	// 初始化日志
	if err := logger.InitLogger(global.Conf.Logger); err != nil {
		log.Fatalln("logger init error: ", err)
	}

	// 初始化mqtt
	provider.MqClient = provider.InitMqtt(global.Conf.Mqtt)
	// provider.MqClient.Connect()
	// 注册一个默认mqtt客户端
	// provider.MqClient.RegisterMqttHandler("test/hello", provider.DefaultMessageHandle())
	// client.Disconnect(250)

	// 3. 调用 Logger 打印日志测试
	// zap.S().Infof("测试 Infof 用法：%s", "111")   // logger Infof 用法
	// zap.S().Debugf("测试 Debugf 用法：%s", "111") // logger Debugf 用法
	// // go func() {
	// // 	for i := 0; i < 100000; i++ {
	// // 		zap.S().Infof("(1)协程内部调用测试 Infof 用法：%s", "111")
	// // 		time.Sleep(time.Millisecond)
	// // 	}
	// // }()
	// zap.S().Errorf("测试 Errorf 用法：%s", "111") // logger Errorf 用法
	// zap.S().Warnf("测试 Warnf 用法：%s", "111")   // logger Warnf 用法
	// zap.S().Infof("测试 Infof 用法：%s, %d, %v, %f", "111", 1111, errors.New("collector returned no data"), 3333.33)
	// // logger With 用法
	// logger := zap.S().With("collector", "cpu", "name", "主机")
	// logger.Infof("测试 (With + Infof) 用法：%s", "测试")
	// zap.S().Errorf("测试 Errorf 用法：%s", "111")
	// 服务配置
	router := router.InitRouter()
	server := &http.Server{
		Addr:           viper.GetString("Server.Port"),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Listen server error: ", err, syscall.Getpid())
		}
	}()

	// c.Client.Disconnect(250)
	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// helper.GetLogger("").Warnln("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		// helper.GetLogger("").Fatalln("Shutdown server error: ", err)
	}
	// helper.GtLogger("").Warnln("Server exiting")
}
