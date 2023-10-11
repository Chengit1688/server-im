package job

import (
	"fmt"
	"go.uber.org/zap"
	"im/config"
	configModel "im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	roleModel "im/internal/cms_api/role/model"
	roleRepo "im/internal/cms_api/role/repo"
	"im/pkg/logger"
	"im/pkg/util"
	"io/ioutil"
	"net/http"
)

func GetCtrlMenu() {
	timeString, err := configRepo.ConfigRepo.MenuGetConfigTime()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		return
	}
	url := fmt.Sprintf("%s/menu_config?timestamp=%s", config.Config.Cms.CtrlApi, timeString)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("Sync Menu", "同步菜单出错 请求出错"), zap.String("err", err.Error()))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("Sync Menu", "同步菜单出错 请求出错"), zap.String("StatusCode", resp.Status))
		return
	}
	if resp.StatusCode != 200 {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("Sync Menu", "同步菜单出错 请求出错"), zap.String("StatusCode", resp.Status))
		return
	}

	menuConf := new(configModel.GetMenuConfigResp)
	util.JsonUnmarshal(body, &menuConf)

	if menuConf.Code == 0 {

		respMenus := menuConf.Data.Menus
		respFixedMenus := fixTimeFormat(respMenus)
		menus := new([]roleModel.CmsMenu)
		util.CopyStructFields(&menus, &respFixedMenus)
		err := roleRepo.RoleRepo.SyncMenu(*menus)
		if err != nil {
			logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("Sync Menu", "同步菜单出错 更新数据库错误"))
		} else {
			logger.Sugar.Info(zap.String("func", util.GetSelfFuncName()), zap.String("Sync Menu", "同步菜单成功"), zap.String("时间戳", util.Int64ToString(menuConf.Data.Timestamp)))
			configRepo.ConfigRepo.MenuUpdateConfigTime(util.Int64ToString(menuConf.Data.Timestamp))
		}
	} else {
		if menuConf.Code == 201 {
			logger.Sugar.Info(zap.String("func", util.GetSelfFuncName()), zap.String("Sync Menu", "最新配置 无需更新"))
		} else {
			logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("Sync Menu", "同步菜单出错 请求出错"))
		}
	}
}

func fixTimeFormat(src []configModel.MenuListRespItem) (dst []configModel.MenuListRespItem) {
	for _, item := range src {

		if item.DeletedAt.Unix() < 0 {
			item.DeletedAt = nil
		}
		dst = append(dst, item)
	}
	return
}
