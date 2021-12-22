package main

import (
	"flag"
	"github.com/zhouqiaokeji/server/boot"
	"github.com/zhouqiaokeji/server/pkg/conf"
	"github.com/zhouqiaokeji/server/pkg/util"
	"github.com/zhouqiaokeji/server/routers"
)

var confPath string

func init() {
	flag.StringVar(&confPath, "c", util.RelativePath(conf.CONF_FILE_NAME), "config file path")
	boot.Init(confPath)
}

func main() {
	api := routers.InitRouter()

	// Enable SSL
	if conf.SSLConfig.CertPath != "" {
		go func() {
			util.Log().Info("Listen On %s", conf.SSLConfig.Listen)
			if err := api.RunTLS(conf.SSLConfig.Listen,
				conf.SSLConfig.CertPath, conf.SSLConfig.KeyPath); err != nil {
				util.Log().Error("Listen Fail [%s]，%s", conf.SSLConfig.Listen, err)
			}
		}()
	}

	// Enable Unix
	if conf.UnixConfig.Listen != "" {
		util.Log().Info("Listen On %s", conf.UnixConfig.Listen)
		if err := api.RunUnix(conf.UnixConfig.Listen); err != nil {
			util.Log().Error("Listen Fail [%s]，%s", conf.UnixConfig.Listen, err)
		}
		return
	}

	util.Log().Info("Listen On %s", conf.SystemConfig.Listen)
	if err := api.Run(conf.SystemConfig.Listen); err != nil {
		util.Log().Error("Listen Fail [%s]，%s", conf.SystemConfig.Listen, err)
	}
}
