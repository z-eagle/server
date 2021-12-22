package models

import (
	"errors"
	"github.com/zhouqiaokeji/server/pkg/bcrypt"
	"gorm.io/gorm"
)

const (
	// Active 账户正常状态
	Active = iota
	// Disabled 禁用
	Disabled
)

var bCrypt = bcrypt.NewDefaultBcrypt()

type User struct {
	Auditable
	Email    string `gorm:"type:varchar(100);unique_index"`
	Nick     string `gorm:"size:50"`
	UserName string
	Password string `json:"-"`
	Status   int
}

func NewUser() User {
	return User{}
}

// GetUserByID 用ID获取用户
func GetUserByID(ID uint64) (User, error) {
	var user User
	result := DB.Set("gorm:auto_preload", true).First(&user, ID)
	return user, result.Error
}

// GetUserByUserName 用UserName获取用户
func GetUserByUserName(userName string) (User, bool) {
	var user User
	result := DB.Where("user_name = ?", userName).First(&user)
	return user, !errors.Is(result.Error, gorm.ErrRecordNotFound)
}

// GetActiveUserByUserName 用UserName获取可登录用户
func GetActiveUserByUserName(email string) (User, error) {
	var user User
	result := DB.Where("status = ? and user_name = ?", Active, email).First(&user)
	return user, result.Error
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) (bool, error) {
	return bCrypt.Matches(password, user.Password), nil
}

// SetPassword RSA解密后设定 User 的 Password 字段
func (user *User) SetPassword(password string) {
	user.Password = bCrypt.Encode(password)
}

// GetPasswordBcrypt 获取密码加密编码器
func GetPasswordBcrypt() *bcrypt.BCrypt {
	return bCrypt
}
