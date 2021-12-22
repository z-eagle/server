package appuseinfo

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	model "github.com/zhouqiaokeji/server/models"
	"github.com/zhouqiaokeji/server/pkg/rsa"
	"github.com/zhouqiaokeji/server/pkg/serializer"
	"github.com/zhouqiaokeji/server/pkg/util"
	"time"
)

type ServiceAppUseInfoDTO struct {
	Id           uint64    `json:"id,string"`
	Name         string    `json:"name"`
	Mobile       string    `json:"mobile"`
	ServerAddr   string    `json:"server_addr"`
	RequestIp    string    `json:"request_ip"`
	UserId       int64     `json:"user_id,string"`
	UserName     string    `json:"user_name"`
	Version      string    `json:"version"`
	Platform     string    `json:"platform"`
	WIFIName     string    `json:"wifi_name"`
	WIFIMac      string    `json:"wifi_mac"`
	BootLoader   string    `json:"boot_loader"`
	LON          float32   `json:"lon"`
	LAT          float32   `json:"lat"`
	OSVersion    string    `json:"os_version"`
	DeviceSN     string    `json:"device_sn"`
	DeviceVendor string    `json:"device_vendor"`
	Time         time.Time `json:"time"`
}

type SignAppUseInfo struct {
	Key        string `json:"key"`
	DeviceInfo string `json:"d_str" binding:"required"`
	UserInfo   string `json:"u_str" binding:"required"`
}

type LicenseUseInfos struct {
	Name        string                 `json:"name"`
	ContainerID string                 `json:"containerId" binding:"required"`
	UseInfos    []ServiceAppUseInfoDTO `json:"use_infos"`
}

func (s *SignAppUseInfo) Create(c *gin.Context) serializer.Response {

	useInfoDTO, err := s.decode(c)
	if err != nil {
		return serializer.Err(serializer.CodeNotFullySuccess, "信息解密异常", nil)
	}
	useInfo := &model.AppUseInfo{
		Name:         useInfoDTO.Name,
		Mobile:       useInfoDTO.Mobile,
		ServerAddr:   useInfoDTO.ServerAddr,
		RequestIp:    useInfoDTO.RequestIp,
		UserId:       useInfoDTO.UserId,
		UserName:     useInfoDTO.UserName,
		Version:      useInfoDTO.Version,
		Platform:     useInfoDTO.Platform,
		WIFIName:     useInfoDTO.WIFIName,
		WIFIMac:      useInfoDTO.WIFIMac,
		BootLoader:   useInfoDTO.BootLoader,
		LON:          useInfoDTO.LON,
		LAT:          useInfoDTO.LAT,
		OSVersion:    useInfoDTO.OSVersion,
		DeviceSN:     useInfoDTO.DeviceSN,
		DeviceVendor: useInfoDTO.DeviceVendor,
	}
	id, _ := useInfo.Create()
	return serializer.Response{
		Code: serializer.OK,
		Data: map[string]interface{}{
			"id": id,
		},
	}
}

func (s *ServiceAppUseInfoDTO) GetAppInfos(containerId string, page, size int, order string, date ...time.Time) serializer.Response {
	if containerId == "" {
		licenses, total := model.GetLicenses(page, size, "", "")
		res := make([]LicenseUseInfos, 0, len(licenses))

		for _, t := range licenses {
			resUseInfos, _ := getLicenseUseInfo(t, 1, 10, "")
			res = append(res, LicenseUseInfos{
				Name:        t.Name,
				ContainerID: t.ContainerID,
				UseInfos:    resUseInfos,
			})
		}
		return serializer.Response{
			Code: serializer.OK,
			Data: &serializer.Page{
				Total:   total,
				Content: res,
				Page:    page,
				Size:    size,
			},
		}
	}
	license, _ := model.GetLicense(containerId)
	res, total := getLicenseUseInfo(license, page, size, order, date...)
	return serializer.Response{
		Code: serializer.OK,
		Data: &serializer.Page{
			Total: total,
			Content: &LicenseUseInfos{
				ContainerID: license.ContainerID,
				Name:        license.Name,
				UseInfos:    res,
			},
			Page: page,
			Size: size,
		},
	}
}

func (s *SignAppUseInfo) decode(c *gin.Context) (ServiceAppUseInfoDTO, error) {
	var (
		err        error
		useInfoDTO ServiceAppUseInfoDTO
	)
	deviceInfoBase64 := rsa.BcryptRSA(s.DeviceInfo)
	if err = json.Unmarshal([]byte(deviceInfoBase64), &useInfoDTO); err != nil {
		util.Log().Error(err.Error())
	}
	userInfoBase64 := rsa.BcryptRSA(s.UserInfo)
	if err = json.Unmarshal([]byte(userInfoBase64), &useInfoDTO); err != nil {
		util.Log().Error(err.Error())
	}
	useInfoDTO.RequestIp = util.GetIpAddr(c.Request)
	return useInfoDTO, nil
}

func getLicenseUseInfo(license model.License, page, size int, order string, date ...time.Time) ([]ServiceAppUseInfoDTO, int64) {
	var (
		serverAddr []string
		total      int64
		infos      []model.AppUseInfo
		res        []ServiceAppUseInfoDTO
	)
	serverAddr = make([]string, 0, 2)
	if license.IP != "" {
		serverAddr = append(serverAddr, license.IP)
	}
	if license.Domain != "" {
		serverAddr = append(serverAddr, license.Domain)
	}
	if len(serverAddr) > 0 {
		infos, total = model.GetAppUseInfoByServerAddr(page, size, order, serverAddr, date)
		res = make([]ServiceAppUseInfoDTO, 0, len(infos))
		for _, t := range infos {
			res = append(res, ServiceAppUseInfoDTO{
				Id:           t.ID,
				Name:         t.Name,
				Mobile:       t.Mobile,
				ServerAddr:   t.ServerAddr,
				RequestIp:    t.RequestIp,
				UserId:       t.UserId,
				UserName:     t.UserName,
				Version:      t.Version,
				Platform:     t.Platform,
				WIFIName:     t.WIFIName,
				WIFIMac:      t.WIFIMac,
				BootLoader:   t.BootLoader,
				LON:          t.LON,
				LAT:          t.LAT,
				OSVersion:    t.OSVersion,
				DeviceSN:     t.DeviceSN,
				DeviceVendor: t.DeviceVendor,
				Time:         t.CreatedAt,
			})
		}
	}
	return res, total
}
