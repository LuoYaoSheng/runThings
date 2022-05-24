# runThings
## IoT中转队列
### 只接收 QueueName: runThings-logs
### 工厂模式接收  RabbitMQ消息
### 将接收到设备消息存储到 influxdb
### 如果设备消息不等于 EqStatusNor 则通过 MQTT 发送异常事件

#### 生成dockerfile
##### goctl docker -go runThings.go
#### docker生成
##### 进入到 go-things 主目录
##### docker build -t luoyaosheng/runthings:1.1 -f app/common/core/cmd/runThings/Dockerfile .
##### docker build -t luoyaosheng/runthings:latest -f app/common/core/cmd/runThings/Dockerfile .