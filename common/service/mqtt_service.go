package service

import (
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

var mqttClient mqtt.Client
var topic string

func GetMqttClient(broker_, username_, password_ string, defaultHandler mqtt.MessageHandler) mqtt.Client {
	var (
		broker   = broker_
		clientID = uuid.NewV4().String()
		username = username_
		password = password_
	)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetDefaultPublishHandler(defaultHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		logx.Error("mqttServerStart fail:", token.Error())
		return nil
	}

	return mqttClient
}

func MqttSubscribe(broker_, username_, password_, topic_ string, qos_ byte, defaultHandler mqtt.MessageHandler) error {

	topic = fmt.Sprintf("%s", topic_)
	//logx.Info("mqtt-topic:", topic)

	if GetMqttClient(broker_, username_, password_, defaultHandler) == nil {
		return errors.New("mqtt连接失败")
	}

	token := mqttClient.Subscribe(topic, qos_, nil)
	token.Wait()

	return nil
}

func MqttSend(topic string, payload interface{}, qos_ byte) error {
	if mqttClient == nil {
		return errors.New("mqtt连接失败")
	}
	token := mqttClient.Publish(topic, qos_, false, payload)
	token.Wait()
	return nil
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logx.Info(`Mqtt Connected & topic:`, topic)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
