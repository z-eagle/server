package conf

const (
	CONF_FILE_NAME = "conf.ini"
	MODE_MASTER    = "master"
	MODE_SLAVE     = "slave"
	MODE_BOTH      = "both"
)

// RedisConfig Redis Serve Config
var RedisConfig = &redis{
	Network:  "tcp",
	Server:   "",
	Password: "",
	DB:       "0",
}

// DatabaseConfig Database Confog
var DatabaseConfig = &database{
	Type:        "UNSET",
	Charset:     "utf8",
	DBFile:      "data.db",
	TablePrefix: "t_",
	Port:        3306,
}

// SystemConfig System Public Config
var SystemConfig = &system{
	Debug:          false,
	AdminContainer: []string{"a6b67a6b8c7bc39b55e8b0af5e771d6e8b6a9aa382d872af611118282d5ab886"},
	Mode:           "master",
	Listen:         ":8099",
}

// CaptchaConfig Captcha Config
var CaptchaConfig = &captcha{
	Height:             60,
	Width:              240,
	Mode:               3,
	ComplexOfNoiseText: 0,
	ComplexOfNoiseDot:  0,
	IsShowHollowLine:   false,
	IsShowNoiseDot:     false,
	IsShowNoiseText:    false,
	IsShowSlimeLine:    false,
	IsShowSineLine:     false,
	CaptchaLen:         6,
}

// CORSConfig CORS Config
var CORSConfig = &cors{
	AllowOrigins:     []string{"UNSET"},
	AllowMethods:     []string{"PUT", "POST", "GET", "OPTIONS"},
	AllowHeaders:     []string{"Cookie", "X-Cr-Policy", "x-token", "Content-Length", "Content-Type", "X-Cr-Path", "X-Cr-FileName"},
	AllowCredentials: false,
	ExposeHeaders:    nil,
}

// ThumbConfig Picture Thumb Config
var ThumbConfig = &thumb{
	MaxWidth:      400,
	MaxHeight:     300,
	FileSuffix:    "._thumb",
	MaxTaskCount:  -1,
	EncodeMethod:  "jpg",
	GCAfterGen:    false,
	EncodeQuality: 85,
}

// SlaveConfig Slave Config
var SlaveConfig = &slave{
	CallbackTimeout: 20,
	SignatureTTL:    60,
}

// SSLConfig SSL Config
var SSLConfig = &ssl{
	Listen:   ":443",
	CertPath: "",
	KeyPath:  "",
}

// UnixConfig Unix Config
var UnixConfig = &unix{
	Listen: "",
}
