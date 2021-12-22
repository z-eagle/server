package boot

import (
	"github.com/gin-gonic/gin"
	model "github.com/zhouqiaokeji/server/models"
	"github.com/zhouqiaokeji/server/pkg/cache"
	"github.com/zhouqiaokeji/server/pkg/conf"
	"github.com/zhouqiaokeji/server/pkg/id"
)

// Init App Initialize
func Init(path string) {
	conf.Init(path)
	// Debug Closed
	if !conf.SystemConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	dependencies := []struct {
		mode    string
		factory func()
	}{
		//{
		//	"both",
		//	func() {
		//		scripts.Init()
		//	},
		//},
		{
			conf.MODE_BOTH,
			func() {
				id.Init()
			},
		},
		{
			conf.MODE_BOTH,
			func() {
				cache.Init()
			},
		},
		{
			conf.MODE_MASTER,
			func() {
				model.Init()
			},
		},
		//{
		//	"both",
		//	func() {
		//		task.Init()
		//	},
		//},
		//{
		//	"master",
		//	func() {
		//		cluster.Init()
		//	},
		//},
		//{
		//	"master",
		//	func() {
		//		aria2.Init(false, cluster.Default, mq.GlobalMQ)
		//	},
		//},
		//{
		//	"master",
		//	func() {
		//		email.Init()
		//	},
		//},
		//{
		//	"master",
		//	func() {
		//		crontab.Init()
		//	},
		//},
		//{
		//	"slave",
		//	func() {
		//		cluster.InitController()
		//	},
		//},
		//{
		//	"both",
		//	func() {
		//		auth.Init()
		//	},
		//},
	}

	for _, dependency := range dependencies {
		switch dependency.mode {
		case conf.MODE_MASTER:
			if conf.SystemConfig.Mode == conf.MODE_MASTER {
				dependency.factory()
			}
		case conf.MODE_SLAVE:
			if conf.SystemConfig.Mode == conf.MODE_SLAVE {
				dependency.factory()
			}
		default:
			dependency.factory()
		}
	}
}
