package service

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	//连接
	conn *amqp.Connection
	//管道
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key Simple模式 几乎用不到
	Key string
	//连接信息
	Mqurl string
}

// 向外部提供方法 --- 临时写的，有用到才添加
type recieveMsgFunc func(string)

// NewRabbitMQ 创建RabbitMQ结构体实例
func NewRabbitMQ(queuename string, exchange string, key string, mqurl string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queuename, Exchange: exchange, Key: key, Mqurl: mqurl}
	var err error
	//创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误！")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

// Destory 断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple 简单模式step：1。创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string, mqurl string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "", mqurl)
}

// NewRabbitMQPubSub 订阅模式创建rabbitmq实例
func NewRabbitMQPubSub(exchangeName string, mqurl string) *RabbitMQ {
	//创建rabbitmq实例
	rabbitmq := NewRabbitMQ("", exchangeName, "", mqurl)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connecct rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel!")
	return rabbitmq
}

// PublishPub 订阅模式生成
func (r *RabbitMQ) PublishPub(message string) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		"fanout",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an excha"+"nge")

	//2 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: []byte(message),
		})
}

// RecieveSub 订阅模式消费端代码
func (r *RabbitMQ) RecieveSub(fc recieveMsgFunc) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		"fanout",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an excha"+"nge")
	//2试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range message {
			log.Printf("Received a message:%s,", d.Body)
			fc(string(d.Body))
		}
	}()
	fmt.Println("退出请按 Ctrl+C")
	<-forever
}

// NewRabbitMQTopic 话题模式 创建RabbitMQ实例
func NewRabbitMQTopic(exchagne string, routingKey string, mqurl string) *RabbitMQ {
	//创建rabbitmq实例
	rabbitmq := NewRabbitMQ("", exchagne, routingKey, mqurl)
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed     to connect rabbingmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// PublishTopic 话题模式发送信息
func (r *RabbitMQ) PublishTopic(message string) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 话题模式
		"topic",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "topic failed to declare an excha"+"nge")
	//2发送信息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: []byte(message),
		})
}

// RecieveTopic 话题模式接收信息
//要注意key
//其中* 用于匹配一个单词，#用于匹配多个单词（可以是零个）
//匹配 表示匹配imooc.* 表示匹配imooc.hello,但是imooc.hello.one需要用imooc.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 话题模式
		"topic",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an excha"+"nge")
	//2试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range message {
			log.Printf("Received a message:%s,", d.Body)
		}
	}()
	fmt.Println("退出请按 Ctrl+C")
	<-forever
}

// NewRabbitMQRouting 路由模式 创建RabbitMQ实例
func NewRabbitMQRouting(exchagne string, routingKey string, mqurl string) *RabbitMQ {
	//创建rabbitmq实例
	rabbitmq := NewRabbitMQ("", exchagne, routingKey, mqurl)
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed     to connect rabbingmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// PublishRouting 路由模式发送信息
func (r *RabbitMQ) PublishRouting(message string) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		"direct",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an excha"+"nge")
	//发送信息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: []byte(message),
		})
}

// RecieveRouting 路由模式接收信息
func (r *RabbitMQ) RecieveRouting() {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		//交换机名称
		r.Exchange,
		//交换机类型 广播类型
		"direct",
		//是否持久化
		true,
		//是否字段删除
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an excha"+"nge")
	//2试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range message {
			log.Printf("Received a message:%s,", d.Body)
		}
	}()
	fmt.Println("退出请按 Ctrl+C")
	<-forever
}

// PublishSimple 简单模式Step:2、简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message string) {
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,
		//是否具有排他性 true表示自己可见 其他用户不能访问
		false,
		//是否阻塞 true表示要等待服务器的响应
		false,
		//额外数学系
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//2.发送消息到队列中
	r.channel.Publish(
		//默认的Exchange交换机是default,类型是direct直接类型
		r.Exchange,
		//要赋值的队列名称
		r.QueueName,
		//如果为true，根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息还给发送者
		false,
		//消息
		amqp.Publishing{
			//类型
			ContentType: "text/plain",
			//消息
			Body: []byte(message),
		})
}

// ConsumeSimple 修改简单模式，其实支持外部函数调用
func (r *RabbitMQ) ConsumeSimple(fc recieveMsgFunc) {
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		//队列名称
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外数学系
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true,表示不能同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)

	//启用协程处理
	go func() {
		for d := range msgs {
			//实现我们要处理的逻辑函数
			log.Printf("Received a message:%s", d.Body)
			//fmt.Println(d.Body)
			fc(string(d.Body))
		}
	}()

	log.Printf("【*】warting for messages, To exit press CCTRAL+C")
	<-forever
}
