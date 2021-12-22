package license

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	model "github.com/zhouqiaokeji/server/models"
	"github.com/zhouqiaokeji/server/pkg/rsa"
	"github.com/zhouqiaokeji/server/pkg/serializer"
	"github.com/zhouqiaokeji/server/pkg/util"
	"time"
)

type ServiceLicenseAuthDTO struct {
	Name        string       `json:"name"`
	ContainerID string       `json:"containerId" `
	Status      model.Status `json:"status"`
}
type ServiceLicenseDTO struct {
	Name        string       `json:"name"`
	ContainerID string       `json:"containerId" binding:"required"`
	Status      model.Status `json:"status"`
	IP          string       `json:"ip" `
	Domain      string       `json:"domain" `
	Expire      string       `json:"expire"`
	Time        time.Time    `json:"time" `
}

type SignLicense struct {
	Key  string `json:"key" binding:"required"`
	Info string `json:"info" binding:"required"`
}

// Create 新增并返回加密授权信息
func (s *SignLicense) Create() serializer.Response {
	var (
		data []byte
		err  error
	)
	license, err := s.decode()
	license.Status = model.Disabled
	if model.CheckExistByContainer(nil, license.ContainerID) {
		verifyLicense, _ := model.GetLicense(license.ContainerID)
		res, err := verifyLicense.Verify(license)
		if err != nil {
			return serializer.Response{
				Code: serializer.CodeCheckLogin,
				Data: "非法授权认证",
			}
		}
		if data, err = json.Marshal(&ServiceLicenseAuthDTO{
			Name:        res.Name,
			ContainerID: res.ContainerID,
			Status:      res.Status,
		}); err != nil {
			util.Log().Error(err.Error())
		}
		return serializer.Response{
			Code: serializer.OK,
			Data: map[string]interface{}{
				"sign": rsa.SignRSA(data),
				"data": string(data),
			},
		}
	}
	_, _ = license.Create()
	if data, err = json.Marshal(&ServiceLicenseAuthDTO{
		Name:        license.Name,
		ContainerID: license.ContainerID,
		Status:      license.Status,
	}); err != nil {
		util.Log().Error(err.Error())
	}
	return serializer.Response{
		Code: serializer.OK,
		Data: map[string]interface{}{
			"sign": rsa.SignRSA(data),
			"data": string(data),
		},
	}
}

// Verify 校验授权信息
func (s *SignLicense) Verify() serializer.Response {
	var (
		data []byte
	)
	license, _ := s.decode()
	if model.CheckExistByContainer(nil, license.ContainerID) {
		verifyLicense, _ := model.GetLicense(license.ContainerID)
		res, err := verifyLicense.Verify(license)
		if err != nil {
			return serializer.Response{
				Code: serializer.CodeCheckLogin,
				Data: "非法授权认证",
			}
		}
		if data, err = json.Marshal(&ServiceLicenseAuthDTO{
			Name:        res.Name,
			ContainerID: res.ContainerID,
			Status:      res.Status,
		}); err != nil {
			util.Log().Error(err.Error())
		}
		return serializer.Response{
			Code: serializer.OK,
			Data: map[string]interface{}{
				"sign": rsa.SignRSA(data),
				"data": string(data),
			},
		}
	}
	return serializer.Response{
		Code: serializer.CodeCheckLogin,
		Data: "非法授权认证",
	}
}

// UpdateLicense 更新授权信息
func (s *ServiceLicenseDTO) UpdateLicense() serializer.Response {
	license := model.License{
		Name:        s.Name,
		ContainerID: s.ContainerID,
		IP:          s.IP,
		Domain:      s.Domain,
		Expire:      s.Expire,
	}
	if model.CheckExistIpOrDomain(s.IP, "", s.ContainerID) {
		return serializer.Err(serializer.CodeNotFullySuccess, "IP 已绑定服务", nil)
	}
	if model.CheckExistIpOrDomain("", s.Domain, s.ContainerID) {
		return serializer.Err(serializer.CodeNotFullySuccess, "域名 已绑定服务", nil)
	}
	updateLicense, _ := model.GetLicense(s.ContainerID)
	updateLicense.Update(license)
	return serializer.Response{
		Code: serializer.OK,
		Data: s,
	}
}

// UpdateLicenseStatus 更新授权状态
func (s *ServiceLicenseDTO) UpdateLicenseStatus() serializer.Response {
	updateLicense, _ := model.GetLicense(s.ContainerID)
	updateLicense.UpdateStatus(s.Status)
	return serializer.Response{
		Code: serializer.OK,
		Data: s,
	}
}

// GetLicense 分页查询授权信息
func (s *ServiceLicenseDTO) GetLicense(name string, page, size int, order string) serializer.Response {
	licenses, total := model.GetLicenses(page, size, order, name)
	res := make([]ServiceLicenseDTO, 0, len(licenses))
	for _, t := range licenses {
		res = append(res, ServiceLicenseDTO{
			Name:        t.Name,
			ContainerID: t.ContainerID,
			Status:      t.Status,
			IP:          t.IP,
			Domain:      t.Domain,
			Time:        t.CreatedAt,
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

// Remove 删除授权信息
func (s *ServiceLicenseDTO) Remove(key string) serializer.Response {
	model.Remove(key)
	return serializer.Response{
		Code: serializer.OK,
	}
}

func (s *SignLicense) decode() (*model.License, error) {
	var err error
	licenseStr := rsa.BcryptRSA(s.Info)
	key := fmt.Sprintf("%x", md5.Sum([]byte(licenseStr)))
	if s.Key != key {
		return nil, errors.New("信息校验失败")
	}
	license := &model.License{}
	if err = json.Unmarshal([]byte(licenseStr), &license); err != nil {
		util.Log().Error(err.Error())
		return nil, errors.New("信息解析失败")
	}

	return license, nil
}
