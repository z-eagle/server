package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhouqiaokeji/server/service/user"
)

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var service user.LoginService
	service.UserName = c.PostForm("username")
	service.Password = c.PostForm("password")
	res := service.Login(c)
	c.JSON(200, res)
}
