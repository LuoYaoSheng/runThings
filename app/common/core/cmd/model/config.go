package model

// RunThingsConf 结构体
type RunThingsConf struct {
	Mqurl     string
	Logs      string // 日志
	Heartbeat string // 心跳
	Cmd       string // 命令
	Threshold string // 阈值
}
