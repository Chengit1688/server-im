package service

import (
	"encoding/json"
	"fmt"
	apiUserRepo "im/internal/api/user/repo"
	"im/internal/cms_api/config/model"
	"im/internal/cms_api/config/repo"
	configRepo "im/internal/cms_api/config/repo"
	ipwhitelistRepo "im/internal/cms_api/ipwhitelist/repo"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/response"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ConfigService = new(configService)

type configService struct{}

func (s *configService) GetRegisterConfig(c *gin.Context) {
	req := new(model.GetLoginConfigReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	config, err := configRepo.ConfigRepo.GetRegisterConfig()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.GetRegisterConfigResp)
	err = util.JsonUnmarshal([]byte(config.Content), &ret)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, ret)
}

func (s *configService) GetLoginConfig(c *gin.Context) {
	req := new(model.GetLoginConfigReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	config, err := configRepo.ConfigRepo.GetLoginConfig()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.GetLoginConfigResp)
	err = util.JsonUnmarshal([]byte(config.Content), &ret)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, ret)
}

func (s *configService) UpdateLoginConfig(c *gin.Context) {
	req := new(model.UpdateLoginConfigReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	config, err := util.JsonMarshal(req.Config)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	err = configRepo.ConfigRepo.UpdateLoginConfig(string(config))
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, nil)
}

func (s *configService) UpdateRegisterConfig(c *gin.Context) {
	req := new(model.UpdateRegisterConfigReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	config, err := util.JsonMarshal(req.Config)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	err = configRepo.ConfigRepo.UpdateRegisterConfig(string(config))
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, nil)
}

func (s *configService) HandleSignConfig(c *gin.Context) {
	req := new(model.SignLogConfigReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	config, err := util.JsonMarshal(req.Config)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	cf := &model.SettingConfig{
		Content:    string(config),
		ConfigType: configRepo.ConfigSignConfig,
	}
	err = configRepo.ConfigRepo.UpdateSignConfig(cf)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, nil)
}

func (s *configService) GetSignConfig(c *gin.Context) {
	var (
		err    error
		config model.SettingConfig
	)
	req := new(model.GetSignLogConfigReq)
	err = c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	config, err = configRepo.ConfigRepo.GetSignConfig()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		config = model.SettingConfig{
			ConfigType: configRepo.ConfigSignConfig,
			Content:    `{"sign":1}`,
		}
	}
	ret := new(model.SignLogConfigResp)
	err = util.JsonUnmarshal([]byte(config.Content), &ret)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, ret)
}

func (s *configService) GetCmsConfig(c *gin.Context) {
	req := new(model.GetCmsConfigReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	status, err := configRepo.ConfigRepo.GetGoogleCodeIsOpen()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	siteInfos, err := configRepo.ConfigRepo.GetCmsSiteInfo()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetCmsConfigResp)
	ret.GoogleCodeIsOpen = status
	for _, item := range siteInfos {
		switch item.Name {
		case configRepo.ConfigCmsSiteName:
			ret.UIInfo.SiteName = item.Value
		case configRepo.ConfigCmsLoginIcon:
			ret.UIInfo.LoginIcon = item.Value
		case configRepo.ConfigCmsLoginBackend:
			ret.UIInfo.LoginBackend = item.Value
		case configRepo.ConfigCmsPageIcon:
			ret.UIInfo.PageIcon = item.Value
		case configRepo.ConfigCmsMenuIcon:
			ret.UIInfo.MenuIcon = item.Value
		}
	}
	http.Success(c, ret)
}

func (s *configService) SetGoogleCodeIsOpen(c *gin.Context) {
	req := new(model.SetGoogleCodeIsOpenReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetGoogleCodeIsOpen(req.Status)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c, nil)
}

func (s *configService) GetJPushConfig(c *gin.Context) {
	req := new(model.GetJPushReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	configs, err := configRepo.ConfigRepo.GetJPushAuthInfo()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetJPushResp)
	for _, config := range configs {
		switch config.Name {
		case repo.ConfigJPushAppKey:
			ret.AppKey = config.Value
		case repo.ConfigJPushMasterSecret:
			ret.MasterSecret = config.Value
		}
	}
	http.Success(c, ret)
}

func (s *configService) SetJPushConfig(c *gin.Context) {
	req := new(model.SetJPushReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetJPushAuthInfo(req.AppKey, req.MasterSecret)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.SetJPushResp)
	ret.AppKey = req.AppKey
	ret.MasterSecret = req.MasterSecret
	http.Success(c, ret)
}

func (s *configService) GetSystemConfig(c *gin.Context) {
	req := new(model.GetLoginConfigReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("query error: %v", err))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	resp, err := configRepo.ConfigRepo.GetSystemConfig()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetParameterConfig error: %v", err))
		http.Failed(c, code.ErrUnknown)
		return
	}

	http.Success(c, resp)
}

func (s *configService) UpdateSystemConfig(c *gin.Context) {
	req := new(model.ParameterConfigReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(c, code.ErrBadRequest)
		return
	}

	config, err := util.JsonMarshal(req.Config)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(c, code.ErrUnknown)
		return
	}
	dbConfig, _ := apiUserRepo.UserCache.GetSystemConfigInfo(configRepo.ConfigRepo.GetSystemConfig)
	if dbConfig != nil {
		if dbConfig.IpRegLimitCount != req.Config.IpRegLimitCount {

			apiUserRepo.UserCache.ClearAllIPCount()
		}
	}
	m := model.ParameterConfigResp{}
	if err = json.Unmarshal([]byte(config), &m); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("make json, error: %v", err))
		http.Failed(c, code.ErrUnknown)
		return
	}
	if m.GroupLimit > 10000 || m.ContactsFriendLimit > 10000 || m.CreateGroupLimit > 10000 {
		http.Failed(c, code.ErrImSiteParmLimit)
		return
	}

	settingConfig := &model.SettingConfig{
		Content:    string(config),
		ConfigType: configRepo.ConfigParameterConfig,
	}
	err = configRepo.ConfigRepo.UpdateParameterConfig(settingConfig)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	apiUserRepo.UserCache.DelSystemConfigOnCache()

	var sysConfig *model.ParameterConfigResp
	if sysConfig, err = apiUserRepo.UserCache.GetSystemConfigInfo(configRepo.ConfigRepo.GetSystemConfig); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " BroadCastSysconfig error:", err)
	}

	if err == nil {

		if err = mqtt.BroadcastMessage(req.OperationID, common.SysConfigBroadcast, sysConfig); err != nil {
			logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " BroadCastSysconfig error:", err)
		}
	}

	http.Success(c)
}

func (s *configService) GetDepositeConfig(c *gin.Context) {
	req := new(model.GetDepositeReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	configs, err := configRepo.ConfigRepo.GetDepositInfo()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetDepositeResp)
	for _, config := range configs {
		switch config.Name {
		case repo.ConfigDepositHTML:
			ret.Html = config.Value
		case repo.ConfigDepositURL:
			ret.Url = config.Value
		case repo.ConfigDepositSWITCH:
			ret.Switch = util.String2Int(config.Value)
		}
	}
	http.Success(c, ret)
}

func (s *configService) SetDepositeConfig(c *gin.Context) {
	req := new(model.SetDepositeReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetDepositInfo(req.Html, req.Url, req.Switch)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetDepositeResp)
	ret.Html = req.Html
	ret.Url = req.Url
	ret.Switch = req.Switch
	http.Success(c, ret)
}

func (s *configService) GetWithdrawConfig(c *gin.Context) {
	req := new(model.GetWithdrawConfigReq)
	err := c.ShouldBindQuery(&req)
	lang := c.GetHeader("Locale")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("query error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	resp, err := configRepo.ConfigRepo.GetWithdrawConfig(lang)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetWithdrawConfig error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	http.Success(c, resp)
}

func (s *configService) SetWithdrawConfig(c *gin.Context) {
	req := new(model.SetWithdrawConfigReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	config := new(model.WithdrawConfigResp)
	util.CopyStructFields(&config, &req)
	configByte, err := util.JsonMarshal(config)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(c, code.ErrUnknown)
		return
	}
	settingConfig := &model.SettingConfig{
		Content:    string(configByte),
		ConfigType: configRepo.ConfigWithdrawConfig,
	}
	err = configRepo.ConfigRepo.UpdateWithdrawConfig(settingConfig)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}

func (s *configService) GetFeihuConfig(c *gin.Context) {
	req := new(model.GetFeihuReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	configs, err := configRepo.ConfigRepo.GetFeihuAuthInfo()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetFeihuResp)
	for _, config := range configs {
		switch config.Name {
		case repo.ConfigFeihuAppKey:
			ret.AppKey = config.Value
		case repo.ConfigFeihuAppSecret:
			ret.AppSecret = config.Value
		}
	}
	http.Success(c, ret)
}

func (s *configService) SetFeihuConfig(c *gin.Context) {
	req := new(model.SetFeihuReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetFeihuAuthInfo(req.AppKey, req.AppSecret)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.SetFeihuResp)
	ret.AppKey = req.AppKey
	ret.AppSecret = req.AppSecret
	http.Success(c, ret)
}

func (s *configService) GetAboutUs(c *gin.Context) {
	req := new(model.GetAboutUsReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	content, err := configRepo.ConfigRepo.GetAboutUs()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetAboutUsResp)
	ret.Content = content
	http.Success(c, ret)
}

func (s *configService) SetAboutUs(c *gin.Context) {
	req := new(model.SetAboutUsReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetAboutUs(req.Content)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetAboutUsResp)
	ret.Content = req.Content
	http.Success(c, ret)
}

func (s *configService) GetPrivacyPolicy(c *gin.Context) {
	req := new(model.GetAboutUsReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	content, err := configRepo.ConfigRepo.GetPrivacyPolicy()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetAboutUsResp)
	ret.Content = content
	http.Success(c, ret)
}

func (s *configService) SetPrivacyPolicy(c *gin.Context) {
	req := new(model.SetAboutUsReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetPrivacyPolicy(req.Content)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetAboutUsResp)
	ret.Content = req.Content
	http.Success(c, ret)
}

func (s *configService) GetUserAgreement(c *gin.Context) {
	req := new(model.GetAboutUsReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	content, err := configRepo.ConfigRepo.GetUserAgreement()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetAboutUsResp)
	ret.Content = content
	http.Success(c, ret)
}

func (s *configService) SetUserAgreement(c *gin.Context) {
	req := new(model.SetAboutUsReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetUserAgreement(req.Content)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetAboutUsResp)
	ret.Content = req.Content
	http.Success(c, ret)
}

func (s *configService) SetIPWhiteListIsOpen(c *gin.Context) {
	req := new(model.SetIPWhiteListIsOpenReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetIPWhiteIsOpen(req.Status)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	if req.Status == 1 {
		ipwhitelistRepo.IpWhiteListCache.Open()
	} else {
		ipwhitelistRepo.IpWhiteListCache.Close()
	}
	http.Success(c, nil)
}

func (s *configService) GetIPWhiteListIsOpen(c *gin.Context) {
	req := new(model.GetIPWhiteListIsOpenReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	status, err := configRepo.ConfigRepo.GetIPWhiteIsOpen()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c, status)
}

func (s *configService) SetDefaultIsOpen(c *gin.Context) {
	req := new(model.SetDefaultIsOpenReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetDefaultIsOpen(req.Status)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c, nil)
}

func (s *configService) GetDefaultIsOpen(c *gin.Context) {
	req := new(model.GetDefaultIsOpenReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	status, err := configRepo.ConfigRepo.GetDefaultIsOpen()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c, status)
}
