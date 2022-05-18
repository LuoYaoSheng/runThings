package api

// 模拟 api
import (
	"runThings/app/eq/THCALC/cmd/rpc"
	"runThings/app/eq/THCALC/model"
	"runThings/common/config"
	"runThings/common/service"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

var apiMessagePubHandler mqtt.MessageHandler = func(mqttClient mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("apiMessagePubHandler: %s from topic: %s\n", msg.Payload(), msg.Topic())

	// 将数据传递给 rpc
	rpc.Revive(msg.Topic(), string(msg.Payload()))
}

func TestApi(t *testing.T) {

	pkey := model.ProductKey

	// MQTT 订阅，设备端
	mqttCfg := &config.MqttConf{
		Broker:   "ws://127.0.0.1:8083/mqtt",
		Username: "",
		Password: "",
		Topic:    `th-calc/` + pkey + `/#`, // 订阅该产品下所有设备
		Qos:      0,
	}

	// 开启订阅模式
	err := service.MqttSubscribe(mqttCfg.Broker, mqttCfg.Username, mqttCfg.Password, mqttCfg.Topic, mqttCfg.Qos, apiMessagePubHandler)
	if err != nil {
		logx.Error(err)
		return
	}

	select {} // 强制等待
}
