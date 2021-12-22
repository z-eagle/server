package models

import (
	"github.com/zhouqiaokeji/server/models/datatypes"
	"github.com/zhouqiaokeji/server/pkg/util"
)

type Holidays struct {
	Auditable
	Year    int            `json:"year" gorm:"index"`
	Holiday datatypes.JSON `json:"holiday"`
}

type Holiday struct {
	Date string `json:"date"`
	Name string `json:"name"`
}

// Create 创建年度假期记录
func (hol *Holidays) Create() (uint64, error) {
	if err := DB.Create(hol).Error; err != nil {
		util.Log().Warning("无法插入任务记录, %s", err)
		return 0, err
	}
	return hol.ID, nil
}

// GetHolidayByYear 获取指定年度假期
func GetHolidayByYear(year int) (*Holidays, error) {
	holiday := &Holidays{}
	result := DB.Where("year = ?", year).First(holiday)
	return holiday, result.Error
}
