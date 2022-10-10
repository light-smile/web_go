package provider

import (
	"dnds_go/config"
	"dnds_go/global"
	"dnds_go/logger"
	"encoding/json"
	"fmt"
	"testing"
)

func TestInitMqtt(t *testing.T) {

	InitConfig()
	_ = logger.InitLogger(global.Conf.Logger)

	conf1 := config.Mqtt{
		Qos:      1,
		Ip:       "193.168.2.33:", // mqtt ip
		Port:     "1883",
		ClientId: "0001", // clientId
	}
	c1 := InitMqtt(conf1)
	c1.RegisterMqttHandler("test/hello", DefaultMessageHandle())
	// c1.Client.Disconnect(250)
	conf2 := config.Mqtt{
		Qos:      1,
		Ip:       "193.168.2.33:", // mqtt ip
		Port:     "1883",
		ClientId: "0002", // clientId
	}
	c2 := InitMqtt(conf2)
	data := Entity{
		Code: "0",
		Data: "ok",
		Msg:  "success",
	}
	pubData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	c2.Publish("test/hello", false, pubData)

	// c1.Client.Disconnect(250)
}
