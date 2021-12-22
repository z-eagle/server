package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhouqiaokeji/server/pkg/conf"
	"github.com/zhouqiaokeji/server/pkg/util"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

// DB 数据库链接单例
var DB *gorm.DB

// Init 初始化数据库链接
func Init() {
	util.Log().Info("初始化数据库连接")

	var (
		db       *gorm.DB
		err      error
		logLevel logger.LogLevel
		dsn      string
	)

	// Debug模式下，输出所有 SQL 日志
	if conf.SystemConfig.Debug {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logLevel,    // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)

	options := &gorm.Config{
		PrepareStmt: true,
		Logger:      newLogger,
		NamingStrategy: schema.NamingStrategy{ // 处理表前缀
			TablePrefix:   conf.DatabaseConfig.TablePrefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,                            // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	}

	if gin.Mode() == gin.TestMode {
		// 测试模式下，使用内存数据库
		dsn = "file:test.db?cache=shared&mode=memory"
		db, err = gorm.Open(sqlite.Open(dsn), options)
	} else {
		switch conf.DatabaseConfig.Type {
		case "UNSET", "sqlite", "sqlite3":
			// 未指定数据库或者明确指定为 sqlite 时，使用 SQLite3 数据库
			dsn = fmt.Sprintf("file:%s?cache=shared", conf.DatabaseConfig.DBFile)
			db, err = gorm.Open(sqlite.Open(dsn), options)
		case "postgres":
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				conf.DatabaseConfig.Host,
				conf.DatabaseConfig.User,
				conf.DatabaseConfig.Password,
				conf.DatabaseConfig.Name,
				conf.DatabaseConfig.Port)
			db, err = gorm.Open(postgres.Open(dsn), options)
		case "mysql":
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
				conf.DatabaseConfig.User,
				conf.DatabaseConfig.Password,
				conf.DatabaseConfig.Host,
				conf.DatabaseConfig.Port,
				conf.DatabaseConfig.Name,
				conf.DatabaseConfig.Charset)
			db, err = gorm.Open(mysql.Open(dsn), options)
		case "mssql":
			dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
				conf.DatabaseConfig.User,
				conf.DatabaseConfig.Password,
				conf.DatabaseConfig.Host,
				conf.DatabaseConfig.Port,
				conf.DatabaseConfig.Name,
				conf.DatabaseConfig.Charset)
			db, err = gorm.Open(sqlserver.Open(dsn), options)
		default:
			util.Log().Panic("不支持数据库类型: %s", conf.DatabaseConfig.Type)
		}
	}

	if err != nil {
		util.Log().Panic("连接数据库不成功, %s", err)
	}

	//设置连接池
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	//空闲
	sqlDB.SetMaxIdleConns(50)
	//打开
	sqlDB.SetMaxOpenConns(100)
	//超时
	sqlDB.SetConnMaxLifetime(time.Second * 60)

	DB = db

	//执行迁移
	migration()
}
