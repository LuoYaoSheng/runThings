package rpc

import (
	model2 "runThings/app/eq/THCALC/model"
	"runThings/common/config"
	"runThings/common/model"
	"runThings/common/service"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/apimachinery/pkg/util/json"
)

func recieveSimple(str string) {
	logx.Info("---rpc_test: ", str)

	cmd := &model.Eq2MqCmd{}
	err := json.Unmarshal([]byte(str), &cmd)
	if err != nil {
		logx.Error(err)
		return
	}

	topic := `th-calc/` + model2.ProductKey + `/` + cmd.Sn + `/cmd`
	service.MqttSend(topic, str, mqttCfg.Qos) // 直接透传 -- 可以的话，去掉 sn 也可以，减少传输内容，毕竟硬件空间比较小
}

var mqttCfg *config.MqttConf

func TestReceiveSimpleCmd(t *testing.T) {

	mqttCfg = &config.MqttConf{
		Broker:   "ws://127.0.0.1:8083/mqtt",
		Username: "",
		Password: "",
		Topic:    "",
		Qos:      0,
	}

	// 开启非订阅模式
	err := service.MqttSubscribe(mqttCfg.Broker, mqttCfg.Username, mqttCfg.Password, mqttCfg.Topic, mqttCfg.Qos, nil)
	if err != nil {
		logx.Error(err)
		return
	}

	product := model2.ProductKey
	rabbitmq := service.NewRabbitMQSimple("runThings-cmd-"+product, "amqp://admin:admin@127.0.0.1:5672/")
	rabbitmq.ConsumeSimple(recieveSimple)
}
