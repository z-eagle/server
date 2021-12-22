package models

import (
	"github.com/zhouqiaokeji/server/pkg/util"
	"time"
)

type AppUseInfo struct {
	Auditable
	Name         string  `json:"name"`
	Mobile       string  `json:"mobile"`
	ServerAddr   string  `json:"server_addr" gorm:"index"`
	RequestIp    string  `json:"request_ip"`
	UserId       int64   `json:"user_id"`
	UserName     string  `json:"user_name"`
	Version      string  `json:"version"`
	Platform     string  `json:"platform"`
	WIFIName     string  `json:"wifi_name"`
	WIFIMac      string  `json:"wifi_mac"`
	BootLoader   string  `json:"boot_loader"`
	LON          float32 `json:"lon"`
	LAT          float32 `json:"lat"`
	OSVersion    string  `json:"os_version"`
	DeviceSN     string  `json:"device_sn"`
	DeviceVendor string  `json:"device_vendor"`
}

// Create 记录服务使用信息
func (useInfo *AppUseInfo) Create() (uint64, error) {
	if err := DB.Create(useInfo).Error; err != nil {
		util.Log().Warning("无法插入服务使用记录, %s", err)
		return 0, err
	}
	return useInfo.ID, nil
}

// GetAppUseInfoByServerAddr 获取指定服务地址使用信息
func GetAppUseInfoByServerAddr(page, size int, order string, serverAddr []string, time []time.Time) ([]AppUseInfo, int64) {
	var (
		useInfos []AppUseInfo
		total    int64
	)
	dbChain := DB
	dbChain = dbChain.Where("server_addr in (?)", serverAddr)
	if len(time) == 2 {
		dbChain = dbChain.Where("created_at BETWEEN ? AND ?", time[0], time[1])
	}

	// 计算总数用于分页
	dbChain.Model(&AppUseInfo{}).Count(&total)

	// 查询记录
	dbChain.Limit(size).Offset((page - 1) * size).Order(order).Find(&useInfos)

	return useInfos, total
}
