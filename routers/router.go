package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zhouqiaokeji/server/middleware"
	"github.com/zhouqiaokeji/server/pkg/conf"
	"github.com/zhouqiaokeji/server/pkg/util"
	"github.com/zhouqiaokeji/server/routers/controllers"
)

// InitRouter Router Initialize
func InitRouter() *gin.Engine {
	if conf.SystemConfig.Mode == conf.MODE_MASTER {
		util.Log().Info("Current Mode ：Master")
		return InitMasterRouter()
	}
	util.Log().Info("Current Mode ：Slave")
	return InitSlaveRouter()

}

// InitSlaveRouter Init Slave Router
func InitSlaveRouter() *gin.Engine {
	arr := gin.Default()
	initCORS(arr)
	//app.POST("/sync", controllers.syncCluster)
	//app.GET("/status", controllers.statusCluster)
	//app.GET("/sync/check", controllers.syncFromCluster)
	return arr
}

// InitMasterRouter Init Master Router
func InitMasterRouter() *gin.Engine {
	app := gin.Default()
	initCORS(app)
	app.POST("/login", controllers.UserLogin)
	app.POST("/check/server", controllers.CreateAppUseInfo)
	holiday := app.Group("/holiday")
	holiday.GET("/get", controllers.GetHolidays)
	// 授权路由
	license := app.Group("/license")
	license.POST("/create", controllers.CreateLicense)
	license.POST("/verify", controllers.VerifyLicense)
	// 添加JWT验证
	app.Use(middleware.CurrentUser())
	license.GET("/list", controllers.GetLicenses)
	license.GET("/getInfo", controllers.GetLicense)
	license.POST("/status", controllers.ChangeStatus)
	license.POST("/bind", controllers.BindLicense)
	license.GET("/remove", controllers.RemoveLicense)
	//holiday.POST("/list", controllers.listHolidays)
	holiday.POST("/create", controllers.CreateHoliday)
	appInfo := app.Group("/appInfo")
	appInfo.POST("/list", controllers.GetAppInfos)
	return app
}

// initCORS CORS Init
func initCORS(router *gin.Engine) {
	if conf.CORSConfig.AllowOrigins[0] != "UNSET" {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     conf.CORSConfig.AllowOrigins,
			AllowMethods:     conf.CORSConfig.AllowMethods,
			AllowHeaders:     conf.CORSConfig.AllowHeaders,
			AllowCredentials: conf.CORSConfig.AllowCredentials,
			ExposeHeaders:    conf.CORSConfig.ExposeHeaders,
		}))
		return
	}

	// slave模式下未启动跨域的警告
	if conf.SystemConfig.Mode == "slave" {
		util.Log().Warning("当前作为存储端（Slave）运行，但未启用跨域配置，可能会导致 Master 端无法正常使用")
	}
}
