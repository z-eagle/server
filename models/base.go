package models

import (
	"github.com/zhouqiaokeji/server/pkg/id"
	"gorm.io/gorm"
	"time"
)

type Identical struct {
	ID uint64 `gorm:"primary_key;autoIncrement:false"`
}

type Auditable struct {
	Identical
	CreatedAt time.Time
	Creator   uint64 `gorm:"index"`
	UpdatedAt time.Time
	Updater   uint64         `gorm:"index"`
	Deleted   gorm.DeletedAt `gorm:"index"`
}

func (u *Identical) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = id.GeneratorId.NextId()
	return
}
