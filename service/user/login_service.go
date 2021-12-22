package user

import (
	"github.com/gin-gonic/gin"
	model "github.com/zhouqiaokeji/server/models"
	"github.com/zhouqiaokeji/server/pkg/jwt"
	"github.com/zhouqiaokeji/server/pkg/rsa"
	"github.com/zhouqiaokeji/server/pkg/serializer"
	"github.com/zhouqiaokeji/server/pkg/util"
)

type LoginService struct {
	UserName string `form:"username" json:"username" binding:"required,username"`
	Password string `form:"password" json:"Password" binding:"required,min=4,max=64"`
}

//Login 用户登录
func (s *LoginService) Login(ctx *gin.Context) serializer.Response {

	expectedUser, exist := model.GetUserByUserName(s.UserName)
	// 一系列校验
	if !exist {
		return serializer.Err(serializer.CodeCredentialInvalid, "用户名未注册或添加", nil)
	}

	if authOK, _ := expectedUser.CheckPassword(rsa.BcryptRSA(s.Password)); !authOK {
		return serializer.Err(serializer.CodeCredentialInvalid, "用户名或密码错误", nil)
	}

	return serializer.Response{
		Code: serializer.OK,
		Data: map[string]interface{}{
			"username": expectedUser.UserName,
			"token": jwt.Prefix + jwt.GenerateToken(jwt.User{
				Id:       expectedUser.ID,
				UserName: expectedUser.UserName,
				IpAddr:   util.GetIpAddr(ctx.Request),
				Terminal: util.GetTerminal(ctx.Request),
			}),
		},
	}
}
