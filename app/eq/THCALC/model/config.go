package model

import "time"

type ThCalcModel struct {
	Temperature float64   `json:"temperature"` // 温度
	Humidity    float64   `json:"humidity"`    // 湿度
	On          bool      `json:"on"`          // 开关
	CurTime     time.Time `json:"curTime"`     // 当前时间
}
