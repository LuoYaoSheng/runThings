package main

import (
	"flag"
	"log"
	"runThings/app/common/core/cmd/runThings/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/runThings.yaml", "Specify the config file")

func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	flag.Parse()
	conf.MustLoad(*configFile, &svc.Conf)
	svc.NewServiceContext()
	select {} // 使程序一直启动着
}
