package rabbitmq

import (
	"fmt"
	"runThings/common/service"
	"testing"
)

func recieveSimple(str string) {
	fmt.Println("---简单模式接收到数据: ", str)
}

func TestReceiveSimple(t *testing.T) {
	rabbitmq := service.NewRabbitMQSimple("runThings", "amqp://admin:admin@127.0.0.1:5672/")
	rabbitmq.ConsumeSimple(recieveSimple)
}

// 接收
func TestReceiveSimpleCmd(t *testing.T) {
	product := "p001"
	rabbitmq := service.NewRabbitMQSimple("runThings-cmd-"+product, "amqp://admin:admin@127.0.0.1:5672/")
	rabbitmq.ConsumeSimple(recieveSimple)
}
