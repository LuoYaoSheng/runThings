package redis

import (
	"fmt"
	"runThings/common/service"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

func testReviceMsg(str string) {
	fmt.Println(str, "解散通知~~~")
}

func TestRedis(t *testing.T) {

	// 初始化 redis 客户端
	service.GetRedisClient("127.0.0.1:6379", "123456", 0)

	// 设置过期
	expiration := time.Duration(1) * time.Second
	err := service.SetRdValueTimeout("QQ群", "925653309", expiration)
	if err != nil {
		logx.Error(err)
		return
	}

	// 过期订阅
	err1 := service.SubscribeKeyExpired(testReviceMsg)
	if err1 != nil {
		logx.Error(err1)
		return
	}

	time.Sleep(expiration + 1*time.Second)
}
