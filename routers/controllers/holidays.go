package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhouqiaokeji/server/service/holidays"
	"strconv"
	"time"
)

// CreateHoliday 新增服务使用记录
func CreateHoliday(c *gin.Context) {
	var service holidays.HolidayDTO
	if err := c.ShouldBindJSON(&service); err == nil {
		res := service.Create()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

func GetHolidays(c *gin.Context) {
	var service holidays.HolidayDTO
	var year int
	year, _ = strconv.Atoi(c.DefaultQuery("year", c.PostForm("year")))
	if year == 0 {
		year = time.Now().Year()
	}
	service.Year = year
	res := service.GetHoliday()
	c.JSON(200, res)
}
