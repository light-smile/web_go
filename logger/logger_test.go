package logger

import (
	"dnds_go/config"
	"dnds_go/global"
	"errors"
	"time"

	"testing"

	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	var conf = config.Logger{
		LogLevel:          "debug",
		LogFormat:         "logfmt",
		LogPath:           "./log",
		LogFileName:       "test.log",
		LogFileMaxSize:    1,
		LogFileMaxBackups: 10,
		LogMaxAge:         1000,
		LogCompress:       false,
		LogStdout:         true,
	}
	// 2. 初始化log
	if err := InitLogger(conf); err != nil {
		t.Fatal(err)
	}

	// 3. 调用 Logger 打印日志测试
	global.Logger.Info("111")                // logger Infof 用法
	zap.S().Debugf("测试 Debugf 用法：%s", "111") // logger Debugf 用法
	go func() {
		for i := 0; i < 100000; i++ {
			zap.S().Infof("(1)协程内部调用测试 Infof 用法：%s", "111")
			time.Sleep(time.Millisecond)
		}
	}()
	zap.S().Errorf("测试 Errorf 用法：%s", "111") // logger Errorf 用法
	zap.S().Warnf("测试 Warnf 用法：%s", "111")   // logger Warnf 用法
	zap.S().Infof("测试 Infof 用法：%s, %d, %v, %f", "111", 1111, errors.New("collector returned no data"), 3333.33)
	// logger With 用法
	logger := zap.S().With("collector", "cpu", "name", "主机")
	logger.Infof("测试 (With + Infof) 用法：%s", "测试")
	zap.S().Errorf("测试 Errorf 用法：%s", "111")
	go func() {
		for i := 0; i < 10000; i++ {
			zap.S().Infof("(2)协程内部调用测试 Infof 用法：%s", "111")
			time.Sleep(time.Millisecond)
		}
	}()
	time.Sleep(5 * time.Second)
}
