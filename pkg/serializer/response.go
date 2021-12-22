package serializer

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"github.com/zhouqiaokeji/server/pkg/util"
)

const OK = 200

// CheckLogin 检查登录
func CheckLogin() Response {
	return Response{
		Code: CodeCheckLogin,
		Msg:  "用户未登录",
	}
}

// CheckJwt 检查jwt
func CheckJwt() Response {
	return Response{
		Code: CodeCheckLogin,
		Msg:  "鉴权信息异常",
	}
}

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

type Page struct {
	Total   int64       `json:"total"`
	Content interface{} `json:"content"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

// NewResponseWithGobData 返回Data字段使用gob编码的Response
func NewResponseWithGobData(data interface{}) Response {
	var w bytes.Buffer
	encoder := gob.NewEncoder(&w)
	if err := encoder.Encode(data); err != nil {
		return Err(CodeInternalSetting, "无法编码返回结果", err)
	}

	return Response{Data: w.Bytes()}
}

// GobDecode 将 Response 正文解码至目标指针
//goland:noinspection GoStandardMethods
func (r *Response) GobDecode(target interface{}) error {
	src := r.Data.(string)
	raw := make([]byte, len(src)*len(src)/base64.StdEncoding.DecodedLen(len(src)))
	_, err := base64.StdEncoding.Decode(raw, []byte(src))
	if err != nil {
		util.Log().Error("StdEncoding Decode err :%s", err.Error())
	}
	decoder := gob.NewDecoder(bytes.NewBuffer(raw))
	if err := decoder.Decode(target); err != nil {
		util.Log().Error("GobDecode err :%s", err.Error())
	}
	return nil
}
