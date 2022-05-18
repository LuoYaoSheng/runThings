package config

//设备状态
const (
	EqStatusNor       = iota // 正常
	EqStatusAck              // 应答
	EqStatusGet              // 查询
	EqStatusOperate          // 设备操作
	EqStatusOnline           // 上线
	EqStatusOffline          // 下线
	EqStatusAbnormal         // 异常
	EqStatusAlarm            // 告警
	EqStatusCmd              // 下发指令
	EqStatusNotActive = 50   // 未激活
	EqStatusDelete    = 51   // 删除
	EqStatusEmpty     = 98   // 空/无需操作
	EqStatusUnknown   = 99   // 未知
)

// 传输协议
const (
	ProtocolUnknown = iota // 未知
	ProtocolHTTP
	ProtocolCoAP
	ProtocolMQTT
	ProtocolDDS
	ProtocolAMQP
	ProtocolXMPP
	ProtocolJMS
)
