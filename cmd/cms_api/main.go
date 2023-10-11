package main

import (
	"fmt"
	"im/config"
	"im/internal/cms_api"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/minio"
	"im/pkg/mqtt"
	"im/pkg/oss"
	"im/pkg/util"
)

const serverName = "im_cms_api"

func main() {
	config.Init()
	cfg := config.Config
	logger.Init(fmt.Sprintf("%s/%s.log", cfg.Log.Path, serverName))
	db.Init()
	cms_api.Init()
	util.SetupCasbin(db.DB, "")
	util.InitIpRegion()
	go minio.Init()
	go oss.InitOss()
	mqtt.Init()

	router := cms_api.NewRouter()
	address := cfg.Server.CmsApiListenAddr
	logger.Sugar.Infow("start api server", "address", address)

	err := router.Run(address)
	if err != nil {
		logger.Sugar.Errorw("api server run error", "error", err)
	}
}
