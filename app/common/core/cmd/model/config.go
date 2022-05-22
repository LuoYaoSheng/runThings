package model

// RunThingsConf 结构体
type RunThingsConf struct {
	Mqurl     string
	Logs      string // 日志
	Heartbeat string // 心跳
	Cmd       string // 命令
	Threshold string // 阈值
}

type RuleContent struct {
	Property  string      `json:"property"`
	Condition int         `json:"condition"`
	Value     interface{} `json:"value"`
}

type Rule struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	Code    string `json:"code"`
	Sn      string `json:"sn"`
	Content string `json:"content"` // 暂不适用 RuleContent
}
