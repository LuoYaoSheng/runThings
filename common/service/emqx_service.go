package service

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"runThings/common/config"
	"runThings/common/model"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmqxParamsConf struct {
	model.Eq2MqLog
	App     string // 暂时未使用，默认 runThings
	User    int64  // 用户id
	Project int64  // 项目id
}

// EmqxApiPublish 使用emqx自带api服务，便于业务管理后台操作
func EmqxApiPublish(emqxCfg *config.EmqxConf, params *EmqxParamsConf) {

	if len(emqxCfg.Url) == 0 {
		logx.Error("未配置Emqx")
		return
	}

	content, _ := json.Marshal(params.Content)

	payload := base64.StdEncoding.EncodeToString(content)
	objUrl := emqxCfg.Url + "/api/v4/mqtt/publish"

	// 订阅规则
	topic := "app/" + params.App + "/" + strconv.FormatInt(params.User, 10) + "/" + strconv.FormatInt(params.Project, 10) + "/" + params.Product + "/" + params.Sn + "/" + strconv.FormatInt(params.Status, 10)
	data := `{"topic":"` + topic + `","payload":"` + payload + `","qos":1,"encoding":"base64"}`

	req, err := http.NewRequest("POST", objUrl, strings.NewReader(data))
	if err != nil {
		logx.Error(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(emqxCfg.User, emqxCfg.Pass)
	_, err1 := http.DefaultClient.Do(req)
	if err1 != nil {
		logx.Error(err1)
		return
	}
}
