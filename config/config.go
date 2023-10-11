package config

import (
	"fmt"
	"im/pkg/logger"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	dbLogger "gorm.io/gorm/logger"
)

var Config config

type config struct {
	Station string `yaml:"station"`

	MySQL struct {
		DSN          string `yaml:"dsn"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		Address      string `yaml:"address"`
		DatabaseName string `yaml:"database_name"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
		MaxLifeTime  int    `yaml:"max_life_time"`
	} `yaml:"mysql"`

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	}

	Server struct {
		ApiListenAddr        string `yaml:"api_listen_addr"`
		CmsApiListenAddr     string `yaml:"cms_api_listen_addr"`
		ConnectWSListenAddr  string `yaml:"connect_ws_listen_addr"`
		ConnectRPCListenAddr string `yaml:"connect_rpc_listen_addr"`
		ConnectRPCAddr       string `yaml:"connect_rpc_addr"`
		ControlListenAddr    string `yaml:"control_listen_addr"`
	}

	EMQXServer struct {
		MQTTAddress       string `yaml:"mqtt_address"`
		MQTTUsername      string `yaml:"mqtt_username"`
		MQTTPassword      string `yaml:"mqtt_password"`
		MQTTMaxConnection int    `yaml:"mqtt_max_connection"`
		APIAddress        string `yaml:"api_address"`
		APIUsername       string `yaml:"api_username"`
		APIPassword       string `yaml:"api_password"`
	} `yaml:"emqx_server"`

	Agora struct {
		AppID           string `yaml:"app_id"`
		AppSecret       string `yaml:"app_secret"`
		TokenExpireTime int64  `yaml:"token_expire_time"`
	}

	Log struct {
		Level          string `yaml:"level"`
		Path           string `yaml:"path"`
		Target         string `yaml:"target"`
		RecordKeepDays int    `yaml:"recordKeepDays"`
	}

	Cms struct {
		CtrlApi   string `yaml:"control"`
		Ip2region string `yaml:"ip2region"`
	}

	TokenPolicy struct {
		AccessSecret string `yaml:"accessSecret"`
		AccessExpire int64  `yaml:"accessExpire"`
	}
	Captcha struct {
		CacheExpireSec                       int    `yaml:"cacheExpireSec"`
		DefaultFont                          string `yaml:"defaultFont"`
		DefaultResourceRoot                  string `yaml:"defaultResourceRoot"`
		DefaultText                          string `yaml:"defaultText"`
		DefaultTemplateImageDirectory        string `yaml:"defaultTemplateImageDirectory"`
		DefaultBackgroundImageDirectory      string `yaml:"defaultBackgroundImageDirectory"`
		DefaultClickBackgroundImageDirectory string `yaml:"defaultClickBackgroundImageDirectory"`
	}

	Sms struct {
		AppID       string `yaml:"app_id"`
		AppSecret   string `yaml:"app_secret"`
		Signature   string `yaml:"signature"`
		TemplateId  string `yaml:"template_id"`
		ExpireTTL   int    `yaml:"expireTTL"`
		CodeTTL     int    `yaml:"codeTTL"`
		SuperCode   string `yaml:"superCode"`
		DxbUsername string `yaml:"dxbUsername"`
		DxbApiKey   string `yaml:"dxbApiKey"`
	}

	Minio struct {
		Endpoint        string `yaml:"endpoint"`
		OssPrefix       string `yaml:"oss_prefix"`
		OssBucket       string `yaml:"oss_bucket"`
		AccessKeyID     string `yaml:"access_key_id"`
		SecretAccessKey string `yaml:"secret_access_key"`
	} `yaml:"minio"`

	AliYun struct {
		AccessKeyId     string `yaml:"accessKeyId"`
		AccessKeySecret string `yaml:"accessKeySecret"`

		OSS struct {
			EndPoint   string `yaml:"endPoint"`
			BucketName string `yaml:"bucketName"`
			Domain     string `yaml:"domain"`
		} `yaml:"oss"`

		STS struct {
			RoleArn    string `yaml:"roleArn"`
			ExpireTime int64  `yaml:"expireTime"`
		}
	} `yaml:"aliyun"`

	Kafka struct {
		KafkaAddress      string `yaml:"kafka_address"`
		KafkaHistoryTopic string `yaml:"kafka_history_topic"`
	} `yaml:"kafka"`

	DefaultIcon struct {
		DiscoverIcon string `yaml:"discover_icon"`
	} `yaml:"default_icon"`

	Jpush string `yaml:"jpush"`
}

func Init() {
	cfgName := os.Getenv("CONFIG_NAME")
	if len(cfgName) == 0 {
		panic("init config error, config name not found")
	}

	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(fmt.Sprintf("init config error, read file error, err: %v", err))
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		panic(fmt.Sprintf("init config error, yaml unmarshal error, err: %v", err))
	}
	switch Config.Log.Level {
	case "debug":
		logger.Level = zap.DebugLevel
		logger.DBLevel = dbLogger.Info
	case "info":
		logger.Level = zap.InfoLevel
		logger.DBLevel = dbLogger.Info
	case "warn":
		logger.Level = zap.WarnLevel
		logger.DBLevel = dbLogger.Warn
	default:
		logger.Level = zap.ErrorLevel
		logger.DBLevel = dbLogger.Error
	}

	switch Config.Log.Target {
	case logger.File:
		logger.Target = logger.File
	default:
		logger.Target = logger.Console
	}
}

func InitControl() {
	cfgName := os.Getenv("CONFIG_NAME_CONTROL")
	if len(cfgName) == 0 {
		panic("init config error, config name not found")
	}

	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(fmt.Sprintf("init config error, read file error, err: %v", err))
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		panic(fmt.Sprintf("init config error, yaml unmarshal error, err: %v", err))
	}

	switch Config.Log.Level {
	case "debug":
		logger.Level = zap.DebugLevel
	case "info":
		logger.Level = zap.InfoLevel
		logger.DBLevel = dbLogger.Info
	case "warn":
		logger.Level = zap.WarnLevel
		logger.DBLevel = dbLogger.Warn
	default:
		logger.Level = zap.ErrorLevel
		logger.DBLevel = dbLogger.Error
	}

	switch Config.Log.Target {
	case logger.File:
		logger.Target = logger.File
	default:
		logger.Target = logger.Console
	}
}

func InitFileControl() {
	cfgName := os.Getenv("FILE_CONFIG_NAME_CONTROL")
	if len(cfgName) == 0 {
		panic("init config error, config name not found")
	}

	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(fmt.Sprintf("init config error, read file error, err: %v", err))
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		panic(fmt.Sprintf("init config error, yaml unmarshal error, err: %v", err))
	}

	switch Config.Log.Level {
	case "debug":
		logger.Level = zap.DebugLevel
	case "info":
		logger.Level = zap.InfoLevel
		logger.DBLevel = dbLogger.Info
	case "warn":
		logger.Level = zap.WarnLevel
		logger.DBLevel = dbLogger.Warn
	default:
		logger.Level = zap.ErrorLevel
		logger.DBLevel = dbLogger.Error
	}

	switch Config.Log.Target {
	case logger.File:
		logger.Target = logger.File
	default:
		logger.Target = logger.Console
	}
}
