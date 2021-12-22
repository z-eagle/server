package holidays

import (
	model "github.com/zhouqiaokeji/server/models"
	"github.com/zhouqiaokeji/server/models/datatypes"
	"github.com/zhouqiaokeji/server/pkg/serializer"
)

type HolidayDTO struct {
	Id      int64          `json:"id"`
	Year    int            `json:"year" binding:"required"`
	Holiday datatypes.JSON `json:"holiday"`
}

func (s *HolidayDTO) Create() serializer.Response {

	holiday := &model.Holidays{
		Year:    s.Year,
		Holiday: s.Holiday,
	}
	id, _ := holiday.Create()
	return serializer.Response{
		Code: serializer.OK,
		Data: map[string]interface{}{
			"id": id,
		},
	}
}

func (s *HolidayDTO) GetHoliday() serializer.Response {
	holiday, err := model.GetHolidayByYear(s.Year)
	if err != nil {
		return serializer.Response{
			Code: serializer.OK,
			Data: map[string]interface{}{
				"year": s.Year,
			},
		}
	}
	return serializer.Response{
		Code: serializer.OK,
		Data: map[string]interface{}{
			"year":    holiday.Year,
			"holiday": holiday.Holiday,
		},
	}
}
