package rpc

// 模拟 RPC 进行接收

import (
	model2 "runThings/app/eq/THCALC/model"

	"github.com/LuoYaoSheng/runThingsCommon/common/config"
	"github.com/LuoYaoSheng/runThingsCommon/common/model"
	"github.com/LuoYaoSheng/runThingsCommon/common/service"

	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/apimachinery/pkg/util/json"
)

var rabbitmqHeartbeat = service.NewRabbitMQSimple("runThings-heartbeat", "amqp://admin:admin@127.0.0.1:5672/")
var rabbitmqLog = service.NewRabbitMQSimple("runThings-logs", "amqp://admin:admin@127.0.0.1:5672/")
var redisClient = service.GetRedisClient("127.0.0.1:6379", "123456", 0) // 直接初始化

func Revive(topic, payload string) {

	logx.Info("rpc:", topic, payload)
	// 需要区分主题： 上报数据[update] / 心跳[heart] / 指令下发[cmd] / 指令应答[ack]
	topics := strings.Split(topic, "/")
	if len(topics) != 4 && topics[0] != "th-calc" {
		return // 过滤非标准
	}

	m := map[string]float64{
		"temperature": model2.TemperatureToplimit,
		"humidity":    model2.HumidityToplimit,
	}
	value, err := service.GetRdValue(topics[2] + "_m")
	if err == nil {
		// 获取值
		err2 := json.Unmarshal([]byte(value), &m)
		if err2 != nil {
			logx.Error(err2)
		}
	}
	logx.Info("rpc-m: ", m)

	// 非下发指令，发送心跳
	if topics[3] != "cmd" {
		heart := model.Eq2MqHeartbeat{
			Sn:        topics[2],
			Product:   topics[1],
			Heartbeat: int64(model2.Heart),
		}
		msg, _ := json.Marshal(heart)
		rabbitmqHeartbeat.PublishSimple(string(msg))
	}

	// 非心跳包，上传数据到日志
	if topics[3] != "heart" && topics[3] != "cmd" && topics[3] != "ack" {
		var tempMap map[string]interface{}

		err = json.Unmarshal([]byte(payload), &tempMap)
		if err != nil {
			logx.Error(err)
			return
		}

		status := config.EqStatusUnknown
		title := ""
		switch topics[3] {
		case "ack":
			status = config.EqStatusAck
		case "cmd":
			status = config.EqStatusCmd // 不用回传
		case "update":
			{
				status = config.EqStatusNor
				temperature := tempMap["temperature"].(float64)
				if temperature > model2.TemperatureToplimit {
					title = "温度高于上限"
					status = config.EqStatusAlarm
				}

				humidity := tempMap["humidity"].(float64)
				if humidity > model2.HumidityToplimit {
					if len(title) > 0 {
						title = title + "|湿度高于上限"
					} else {
						title = "湿度高于上限"
					}
					status = config.EqStatusAlarm
				}
			}
		}

		// 其他情况，发送日志
		log := &model.Eq2MqLog{
			Sn:       topics[2],
			Product:  topics[1],
			Protocol: config.ProtocolMQTT,
			Status:   int64(status),
			Content:  tempMap,
			Title:    title,
			Link:     false,
		}

		msg, _ := json.Marshal(log)
		rabbitmqLog.PublishSimple(string(msg))
	}

}
