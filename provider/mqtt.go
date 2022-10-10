package provider

import (
	"crypto/tls"
	"dnds_go/config"
	"dnds_go/global"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	// "go.uber.org/zap/zapcore"
)

var MqClient *Agent

// mqtt 配置信息
var conf config.Mqtt

// Agent runs an mqtt client
type Agent struct {
	Client MQTT.Client
}

type subscription struct {
	topic   string
	handler MQTT.MessageHandler
}

// NewAgent creates an agent
func InitMqtt(cf config.Mqtt) *Agent {
	conf = cf
	a := new(Agent)
	//opts:= MQTT.NewClientOptions().AddBroker(MQTTconfig.conf.Ip).SetClientID(MQTTconfig.conf.ClientId)

	// brokerInfo: broker地址
	brokerInfo := conf.Ip + conf.Port
	opts := MQTT.NewClientOptions().AddBroker(brokerInfo).SetClientID(conf.ClientId)
	//opts.SetUsername(MQTTconfig.MQTTUser)
	//opts.SetPassword(MQTTconfig.MQTTPasswd)
	opts.SetKeepAlive(5 * time.Second)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetUsername("dnds")
	opts.SetTLSConfig(&tls.Config{
		ClientAuth:         tls.NoClientCert, // 0: 不需要证书
		ClientCAs:          nil,
		InsecureSkipVerify: true,
	})

	opts.OnConnectionLost = func(c MQTT.Client, err error) {
		// log.WithField("error", err).Info("Lost connection")
		// Logger.Error("connect mqtt fail", zap.Error(err))
		zap.S().Errorf("mqtt disconnect: %s", err.Error())
	}

	a.Client = MQTT.NewClient(opts)

	if token := a.Client.Connect(); token.Wait() && token.Error() != nil {
		global.Logger.Error("token err", zap.Error(token.Error()))
		// panic(token.Error())
	}
	// // mqtt.DEBUG = log.New(os.Stdout, "", 0)
	// // mqtt.ERROR = log.New(os.Stdout, "", 0)
	zap.S().Infof("ClientId:%s 连接到broker", conf.ClientId)

	return a
	// return nil
}

type HandlerFunc func(string, *Entity) *myRespEntity

// 响应结构体
type Entity struct {
	Code string
	Data interface{}
	Msg  string
}

// 指定主题，注册处理方法
func (a *Agent) RegisterMqttHandler(topic string, handler HandlerFunc) {
	if token := a.Client.Subscribe(topic, conf.Qos,
		func(c MQTT.Client, msg MQTT.Message) {
			// var msgmqtt Entity
			var respentity *myRespEntity
			// var ic infra.IComponent
			var msgdata *Entity
			//unmarshal the bgp json, execute action and return status
			for true {
				LogPush(topic, msg.Payload())
				if err := json.Unmarshal(msg.Payload(), &msgdata); err != nil {
					global.Logger.Error("RegisterMqttHandler Unmarshal fail")
					break
				}
				//Syslog.LogOutObj(msg.Topic(), msgBgp, Syslog.RESP_OK)
				// msgdate.Entity = &msgmqtt
				respentity = handler(topic, msgdata)

				break
			}
			if respentity != nil {
				//Publish resp to server
				resp_topic := topic + "/resp"
				a.SendMqttResp(resp_topic, respentity)
			}

		}); token.Wait() && token.Error() != nil {
		global.Logger.Sugar().With("mqtt subscribe error", token.Error()).Error("Can't subscribe")
		// log.WithField("error", token.Error()).Error("Can't subscribe")
		os.Exit(1)
	}
}

//默认消息处理方法
func DefaultMessageHandle() HandlerFunc {
	return func(topic string, data *Entity) *myRespEntity {
		return &myRespEntity{
			Code: "0",
			Data: data.Data,
		}
	}

}

// func (a *Agent) RegisterMqttHandler(topic string, handler HandlerFunc) {
// 	if token := a.Client.Subscribe(topic, conf.Qos,
// 		func(c MQTT.Client, msg MQTT.Message) {
// 			var msgmqtt infra.Entity
// 			var respentity *infra.RespEntity
// 			var ic infra.IComponent
// 			var msgdate infra.Message
// 			//unmarshal the bgp json, execute action and return status
// 			for true {
// 				LogOutResp(topic, msg.Payload())
// 				if err := json.Unmarshal(msg.Payload(), &msgmqtt); err != nil {
// 					log.Error("RegisterMqttHandler Unmarshal fail")
// 					break
// 				}
// 				//Syslog.LogOutObj(msg.Topic(), msgBgp, Syslog.RESP_OK)
// 				msgdate.Entity = &msgmqtt
// 				respentity = handler(ic, topic, &msgdate)

// 				break
// 			}
// 			if respentity != nil {
// 				//Publish resp to server
// 				resp_topic := topic + "/resp"
// 				a.SendMqttResp(resp_topic, respentity)
// 			}

// 		}); token.Wait() && token.Error() != nil {
// 		log.WithField("error", token.Error()).Error("Can't subscribe")
// 		os.Exit(1)
// 	}

// }

// Connect opens a new connection
func (a *Agent) Connect() (err error) {
	token := a.Client.Connect()
	if token.WaitTimeout(2*time.Second) == false {
		return errors.New("Open timeout")
	}
	if token.Error() != nil {
		return token.Error()
	}

	go func() {
		done := make(chan os.Signal)
		signal.Notify(done, os.Interrupt)
		<-done
		// log.Info("Shutting down agent")
		global.Logger.Info("Shutting down agent")
		a.Close()
	}()

	return
}

// Close agent
func (a *Agent) Close() {
	a.Client.Disconnect(250)
}

// 响应结构体，测试用，实际看需求自定义
type RespEntity struct {
	Rid       string
	Timestamp uint64
	Object    interface{}
	Action    interface{}
	Code      string
	Data      interface{}
}

// 同上，测试响应结构体
type myRespEntity struct {
	Time string
	Code string
	Data interface{}
}

// 返回响应
func (a *Agent) SendMqttResp(topic string, msg *myRespEntity) {
	// msgResp := &RespEntity{
	// 	Rid:       msg.Rid,
	// 	Timestamp: uint64(time.Now().Unix()),
	// 	Object:    msg.Object,
	// 	Action:    msg.Action,
	// 	Code:      msg.Code,
	// 	Data:      msg.Data,
	// }
	msgResp := &myRespEntity{
		Time: time.Now().Format("2006-01-02 15:04:05"),
		Code: msg.Code,
		Data: msg.Data,
	}
	payload, err := json.Marshal(msgResp)
	zap.S().Infof("mqtt 应答主题: %s ", topic)
	zap.S().Infof("mqtt 应答数据: %s ", string(payload))
	if err != nil {
		// Logger.With(
		// 	zap.Field{
		// 		Key:    "topic",
		// 		String: topic,
		// 	},
		// 	zap.Field{
		// 		Key:    "rid",
		// 		String: msg.Rid,
		// 	},
		// 	zap.Field{
		// 		Key:       "obj",
		// 		Interface: msg.Object,
		// 	}, zap.Field{
		// 		Key:       "act",
		// 		Interface: msg.Action,
		// 	}, zap.Field{
		// 		Key:    "code",
		// 		String: msg.Code,
		// 	}, zap.Field{
		// 		Key:       "date",
		// 		Interface: msg.Data,
		// 	},
		// ).Info("SendMqttResp: Marshal failed")

		return
	}
	a.Publish(topic, true, payload)
}

// Publish things
func (a *Agent) Publish(topic string, retain bool, payload []byte) (err error) {
	token := a.Client.Publish(topic, conf.Qos, retain, payload)
	if token.WaitTimeout(2*time.Second) == false {
		return errors.New("Publish timout")
	}
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

// 记录发送日志
func LogPush(topic string, msg []byte) {
	msg_val := fmt.Sprintf("Topic: %s Payload: %s", topic, msg)
	// log.WithField(msg_resp, msg_val).Info(topic)
	// fmt.Println("LO", msg_val)
	zap.S().Infof("mqtt 发送数据: %s, %s", msg_val, topic)
}
