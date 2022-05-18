package emqx

import (
	"encoding/json"
	"fmt"
	"runThings/common/config"
	"runThings/common/model"
	"runThings/common/service"
	"strconv"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

var MessagePubHandler mqtt.MessageHandler = func(mqttClient mqtt.Client, msg mqtt.Message) {
	fmt.Printf("+++++++++++Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func TestMqtt(t *testing.T) {

	mqttCfg := &config.MqttConf{
		Broker:   "ws://127.0.0.1:8083/mqtt",
		Username: "",
		Password: "",
		Topic:    "log/#",
		Qos:      0,
	}

	// 开启订阅模式
	err := service.MqttSubscribe(mqttCfg.Broker, mqttCfg.Username, mqttCfg.Password, mqttCfg.Topic, mqttCfg.Qos, MessagePubHandler)
	if err != nil {
		logx.Error(err)
		return
	}

	// 发送一条 日志
	m := make(map[string]interface{})
	m["name"] = "智能井盖"
	m["location"] = "智慧展厅"
	m["time"] = time.Now()

	log := model.Eq2MqLog{
		Sn:      "eq001",
		Product: "p1001",
		Status:  config.EqStatusAlarm,
		Content: m,
		Title:   "发生倾斜",
	}

	topic := "log/" + log.Product + "/" + log.Sn + "/" + strconv.FormatInt(log.Status, 10)
	content, _ := json.Marshal(log)

	err2 := service.MqttSend(topic, content, mqttCfg.Qos)
	if err2 != nil {
		logx.Error(err2)
		return
	}
}
