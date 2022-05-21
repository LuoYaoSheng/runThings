package hub

// 模拟一个业务中心
import (
	"fmt"
	model2 "runThings/app/eq/THCALC/model"

	"github.com/LuoYaoSheng/runThingsCommon/common/config"
	"github.com/LuoYaoSheng/runThingsCommon/common/model"
	"github.com/LuoYaoSheng/runThingsCommon/common/service"

	"strconv"
	"strings"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/apimachinery/pkg/util/json"
)

var MessagePubHandler mqtt.MessageHandler = func(mqttClient mqtt.Client, msg mqtt.Message) {
	fmt.Printf("+++++++++++Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	topics := strings.Split(msg.Topic(), "/")
	if len(topics) != 5 {
		logx.Error("过滤非标准")
		return
	}

	sn := topics[3]
	productKey := topics[2]
	status, _ := strconv.Atoi(topics[4])

	log := &model.Eq2MqLog{}
	err := json.Unmarshal(msg.Payload(), &log)
	if err != nil {
		logx.Error(err)
		return
	}
	content, _ := json.Marshal(log.Content)

	// 存储表 -- 需要考虑异常事件重复问题: 可能又要用到Redis，压力挺大
	InsertMySQL(sn, productKey, log.Title, string(content), status)

	// 业务端发送告警到前端

	msf := map[string]interface{}{}
	err = json.Unmarshal(content, &msf)
	if err != nil {
		logx.Error(err)
		return
	}

	// 模拟更改阈值
	if status == config.EqStatusAlarm {

		m := map[string]float64{
			"temperature": model2.TemperatureToplimit,
			"humidity":    model2.HumidityToplimit,
		}
		value, err := service.GetRdValue(sn + "_m")
		if err == nil {
			// 获取值
			err2 := json.Unmarshal([]byte(value), &m)
			if err2 != nil {
				logx.Error(err2)
			}
		}

		// 通过 告警时，自增 +1 进行修改
		content := map[string]interface{}{
			"temperature": m["temperature"] + 1,
			"humidity":    m["humidity"] + 1,
		}
		//dataType, _ := json.Marshal(content)
		//err = service.SetRdValue(sn+"_m", dataType)
		//if err != nil {
		//	logx.Error(err)
		//	return
		//}
		// 直接修改Redis，修改成通过发送到 runTings 操作

		threshold := model.Eq2MqThreshold{
			Sn:      sn,
			Content: content,
		}
		dataType, _ := json.Marshal(threshold)
		thresholdMQ(string(dataType))
	}

	if status == config.EqStatusAlarm && msf["temperature"].(float64) > 50.0 {
		// 模拟下发指令

		cmdContent := map[string]interface{}{}
		cmdContent["on"] = !msf["on"].(bool)

		cmd := model.Eq2MqCmd{
			Sn:      sn,
			Content: cmdContent,
		}

		cmdData, _ := json.Marshal(cmd)

		logx.Info("---- 下发指令 ----", string(cmdData))
		cmdMQ(string(cmdData))
	}
}

func InsertMySQL(sn, productKey, title, content string, status int) {
	sql := "insert into eq_log (sn,product_key, status,title, content,create_time)values (?,?,?,?,?,?)"
	value := [6]interface{}{sn, productKey, status, title, content, time.Now()}

	//执行SQL语句
	_, err := db.Exec(sql, value[0], value[1], value[2], value[3], value[4], value[5])
	if err != nil {
		logx.Error("exec failed,", err)
	}
}

func cmdMQ(content string) {
	rabbitmqCmd.PublishSimple(content)
}
func thresholdMQ(content string) {
	rabbitmqThreshold.PublishSimple(content)
}

var (
	rabbitmqCmd       *service.RabbitMQ
	rabbitmqThreshold *service.RabbitMQ
	db                *sqlx.DB
)

func TestHub(t *testing.T) {

	// 订阅 MQTT ，获取设备异常情况
	topic := "event/runTings/" + model2.ProductKey + "/#"
	mqttCfg := &config.MqttConf{
		Broker: "ws://127.0.0.1:8083/mqtt",
		Topic:  topic,
	}

	// 开启订阅模式
	err := service.MqttSubscribe(mqttCfg.Broker, mqttCfg.Username, mqttCfg.Password, mqttCfg.Topic, mqttCfg.Qos, MessagePubHandler)
	if err != nil {
		logx.Error(err)
		return
	}

	// 获取 Redis
	redisCfg := &config.RedisConf{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	}
	service.GetRedisClient(redisCfg.Addr, redisCfg.Password, redisCfg.DB)

	// 获取 rabbitmq
	rabbitmqCmd = service.NewRabbitMQSimple("runThings-cmd-"+model2.ProductKey, "amqp://admin:admin@127.0.0.1:5672/")
	rabbitmqThreshold = service.NewRabbitMQSimple("runThings-threshold", "amqp://admin:admin@127.0.0.1:5672/")

	// 获取 mysql
	database, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/eq")
	if err != nil {
		logx.Error("open mysql failed,", err)
	}
	db = database

	// 获取告警规则

	select {}
}
