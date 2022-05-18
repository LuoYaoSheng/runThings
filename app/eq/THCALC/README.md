### 此处模拟 温湿度计 [MQTT 接入]

#### 结构说明
```
- THCALC        // 温湿度计演示
 - cmd          // 产品解析
  - api         // 模拟API网关
  - rpc         // 模拟RPC服务
 - hub          // 业务服务中心
 - model        // 设备数据结构
 - simulation   // 设备模拟
```
#### 温湿度计 功能表
```
1、定时上传 温度、湿度、开关状态
2、可以设置 温度、湿度 阈值
3、可以下发 控制开关
```

#### 操作流程
```
前置条件：
    连接mysql，创建库：eq，将 doc/eq.sql 导入
具体步骤：
1、启动 业务服务中心 hub_test.go
2、启动 产品解析 rpc_test.go / api_test.go
3、启动 设备模拟 simulation_test.go
```

#### 设备topic
```
#上报数据
th-calc/${ProductKey}/${DeviceName}/update  
#心跳
th-calc/${ProductKey}/${DeviceName}/heart
#指令下发
th-calc/${ProductKey}/${DeviceName}/cmd
#指令应答
th-calc/${ProductKey}/${DeviceName}/ack
```

#### 设备心跳
```
headerbeat 10s
```
        