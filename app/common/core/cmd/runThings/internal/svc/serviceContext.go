package svc

import (
	"encoding/json"
	"fmt"
	"runThings/app/common/core/cmd/runThings/internal/config"
	config2 "runThings/common/config"
	"runThings/common/model"
	"runThings/common/service"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	rabbitmqLog       *service.RabbitMQ
	rabbitmqHeartbeat *service.RabbitMQ
	rabbitmqThreshold *service.RabbitMQ

	Conf config.Config
)

func NewServiceContext() {
	// 开启redis订阅模式
	service.GetRedisClient(Conf.Redis.Addr, Conf.Redis.Password, Conf.Redis.DB)
	err := service.SubscribeKeyExpired(receiveRedis)
	if err != nil {
		logx.Error(err)
		return
	}

	// 初始化 influxdb
	service.GetClient(Conf.Influx.Addr, Conf.Influx.Username, Conf.Influx.Password, Conf.Influx.Database, Conf.Influx.Precision)

	// 开启非订阅模式
	err1 := service.MqttSubscribe(Conf.Mqtt.Broker, Conf.Mqtt.Username, Conf.Mqtt.Password, Conf.Mqtt.Topic, Conf.Mqtt.Qos, nil)
	if err1 != nil {
		logx.Error(err1)
	}

	// 订阅日志
	go func() {
		rabbitmqLog = service.NewRabbitMQSimple(Conf.RunThings.Logs, Conf.RunThings.Mqurl)
		rabbitmqLog.ConsumeSimple(receiveLog)
	}()

	// 订阅心跳
	go func() {
		rabbitmqHeartbeat = service.NewRabbitMQSimple(Conf.RunThings.Heartbeat, Conf.RunThings.Mqurl)
		rabbitmqHeartbeat.ConsumeSimple(receiveHeartbeat)
	}()

	// 订阅阈值
	go func() {
		rabbitmqThreshold = service.NewRabbitMQSimple(Conf.RunThings.Threshold, Conf.RunThings.Mqurl)
		rabbitmqThreshold.ConsumeSimple(receiveThreshold)
	}()
}

func receiveLog(str string) {
	var log model.Eq2MqLog
	err := json.Unmarshal([]byte(str), &log)
	if err != nil {
		logx.Error(err)
		return
	}
	//fmt.Println("---接收到数据: ", log.Sn, log.Status, log.Title, log.Content)
	receiveLogModel(&log)
}

func receiveHeartbeat(str string) {
	var hb model.Eq2MqHeartbeat
	err5 := json.Unmarshal([]byte(str), &hb)
	if err5 != nil {
		logx.Error(err5)
		return
	}
	receiveHeartbeatModel(&hb)
}

func receiveThreshold(str string) {

	fmt.Println("------------------")
	fmt.Println(str)
	fmt.Println("------------------")

	var threshold model.Eq2MqThreshold
	err5 := json.Unmarshal([]byte(str), &threshold)
	if err5 != nil {
		logx.Error(err5)
		return
	}
	receiveThresholdModel(&threshold)
}

func receiveRedis(str string) {
	//logx.Info("key:", str, "过期通知")
	v, err := service.GetRdValue(str + "_2")
	if err != nil {
		logx.Error(err)
		return
	}

	// 清除附表key
	_, err = service.DelRdValue(str + "_2")
	if err != nil {
		logx.Error(err)
	}

	// 发送离线通知
	log := model.Eq2MqLog{
		Sn:      str,
		Product: v,
		Status:  config2.EqStatusOffline,
		Title:   "离线",
		Content: map[string]interface{}{
			"Sn": str,
		},
	}
	content, _ := json.Marshal(log)
	//receiveLog(string(content))
	rabbitmqLog.PublishSimple(string(content))
}

func receiveLogModel(log *model.Eq2MqLog) {

	// 写入时序数据库
	_, err2 := service.WirteInflux(log.Sn, log.Product, log.Status, log.Content, Conf.Influx.Database, Conf.Influx.Prefix, Conf.Influx.Precision)
	if err2 != nil {
		fmt.Println(err2)
	}

	// 异常通过mqtt推送
	if log.Status != config2.EqStatusNor {
		topic := Conf.Mqtt.Topic + log.Product + "/" + log.Sn + "/" + strconv.FormatInt(log.Status, 10)
		content, _ := json.Marshal(log)

		err4 := service.MqttSend(topic, content, Conf.Mqtt.Qos)
		if err4 != nil {
			logx.Error(err4)
			return
		}

		// 下发指令进行转发
		if log.Status == config2.EqStatusCmd && len(Conf.RunThings.Cmd) > 0 {
			rabbitmq := service.NewRabbitMQSimple(Conf.RunThings.Cmd+log.Product, Conf.RunThings.Mqurl)
			rabbitmq.PublishSimple(string(content))
		}
	}
}

func receiveHeartbeatModel(hb *model.Eq2MqHeartbeat) {
	// 开始 Redis 操作
	_, err := service.GetRdValue(hb.Sn)
	if err == redis.Nil {
		// 存储设备 & 设备上线
		err = service.SetRdValueTimeout(hb.Sn, hb.Heartbeat, time.Duration(hb.Heartbeat+1)*time.Second)
		if err != nil {
			logx.Error(err)
			return
		}

		err = service.SetRdValue(hb.Sn+"_2", hb.Product)
		if err != nil {
			logx.Error(err)
			return
		}

		// 发送上线通知
		log := model.Eq2MqLog{
			Sn:      hb.Sn,
			Product: hb.Product,
			Status:  config2.EqStatusOnline,
			Title:   "上线",
			Content: map[string]interface{}{
				"Sn": hb.Sn,
			},
		}
		content, _ := json.Marshal(log)
		//receiveLog(string(content))
		rabbitmqLog.PublishSimple(string(content))
	} else if err != nil {
		logx.Error(err)
		return
	} else {
		// 更新过期时间
		err = service.SetRdValueTimeout(hb.Sn, hb.Heartbeat, time.Duration(hb.Heartbeat+1)*time.Second)
		if err != nil {
			logx.Error(err)
			return
		}
	}
}

func receiveThresholdModel(threshold *model.Eq2MqThreshold) {
	dataType, _ := json.Marshal(threshold.Content)
	err := service.SetRdValue(threshold.Sn+"_m", string(dataType))
	if err != nil {
		logx.Error(err)
		return
	}
}
