package influxdb

import (
	"fmt"
	"runThings/common/config"
	"runThings/common/service"
	"testing"
)

func TestInfluxdb(t *testing.T) {
	influxdbCfg := &config.InfluxdbConf{
		Addr:      "http://127.0.0.1:8086",
		Username:  "root",
		Password:  "root",
		Database:  "runThings",
		Precision: "",
		Prefix:    "test_",
	}

	sn := "925653309"
	imei := ""
	payload := make(map[string]interface{})

	payload["qq"] = "1034639560"
	payload["author"] = "寺西"

	service.GetClient(influxdbCfg.Addr, influxdbCfg.Username, influxdbCfg.Password, influxdbCfg.Database, influxdbCfg.Precision)

	_, err := service.WirteInflux(sn, imei, config.EqStatusNor, payload, influxdbCfg.Database, influxdbCfg.Prefix, influxdbCfg.Precision)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功")

	list, err := service.SelectInflux(sn, "", influxdbCfg.Database, influxdbCfg.Prefix)
	if err != nil {
		return
	}
	fmt.Println("读取内容:", list)
}
