# runThings

#### 介绍
runThings：一个小型IoT管理中间件，可多开。

计划： 
1、接入tdengine
2、接入集群
3、...

#### 软件架构
![readme](https://gitee.com/luoyaosheng/run-things/raw/master/doc/readme.png)

#### 使用教程[docker版本]

1. 进入项目内，启动 docker-compose
```
docker-compose up -d
```
2. 使用 [demo-gitee](https://gitee.com/luoyaosheng/run-things-demo) / [demo-github](https://github.com/LuoYaoSheng/runThingsDemo)，测试 runThings 

#### 本地教程

1. 进入项目内，创建 mod 并 更新库
```
go mod init runThings
go mod tidy
```
2. 启动自建 mqtt、RabbitMQ、influxdb1.8 、Redis 环境，或通过 docker-compose 快速配置环境
```
docker-compose up -d -f docker-compose-env.yml
```
4. 修改配置文件 run-things/app/common/core/cmd/runThings/etc
5. 使用 [demo-gitee](https://gitee.com/luoyaosheng/run-things-demo) / [demo-github](https://github.com/LuoYaoSheng/runThingsDemo)，测试 runThings

#### 常见错误
1、网络启动失败
```
ERROR: Pool overlaps with other one on this address space

修改 subnet 即可
```
