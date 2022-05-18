package rabbitmq

import (
	"fmt"
	"runThings/common/service"
	"strconv"
	"testing"
	"time"
)

// 测试发送实时日志
func TestSendSimple(t *testing.T) {
	// 发送多个，便于测试工厂模式
	rabbitmq := service.NewRabbitMQSimple("runThings-logs", "amqp://admin:admin@127.0.0.1:5672/")
	for i := 0; i < 10; i++ {
		// 发送 log
		msg := `{"sn":"1034639560","product":"p001","protocol":0,"status":` + strconv.Itoa(i) + `,"content":{"name":"runThings","value":1231},"title":"万物互联从此开始","link":false}`
		rabbitmq.PublishSimple(msg)
	}
}

// 测试发送心跳包
func TestSendSimple2(t *testing.T) {
	// 发送 设备上下线
	rabbitmq := service.NewRabbitMQSimple("runThings-heartbeat", "amqp://admin:admin@127.0.0.1:5672/")

	heartbeat := 10
	msg := `{"sn":"1034639560","product":"p001", "heartbeat":` + strconv.Itoa(heartbeat) + `}`

	// 6 * 10 等于 1分钟
	for i := 0; i < 6; i++ {
		rabbitmq.PublishSimple(msg)
		// 模拟心跳间隔
		fmt.Println("发送心跳", msg)
		time.Sleep(time.Duration(heartbeat) * time.Second)
	}
}
