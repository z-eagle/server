package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhouqiaokeji/server/pkg/util"
	"github.com/zhouqiaokeji/server/service/license"
	"strconv"
)

func GetLicenses(ctx *gin.Context) {
	var (
		service license.ServiceLicenseDTO
	)
	start, _ := strconv.Atoi(ctx.DefaultQuery("page", ctx.DefaultPostForm("page", "1")))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("size", ctx.DefaultPostForm("size", "20")))
	order := ctx.Query("order")
	res := service.GetLicense("", start, limit, order)
	ctx.JSON(200, res)

}
func CreateLicense(ctx *gin.Context) {
	var service = &license.SignLicense{}
	if err := ctx.BindJSON(service); err != nil {
		util.Log().Error(err.Error())
		ctx.JSON(200, ErrorResponse(err))
	}
	res := service.Create()
	ctx.JSON(200, res)
}
func GetLicense(ctx *gin.Context) {
	id := ctx.Query("id")
	var service = &license.ServiceLicenseDTO{}
	if id == "" {
		ctx.JSON(200, ErrorResponse(errors.New("id must not null")))
	}
	res := service.Remove(id)
	ctx.JSON(200, res)
}
func ChangeStatus(ctx *gin.Context) {
	var service = &license.ServiceLicenseDTO{}
	if err := ctx.ShouldBindJSON(&service); err != nil {
		util.Log().Error(err.Error())
		ctx.JSON(200, ErrorResponse(err))
	}
	res := service.UpdateLicenseStatus()
	ctx.JSON(200, res)
}
func BindLicense(ctx *gin.Context) {
	var service = &license.ServiceLicenseDTO{}
	if err := ctx.ShouldBindJSON(&service); err != nil {
		util.Log().Error(err.Error())
		ctx.JSON(200, ErrorResponse(err))
	}
	res := service.UpdateLicense()
	ctx.JSON(200, res)
}
func RemoveLicense(ctx *gin.Context) {
	id := ctx.Query("id")
	var service = &license.ServiceLicenseDTO{}
	if id == "" {
		ctx.JSON(200, ErrorResponse(errors.New("id must not null")))
	}
	res := service.Remove(id)
	ctx.JSON(200, res)
}
func VerifyLicense(ctx *gin.Context) {
	var service = &license.SignLicense{}
	if err := ctx.BindJSON(service); err != nil {
		util.Log().Error(err.Error())
		ctx.JSON(200, ErrorResponse(err))
	}
	res := service.Verify()
	ctx.JSON(200, res)
}
