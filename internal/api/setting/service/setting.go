package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"im/config"
	settingModel "im/internal/api/setting/model"
	settingRepo "im/internal/api/setting/repo"
	"im/internal/api/user/usecase"
	"im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	configUsecase "im/internal/cms_api/config/usecase"
	discoverRepo "im/internal/cms_api/discover/repo"
	domainModel "im/internal/control/domain/model"
	"im/pkg/common/constant"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
)

var SettingService = new(settingService)

type settingService struct{}

func (s *settingService) RegisterAndLoginConfig(ctx *gin.Context) {
	var (
		err         error
		OperationID string
		r           []*settingModel.SettingConfig
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&OperationID); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	configOpt := settingRepo.WhereOption{
		ConfigType: []string{configRepo.ConfigRegisterConfig, configRepo.ConfigLoginConfig, configRepo.ConfigParameterConfig},
	}
	if r, err = settingRepo.SettingRepo.GetConfigInfo(configOpt); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get list, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	res := make(map[string]interface{})
	var m interface{}
	for _, v := range r {
		switch v.ConfigType {
		case configRepo.ConfigLoginConfig:
			m = &model.GetLoginConfigResp{}
		case configRepo.ConfigRegisterConfig:
			m = &model.GetRegisterConfigResp{}
		case configRepo.ConfigParameterConfig:
			m = &model.ParameterConfigResp{}
		}
		if err = json.Unmarshal([]byte(v.Content), &m); err != nil {
			logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("make json, error: %v", err))
			http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
			return
		}
		res[v.ConfigType] = m
	}

	if _, ok := res[configRepo.ConfigLoginConfig]; !ok {
		res[configRepo.ConfigLoginConfig] = model.GetLoginConfigResp{
			Pc:     []int{1, 2},
			Mobile: []int{1, 2, 3},
		}
	}
	if _, ok := res[configRepo.ConfigRegisterConfig]; !ok {
		res[configRepo.ConfigRegisterConfig] = model.GetRegisterConfigResp{
			CheckInviteCode:    2,
			IsInviteCode:       2,
			IsVerificationCode: 2,
			IsSmsCode:          2,
			IsAllAccount:       1,
		}
	}
	if _, ok := res[configRepo.ConfigParameterConfig]; !ok {
		res[configRepo.ConfigParameterConfig] = model.ParameterConfigResp{
			IpRegLimitCount:       0,
			IpRegLimitTime:        0,
			DeviceRegLimit:        0,
			GroupLimit:            1000,
			ContactsFriendLimit:   1000,
			CreateGroupLimit:      1000,
			UserIdPrefix:          "用户",
			FileSizeLimit:         100,
			HistoryTime:           30,
			IsMemberAddFriend:     constant.SwitchOff,
			IsMemberDelFriend:     constant.SwitchOff,
			IsMemberAddGroup:      constant.SwitchOff,
			IsOpenRegister:        constant.SwitchOn,
			IsShowMemberStatus:    constant.SwitchOn,
			IsMsgReadStatus:       constant.SwitchOn,
			IsOpenRedPack:         constant.SwitchOn,
			IsOpenRedPackSingle:   constant.SwitchOn,
			IsNormalSeeAddress:    constant.SwitchOff,
			IsOpenVoiceCall:       constant.SwitchOn,
			IsOpenCameraCall:      constant.SwitchOn,
			IsNormalSeeId:         constant.SwitchOff,
			IsNormalAddPrivilege:  constant.SwitchOn,
			IsAddNormalVerify:     constant.SwitchOn,
			IsAddPrivilegeVerify:  constant.SwitchOn,
			IsPrivilegeAddVerify:  constant.SwitchOff,
			IsNormalJoinGroup:     constant.SwitchOn,
			IsNormalMulSend:       constant.SwitchOff,
			IsShowRevoke:          constant.SwitchOn,
			IsPictureOpen:         constant.SwitchOn,
			IsOssOpen:             constant.SwitchOn,
			IsDisplayNicknameOpen: constant.SwitchOn,
			IsOpenNews:            constant.SwitchOn,
			IsOpenSign:            constant.SwitchOn,
			IsOpenPlay:            constant.SwitchOn,
			IsOpenOperator:        constant.SwitchOn,
			IsOpenShop:            constant.SwitchOn,
			IsOpenTopStories:      constant.SwitchOn,
			IsOpenWallet:          constant.SwitchOn,
			IsOpenPaymentCode:     constant.SwitchOn,
		}
	}

	feihuConfig, _ := configUsecase.ConfigUseCase.GetFeihuConfig()
	res["feihu_config"] = feihuConfig
	http.Success(ctx, res)
}

func (s *settingService) AboutUs(ctx *gin.Context) {
	var (
		err  error
		req  settingModel.AboutUsReq
		resp settingModel.AboutUsResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	resp.Content, err = configUsecase.ConfigUseCase.GetAboutUs()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetAboutUs, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}
	http.Success(ctx, resp)
}

func (s *settingService) PrivacyPolicy(ctx *gin.Context) {
	var (
		err  error
		req  settingModel.AboutUsReq
		resp settingModel.AboutUsResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	resp.Content, err = configUsecase.ConfigUseCase.GetPrivacyPolicy()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetAboutUs, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}
	http.Success(ctx, resp)
}

func (s *settingService) UserAgreement(ctx *gin.Context) {
	var (
		err  error
		req  settingModel.AboutUsReq
		resp settingModel.AboutUsResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	resp.Content, err = configUsecase.ConfigUseCase.GetUserAgreement()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetAboutUs, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}
	http.Success(ctx, resp)
}

func (s *settingService) Version(ctx *gin.Context) {
	var (
		err  error
		req  settingModel.VersionReq
		resp *settingModel.VersionResp
		app  *model.AppVersion
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	opt := settingRepo.WhereOptionByVersion{
		Platform: req.Platform,
		Status:   constant.SwitchOn,
	}
	if app, err = settingRepo.VersionRepo.GetVersionInfo(opt); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetVersionInfo, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrVersionNotExist, lang))
		return
	}
	if err = util.Copy(&app, &resp); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetVersionInfo copy, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}

	http.Success(ctx, resp)
}

func (s *settingService) SmsCode(ctx *gin.Context) {
	var (
		err error
		req settingModel.SmsReq
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind query, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if req.UsedFor == constant.PhoneNumberForChangeWd {

		if UserInfo, _ := usecase.UserUseCase.GetBaseInfoByPhoneNumber(req.PhoneNumber); UserInfo == nil {
			http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
			return
		}
	}

	if err = usecase.GenSmsCode(req.OperationID, req.PhoneNumber, req.CountryCode, "bao"); err != nil {
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}

func (s *settingService) GetDiscoverInfo(ctx *gin.Context) {
	var (
		err  error
		req  settingModel.OperationIDReq
		resp settingModel.GetDiscoverInfoResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind query, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	status, err := configRepo.ConfigRepo.GetDiscoverIsOpen()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db query, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}
	resp.IsOpen = status
	discovers, err := discoverRepo.DiscoverRepo.GetDiscovers()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db query, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}
	if err = util.CopyStructFields(&resp.List, &discovers); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("copy json, error: %v", err))
	}
	http.Success(ctx, resp)
}

func (s *settingService) GetShieldList(ctx *gin.Context) {
	var (
		err   error
		count int64
		req   settingModel.ShieldListReq
		resp  settingModel.ShieldListResp
		list  []model.ShieldWords
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	list, count, err = settingRepo.ShieldRepo.GetList(req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("getlist, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}
	if count > 0 {
		for _, words := range list {
			resp.List = append(resp.List, words.ShieldWords)
		}
	} else {
		resp.List = []string{}
	}
	resp.Count = count
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	http.Success(ctx, resp)
}

func (s *settingService) DomainList(ctx *gin.Context) {
	var (
		req  domainModel.DomainListReq
		resp domainModel.HttpDomainListReturn
		err  error
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.BindJSON(&req); err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	url := fmt.Sprintf("%s/domain/domain_list", config.Config.Cms.CtrlApi)
	var data []byte
	if data, err = http.Post(url, req, 30); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		return
	}
	logger.Sugar.Infof("return DomainList %s", string(data))
	_ = json.Unmarshal(data, &resp)
	if resp.Code != 0 {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", resp.Message)
		return
	}
	http.Success(ctx, resp.Resp)
}
