package middleware

import (
	"github.com/gin-gonic/gin"
	model "github.com/zhouqiaokeji/server/models"
	"github.com/zhouqiaokeji/server/pkg/cache"
	"github.com/zhouqiaokeji/server/pkg/jwt"
	"github.com/zhouqiaokeji/server/pkg/serializer"
	"github.com/zhouqiaokeji/server/pkg/util"
	"strconv"
	"strings"
)

const (
	CurrUser        = "user"
	CacheUserPerfix = "USER_DETAIL:"
	expire          = 30 * 60 * 1000
)

// CurrentUser fetch current user
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		if token == "" {
			c.JSON(serializer.CodeCheckLogin, serializer.CheckLogin())
			c.Abort()
			return
		}
		tokenParts := strings.Split(token, " ")
		// 验证token前缀
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != strings.ToLower(strings.TrimSpace(jwt.Prefix)) {
			c.JSON(serializer.CodeNoPermissionErr, serializer.Response{Code: serializer.CodeNoPermissionErr, Error: "authorization header format must be Bearer {token}"})
			c.Abort()
			return

		}
		if jwtUser, err := jwt.ValidateToken(tokenParts[1], jwt.User{IpAddr: util.GetIpAddr(c.Request), Terminal: util.GetTerminal(c.Request)}); jwtUser == nil || err != nil {
			c.JSON(serializer.CodeNoPermissionErr, serializer.Response{Code: serializer.CodeNoPermissionErr, Error: err.Error()})
			c.Abort()
			return
		} else {
			user, exist := cache.Get(CacheUserPerfix + strconv.Itoa(int(jwtUser.Id)))
			if !exist {
				user, err = model.GetUserByID(jwtUser.Id)
				_ = cache.Set(CacheUserPerfix+strconv.Itoa(int(jwtUser.Id)), &user, expire)
			}
			if user == nil {
				c.Set(CurrUser, &user)
			}
			c.Next()
		}
	}
}

// Fetch the security token set by the client.
func getToken(c *gin.Context) (token string) {
	token = c.GetHeader(jwt.Header)
	if token != "" {
		return token
	}
	token = c.DefaultQuery(jwt.Param, c.PostForm(jwt.Param))
	if token != "" {
		return token
	}
	if cookie, err := c.Cookie(jwt.Cookie); err == nil && cookie != "" {
		return cookie
	}
	return ""
}
