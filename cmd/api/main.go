package main

import (
	"fmt"
	"im/config"
	"im/internal/api"
	"im/pkg/oss"
	//"im/internal/cms_api/config/usecase"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/minio"
	"im/pkg/mqtt"

	//"im/pkg/push"
	"im/pkg/util"
)

const serverName = "im_api"

func main() {
	config.Init()
	cfg := config.Config
	logger.Init(fmt.Sprintf("%s/%s.log", cfg.Log.Path, serverName))
	db.Init()
	api.Init()
	go minio.Init()
	go oss.InitOss()
	util.InitIpRegion()
	mqtt.Init()
	//极光推送 注释
	//appkey, masterSecret, _ := usecase.ConfigUseCase.GetJPushAuthInfo()
	//push.Init(appkey, masterSecret)
	logger.Sugar.Infow("start api server", "启动参数是", cfg)
	router := api.NewRouter()
	address := cfg.Server.ApiListenAddr
	logger.Sugar.Infow("start api server", "address", address)

	err := router.Run(address)
	if err != nil {
		logger.Sugar.Errorw("api server run error", "error", err)
	}
}
