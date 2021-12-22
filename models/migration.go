package models

import (
	"github.com/fatih/color"
	"github.com/zhouqiaokeji/server/pkg/cache"
	"github.com/zhouqiaokeji/server/pkg/conf"
	"github.com/zhouqiaokeji/server/pkg/util"
)

// 是否需要迁移
func needMigration() bool {
	var setting Setting
	result := DB.Where("name = ?", "db_version_"+conf.RequiredDBVersion).First(&setting)
	return result.Error != nil
}

//执行数据迁移
func migration() {

	util.Log().Info("开始进行数据库初始化...")

	// 确认是否需要执行迁移
	if !needMigration() {
		util.Log().Info("数据库版本匹配，跳过数据库迁移")
		return
	}

	// 清除所有缓存
	if instance, ok := cache.Store.(*cache.RedisStore); ok {
		_ = instance.DeleteAll()
	}

	// 自动迁移模式
	//if conf.DatabaseConfig.Type == "mysql" {
	//	DB = DB.Set("gorm:table_options", "ENGINE=InnoDB")
	//}

	_ = DB.AutoMigrate(&User{}, &Setting{}, &License{}, &Holidays{}, &AppUseInfo{})

	// 创建初始管理员账户
	initAdminUser()

	// 向设置数据表添加初始设置
	initDefaultSettings()

	// 执行数据库升级脚本
	//execUpgradeScripts()

	util.Log().Info("数据库初始化结束")

}

func initAdminUser() {
	_, exist := GetUserByUserName("admin")
	password := "zqkj123456"

	// 未找到初始用户时，则创建
	if !exist {
		defaultUser := NewUser()
		defaultUser.Nick = "超级管理员"
		defaultUser.UserName = "admin"
		defaultUser.Status = Active
		defaultUser.SetPassword(password)

		if err := DB.Model(&User{}).Create(&defaultUser).Error; err != nil {
			util.Log().Panic("无法创建初始用户, %s", err)
		}

		c := color.New(color.FgWhite).Add(color.BgBlack).Add(color.Bold)
		util.Log().Info("初始管理员账号：" + c.Sprint("admin"))
		util.Log().Info("初始管理员密码：" + c.Sprint(password))
	}
}

func initDefaultSettings() {
	defaultSettings := []Setting{
		{Name: "db_version_" + conf.RequiredDBVersion, Value: `installed`, Type: "version"},
	}
	for _, setting := range defaultSettings {
		DB.Model(&Setting{}).Where(Setting{Name: setting.Name}).Create(&setting)
	}
}
