package main

import (
	"fmt"
	"im/config"
	"im/internal/upload"
	"im/pkg/logger"
	"im/pkg/minio"
)

const serverName = "im_file_control"

func main() {
	config.InitFileControl()
	cfg := config.Config
	logger.Init(fmt.Sprintf("%s/%s.log", cfg.Log.Path, serverName))
	//db.InitMysql()
	//control.Init()
	go minio.Init()
	router := control.NewRouter()
	address := cfg.Server.ControlListenAddr
	logger.Sugar.Infow("start control server", "address", address)
	err := router.Run(address)
	if err != nil {
		logger.Sugar.Errorw("control server run error", "error", err)
	}
}
