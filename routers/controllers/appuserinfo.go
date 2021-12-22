package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhouqiaokeji/server/pkg/util"
	"github.com/zhouqiaokeji/server/service/appuseinfo"
	"strconv"
	"time"
)

// CreateAppUseInfo 新增服务使用记录
func CreateAppUseInfo(c *gin.Context) {
	var service appuseinfo.SignAppUseInfo
	if err := c.ShouldBindJSON(&service); err == nil {
		res := service.Create(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

func GetAppInfos(c *gin.Context) {
	type param struct {
		ContainerId string   `json:"containerId"`
		Time        []string `json:"time"`
	}
	var service appuseinfo.ServiceAppUseInfoDTO
	var params = &param{}
	page, _ := strconv.Atoi(c.DefaultQuery("page", c.DefaultPostForm("page", "1")))
	limit, _ := strconv.Atoi(c.DefaultQuery("size", c.DefaultPostForm("size", "20")))
	order := c.Query("order")
	if err := c.ShouldBindJSON(params); err == nil {
		var (
			start time.Time
			end   time.Time
		)
		if len(params.Time) == 2 {
			loc, _ := time.LoadLocation("Local")
			start, _ = time.ParseInLocation(util.FORMAT_DATETIME_y4Md, params.Time[0], loc)
			end, _ = time.ParseInLocation(util.FORMAT_DATETIME_Y4MDHMS, params.Time[1]+" 23:59:59", loc)
		}
		res := service.GetAppInfos(params.ContainerId, page, limit, order, start, end)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}
