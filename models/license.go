package models

import (
	"errors"
	"github.com/zhouqiaokeji/server/pkg/conf"
	"github.com/zhouqiaokeji/server/pkg/util"
	"time"
)

type Status int

const (
	VALID Status = iota // 开始生成枚举值, 默认为0
	INVALID
	EXPIRE
)

type License struct {
	Auditable
	Name           string `json:"name"`
	ContainerID    string `json:"containerId" gorm:"index"`
	Status         Status `json:"status"`
	IP             string `json:"ip" gorm:"index"`
	Domain         string `json:"domain" gorm:"index"`
	Expire         string `json:"expire"`
	LastOnlineTime string `json:"last_online_time"`
}

// Create 记录服务使用信息
func (lic *License) Create() (uint64, error) {
	if err := DB.Create(lic).Error; err != nil {
		util.Log().Warning("无法插入授权记录, %s", err)
		return 0, err
	}
	return lic.ID, nil
}

// GetLicenses 分页查询授权信息
func GetLicenses(page, size int, order string, name string) ([]License, int64) {
	var (
		licenses []License
		total    int64
	)
	dbChain := DB
	if name != "" {
		dbChain = dbChain.Where("name like ?", "%"+name+"%")
	}

	// 计算总数用于分页
	dbChain.Model(&License{}).Count(&total)

	// 查询记录
	dbChain.Limit(size).Offset((page - 1) * size).Order(order).Find(&licenses)

	return licenses, total
}

// CheckExistByContainer 判断是否存在授权
func CheckExistByContainer(id *uint64, containerId string) bool {
	var (
		license License
	)
	var res = false
	if id != nil {
		result := DB.First(&license, id)
		if result.Error == nil && &license.ID != nil {
			res = true
		}
	}
	if !res && containerId != "" {
		result := DB.Set("gorm:auto_preload", true).Where("container_id = ?", containerId).First(&license)
		if result.Error == nil && &license.ID != nil {
			res = true
		}
	}
	return res
}

// CheckExistIpOrDomain 判断IP或Domain 是否存在
func CheckExistIpOrDomain(ip, domain, containerId string) bool {
	var (
		license License
	)
	if ip != "" {
		result := DB.Set("gorm:auto_preload", true).Where("ip = ? AND container_id <> ?", ip, containerId).First(&license)
		if result.Error == nil && &license.ID != nil {
			return true
		}
	}
	if domain != "" {
		result := DB.Set("gorm:auto_preload", true).Where("domain = ? AND container_id <> ?", domain, containerId).First(&license)
		if result.Error == nil && &license.ID != nil {
			return true
		}
	}
	return false
}

// GetLicense 获取容器授权信息
func GetLicense(key string) (License, error) {
	var license License
	result := DB.Set("gorm:auto_preload", true).Where("container_id = ?", key).First(&license)
	return license, result.Error
}

// Verify 校验授权信息
func (lic *License) Verify(license *License) (*License, error) {
	online := false
	if lic.LastOnlineTime != "" && !util.ContainsString(conf.SystemConfig.AdminContainer, license.ContainerID) {
		loc, _ := time.LoadLocation("Local")
		lastTime, _ := time.ParseInLocation(util.FORMAT_DATETIME_Y4MDHMS, license.LastOnlineTime, loc)
		duration := time.Now().Sub(lastTime)
		online = duration <= time.Minute*5
	}
	if lic.Expire != "" {
		loc, _ := time.LoadLocation("Local")
		expire, _ := time.ParseInLocation(util.FORMAT_DATE_y4Md, lic.Expire, loc)
		duration := time.Now().After(expire)
		if duration {
			lic.Status = EXPIRE
		}
	}
	if lic.Status == Active && lic.Name == license.Name && !online {
		lic.updateOnline(lic.ID)
		return lic, nil
	}
	return nil, errors.New("授权信息异常 ")
}

// UpdateStatus 修改授权状态
func (lic *License) UpdateStatus(status Status) {
	DB.Model(&lic).Update("status", status)
}

// Update 更新授权信息
func (lic *License) Update(license License) {
	DB.First(lic)

	if license.IP != "" {
		lic.IP = license.IP
	}
	if license.Domain != "" {
		lic.Domain = license.Domain
	}
	if license.Expire != "" {
		lic.Expire = license.Expire
	}

	DB.Save(lic)
}

func (lic *License) updateOnline(id uint64) {
	lic.ID = id
	DB.Model(lic).Update("last_online_time", time.Now().Format(util.FORMAT_DATETIME_Y4MDHMS))
}

// Remove 删除指定授权信息
func Remove(key string) {
	DB.Where("container_id = ?", key).Delete(&License{})
}
