package conf

import (
	"github.com/go-ini/ini"
	"github.com/zhouqiaokeji/server/pkg/util"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
)

type database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
	DBFile      string
	Port        int
	Charset     string
}

type system struct {
	Mode           string `validate:"eq=master|eq=slave"`
	Listen         string `validate:"required"`
	Peer           int64
	AdminContainer []string
	Debug          bool
	SessionSecret  string
	HashIDSalt     string
}

type ssl struct {
	CertPath string `validate:"omitempty,required"`
	KeyPath  string `validate:"omitempty,required"`
	Listen   string `validate:"required"`
}

type unix struct {
	Listen string
}

type slave struct {
	Secret          string `validate:"omitempty,gte=64"`
	CallbackTimeout int    `validate:"omitempty,gte=1"`
	SignatureTTL    int    `validate:"omitempty,gte=1"`
}

type captcha struct {
	Height             int `validate:"gte=0"`
	Width              int `validate:"gte=0"`
	Mode               int `validate:"gte=0,lte=3"`
	ComplexOfNoiseText int `validate:"gte=0,lte=2"`
	ComplexOfNoiseDot  int `validate:"gte=0,lte=2"`
	IsShowHollowLine   bool
	IsShowNoiseDot     bool
	IsShowNoiseText    bool
	IsShowSlimeLine    bool
	IsShowSineLine     bool
	CaptchaLen         int `validate:"gt=0"`
}

// redis 配置
type redis struct {
	Network  string
	Server   string
	Password string
	DB       string
}

type thumb struct {
	MaxWidth      uint
	MaxHeight     uint
	FileSuffix    string `validate:"min=1"`
	MaxTaskCount  int
	EncodeMethod  string `validate:"eq=jpg|eq=png"`
	EncodeQuality int    `validate:"gte=1,lte=100"`
	GCAfterGen    bool
}

type cors struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
}

var cfg *ini.File

const defaultConf = `[System]
Mode   = master
Listen = :8089
Peer   = {Peer}
SessionSecret = {SessionSecret}
HashIDSalt = {HashIDSalt}
`

// Init Config Initialize
func Init(path string) {
	var err error

	if path == "" || !util.Exists(path) {
		// Create Config File
		confContent := util.Replace(map[string]string{
			"{SessionSecret}": util.RandStringRunes(64),
			"{HashIDSalt}":    util.RandStringRunes(64),
			"{Peer}":          strconv.Itoa(util.RandInt(0, 9)),
		}, defaultConf)
		f, err := util.CreatNestedFile(path)
		if err != nil {
			util.Log().Panic("Create Config File Fail, %s", err)
		}

		// Write Config
		_, err = f.WriteString(confContent)
		if err != nil {
			util.Log().Panic("Write Config Fail, %s", err)
		}

		f.Close()
	}

	cfg, err = ini.Load(path)
	if err != nil {
		util.Log().Panic("Resolve Config File Fail '%s': %s", path, err)
	}

	sections := map[string]interface{}{
		"Database":   DatabaseConfig,
		"System":     SystemConfig,
		"SSL":        SSLConfig,
		"UnixSocket": UnixConfig,
		"Captcha":    CaptchaConfig,
		"Redis":      RedisConfig,
		"Thumbnail":  ThumbConfig,
		"CORS":       CORSConfig,
		"Slave":      SlaveConfig,
	}
	for sectionName, sectionStruct := range sections {
		err = mapSection(sectionName, sectionStruct)
		if err != nil {
			util.Log().Panic("Config  %s Section Resolve Fail: %s", sectionName, err)
		}
	}

	// Reset Loglevel
	if !SystemConfig.Debug {
		util.Level = util.LevelInformational
		util.GloablLogger = nil
		util.Log()
	}

}

// mapSection Mapping Config Section To Struct
func mapSection(section string, confStruct interface{}) error {
	err := cfg.Section(section).MapTo(confStruct)
	if err != nil {
		return err
	}

	// Validate
	validate := validator.New()
	err = validate.Struct(confStruct)
	if err != nil {
		return err
	}

	return nil
}
