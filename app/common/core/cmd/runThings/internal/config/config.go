package config

import (
	"runThings/app/common/core/cmd/model"

	"github.com/LuoYaoSheng/runThingsCommon/common/config"

	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf
	RunThings model.RunThingsConf
	Redis     config.RedisConf
	Mqtt      config.MqttConf
	Influx    config.InfluxdbConf
}
