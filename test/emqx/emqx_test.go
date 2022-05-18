package emqx

import (
	"runThings/common/config"
	"runThings/common/model"
	"runThings/common/service"
	"testing"
	"time"
)

func TestEmqx(t *testing.T) {

	emqxCfg := &config.EmqxConf{
		Url:  "http://127.0.0.1:8081",
		User: "admin",
		Pass: "public",
	}

	m := make(map[string]interface{})
	m["name"] = "智能井盖"
	m["location"] = "智慧展厅"
	m["time"] = time.Now()

	log := model.Eq2MqLog{
		Sn:      "eq001",
		Product: "p001",
		Status:  config.EqStatusAlarm,
		Content: m,
		Title:   "发生倾斜",
	}

	params := &service.EmqxParamsConf{
		App:      "runThings",
		User:     1,
		Project:  1,
		Eq2MqLog: log,
	}

	service.EmqxApiPublish(emqxCfg, params)
}
