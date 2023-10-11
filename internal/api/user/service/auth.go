package service

import (
	"encoding/json"
	"fmt"
	"im/config"
	friendUseCase "im/internal/api/friend/usecase"
	groupUseCase "im/internal/api/group/usecase"
	settingModel "im/internal/api/setting/model"
	"im/pkg/response"
	"time"

	settingRepo "im/internal/api/setting/repo"
	"im/internal/api/user/model"
	userRepo "im/internal/api/user/repo"
	userUseCase "im/internal/api/user/usecase"
	configModel "im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	"strings"

	"errors"
	"github.com/gin-gonic/gin"
)

var AuthService = new(authService)

type authService struct{}

func (s *authService) Login(ctx *gin.Context) {
	var (
		err      error
		req      model.LoginReq
		resp     model.LoginResp
		userInfo *model.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if err = s.loginSettingRule(&req, lang); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("login setting, error: %v", err))
		http.Failed(ctx, err)
		return
	}
	logger.Sugar.Infof("输入的参数是:%+v", req)
	ip := ctx.ClientIP()
	switch req.LoginType {
	case 3:
		if userInfo, err = s.visitorUser(ctx, &req, ip); err != nil {
			if errors.Is(err, response.GetError(response.ErrRegisterTimeLimit, lang)) {
				http.Failed(ctx, response.GetError(response.ErrRegisterTimeLimit, lang))
				return
			}
			if errors.Is(err, response.GetError(response.ErrRegisterLimit, lang)) {
				http.Failed(ctx, response.GetError(response.ErrRegisterLimit, lang))
				return
			}
			http.Failed(ctx, response.GetError(response.ErrUnauthorized, lang))
			return
		}

	default:
		if userInfo, err = s.normalUser(ctx, &req, ip); err != nil {
			http.Failed(ctx, err)
			return
		}
	}

	station := config.Config.Station
	username := fmt.Sprintf("%s_%s", station, userInfo.UserID)
	var clients []mqtt.Client
	clients, err = mqtt.GetClients(username, "", "", "", 0, 0)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("login error, userId:%s GetClients, error: %v", userInfo.UserID, err))
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}

	if len(clients) == 0 {
		if err = userUseCase.AuthUseCase.CreateAuthUsername(userInfo.UserID); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("login error, create auth guest username error, error: %v", err))
			http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
			return
		}
	}

	if resp.Token, err = userRepo.UserCache.AddTokenCache(userInfo); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("auth take login token, error: %v", err))
		http.Failed(ctx, err)
		return
	}

	resp.UserId = userInfo.UserID

	http.Success(ctx, resp)
}

func (s *authService) checkUserRule(ctx *gin.Context, req *model.LoginReq, userId string, verPwd bool) (*model.User, error) {
	var (
		err      error
		userInfo *model.User
	)
	opt := userRepo.WhereOption{
		Account:     req.Account,
		PhoneNumber: req.PhoneNumber,
		UserId:      userId,
		CountryCode: req.CountryCode,
	}
	lang := ctx.GetHeader("Locale")
	userInfo, err = userRepo.UserRepo.GetByUserID(opt)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("auth GetByUserID , error: %v", err))
		if err == response.GetError(response.ErrUserIdNotExist, lang) && req.LoginType != constant.GuestStr && ctx.GetHeader(constant.ImSiteHeaderStr) != "" {
			if userInfo, err = s.imSiteAutoRegister(ctx, req); err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("第三方站点注册user失败, error: %v", err))
				return nil, err
			}
		}
		if userInfo == nil {
			if req.LoginType == constant.PhoneNumberStr {
				return nil, response.GetError(response.ErrBadPhoneNumPwd, lang)
			}
			return nil, response.GetError(response.ErrBadAccountPwd, lang)
		}
	}
	if userInfo.Status == constant.UserStatusFreeze {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("auth Status  , error: %v", err))
		return nil, response.GetError(response.ErrUserFreeze, lang)
	}
	if verPwd && !util.CheckPassword(userInfo.Password, req.Password, userInfo.Salt) {
		if req.LoginType == constant.PhoneNumberStr {
			return nil, response.GetError(response.ErrBadPhoneNumPwd, lang)
		}
		return nil, response.GetError(response.ErrBadAccountPwd, lang)
	}

	return userInfo, nil
}

func (s *authService) normalUser(ctx *gin.Context, req *model.LoginReq, ip string) (*model.User, error) {
	var (
		err      error
		userInfo *model.User
	)
	lang := ctx.GetHeader("Locale")
	if req.PhoneNumber == "" && req.Account == "" {
		return nil, response.GetError(response.ErrBadAccount, lang)
	}
	if req.LoginType == 2 {
		if len(req.Account) < constant.AccountLen || len(req.Account) > 16 {
			return nil, response.GetError(response.ErrBadAccount, lang)
		}
		req.Account = strings.ToLower(req.Account)
	} else {
		if len(req.PhoneNumber) < 7 {
			return nil, response.GetError(response.ErrBadPhoneNumber, lang)
		}
		if strings.HasPrefix(req.PhoneNumber, constant.ChinaCountryCode) {
			req.CountryCode = constant.ChinaCountryCode
			req.PhoneNumber = string([]byte(req.PhoneNumber)[3:])
		} else {
			if strings.HasPrefix(req.PhoneNumber, constant.PhoneNumberPrefix) {
				req.CountryCode = ""
			}
		}
	}
	if req.Password == "" || len(req.Password) < 6 || len(req.Password) > 16 {
		return nil, response.GetError(response.ErrBadPassword, lang)
	}
	if userInfo, err = s.checkUserRule(ctx, req, req.Account, true); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check user, error: %v", err))
		return nil, err
	}
	d := &model.User{
		Platform:        req.Platform,
		DeviceId:        req.DeviceId,
		LoginIp:         ip,
		LatestLoginTime: time.Now().Unix(),
	}
	opt := userRepo.WhereOption{
		Id: userInfo.ID,
	}
	_, err = userRepo.UserRepo.UpdateBy(opt, d)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("UpdateBy user, error: %v", err))
		return nil, response.GetError(response.ErrUnknown, lang)
	}
	userRepo.UserCache.DelUserInfoOnCache(userInfo.UserID)

	err = userRepo.UserRepo.RecordUserDeviceAndIP(userInfo.UserID, req.DeviceId, ip, req.Platform)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("RecordUserDeviceAndIP, error: %v", err))
	}

	err = userRepo.UserRepo.RecordUserLoginHistory(userInfo.UserID, req.DeviceId, ip, req.Brand, req.Platform)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("RecordUserLoginHistory, error: %v", err))
	}
	return userInfo, nil
}

func (s *authService) imSiteAutoRegister(ctx *gin.Context, req *model.LoginReq) (*model.User, error) {
	var (
		err error
	)
	lang := ctx.GetHeader("Locale")
	imSite := ctx.GetHeader(constant.ImSiteHeaderStr)
	regReq := model.RegisterReq{
		ImSite: strings.ToLower(imSite),
	}
	switch regReq.ImSite {
	case constant.ImSiteZhaoCai:
		if err = util.CopyStructFields(&regReq, req); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("CopyStructFields，error: %v", err))
			return nil, err
		}
		return userUseCase.NewRegisterUseCase(ctx, friendUseCase.FriendUseCase, groupUseCase.GroupUseCase).Register(regReq)
	}
	logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("自动注册失败，error: 暂不支持 %s 站点登录", imSite))
	return nil, response.GetError(response.ErrImSitePermissions, lang)
}

func (s *authService) visitorUser(ctx *gin.Context, req *model.LoginReq, ip string) (*model.User, error) {
	var (
		err        error
		user       *model.User
		userCheck  *model.User
		userResult *model.User
	)
	lang := ctx.GetHeader("Locale")

	guestOpt := userRepo.WhereOption{
		DeviceId:  req.DeviceId,
		Platform:  req.Platform,
		UserModel: constant.GuestUserModel,
	}
	userResult, _ = userRepo.UserRepo.GetByDeviceID(guestOpt)
	if userResult != nil {
		if userResult.Status == 2 {
			return nil, response.GetError(response.ErrUserFreeze, lang)
		}
		if userResult.Status != 1 {
			return nil, response.GetError(response.ErrUserPermissions, lang)
		}
	} else {
		regReq := model.RegisterReq{}
		regReq.DeviceId = req.DeviceId
		regReq.Platform = req.Platform
		err = userUseCase.NewRegisterUseCase(ctx, friendUseCase.FriendUseCase, groupUseCase.GroupUseCase).SystemConfigRule(&regReq)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("guest login error, create guest account is limit, error: %v", err))
			return nil, err
		}
		user = &model.User{
			UserModel: constant.GuestUserModel,
			Platform:  req.Platform,
			DeviceId:  req.DeviceId,
		}
		user.UserID = util.RandID(db.UserIDSize)
		user.Account = "G" + util.RandStringInt(9)
		user.Salt = util.RandString(6)
		user.Password = util.GetPassword(util.RandString(6), user.Salt)
		user.NickName = user.Account
		user.LoginIp = ip
		user.RegisterIp = ip
		user.RegisterDeviceId = req.DeviceId
		user.LatestLoginTime = time.Now().Unix()
		req.Account = user.Account
		if ctx.GetHeader(constant.ImSiteHeaderStr) != "" {
			user.ImSite = ctx.GetHeader(constant.ImSiteHeaderStr)
		}
		if userCheck, err = s.checkUserRule(ctx, req, user.UserID, false); userCheck != nil {
			return nil, response.GetError(response.ErrUserIdExist, lang)
		}
		if err != nil {
			logger.Sugar.Warnw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check guest user, error: %v", err))
		}

		if err = userUseCase.AuthUseCase.CreateAuthUsername(user.UserID); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("login error, create auth guest username error, error: %v", err))
			return nil, err
		}

		if userResult, err = userRepo.UserRepo.Create(user); err != nil {
			return nil, err
		}

		friendUseCase.FriendUseCase.CreateFriendLabel(userResult.UserID, userResult.UserID, "我的好友")

		userRepo.UserCache.RecordRegisterDeviceIDCount(userResult.DeviceId)
		userRepo.UserCache.RecordRegisterIPCount(userResult.RegisterIp)
	}

	s.joinDefaultGroups(req.OperationID, userResult.UserID)
	s.addDefaultFriends(req.OperationID, userResult.UserID)

	err = userRepo.UserRepo.RecordUserDeviceAndIP(userResult.UserID, userResult.RegisterDeviceId, userResult.LoginIp, userResult.Platform)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("RecordUserDeviceAndIP, error: %v", err))
	}

	err = userRepo.UserRepo.RecordUserLoginHistory(userResult.UserID, req.DeviceId, ip, req.Brand, req.Platform)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("RecordUserLoginHistory, error: %v", err))
	}
	return userResult, nil
}

func (s *authService) loginSettingRule(req *model.LoginReq, lang string) error {
	var (
		err error
		g   map[string][]int64
		c   *settingModel.SettingConfig
	)
	if req.LoginType != 3 && req.Password == "" {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "login params error:非游客模式下密码必须输入")
		return response.GetError(response.ErrBadPassword, lang)
	}
	opt := settingRepo.WhereOption{
		ConfigType: []string{"login_config"},
	}

	if c, err = settingRepo.SettingRepo.GetInfo(opt); c != nil {
		if c.Content != "" {
			var ig []int64
			if err = json.Unmarshal([]byte(c.Content), &g); err != nil {
				return err
			}
			for _, v := range g {
				ig = append(ig, v...)
			}
			if !util.InSliceForInt64(ig, req.LoginType) {
				return response.GetError(response.ErrLoginDenied, lang)
			}
		}
	}
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "GetInfo error:", err)
	}

	return nil
}

func (s *authService) joinDefaultGroups(operationId, userId string) {
	go func() {
		groupList := groupUseCase.GroupUseCase.GetDefaultGroups()
		for _, g := range groupList {
			if err := groupUseCase.GroupUseCase.JoinGroup(operationId, g.GroupId, userId); err != nil {
				logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" user joinDefaultGroups, error: %v", err))
				continue
			}
		}
	}()
}

func (s *authService) addDefaultFriends(operationId, userId string) {
	var (
		err error
		fin *configModel.DefaultFriend
	)
	go func() {
		if fin, err = configRepo.DefaultFriendRepo.RandFriend(); err != nil {
			logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" user addDefaultFriends, error: %v", err))
			return
		}
		if err = friendUseCase.FriendUseCase.AddFriend(operationId, userId, fin.UserId, "", fin.GreetMsg, false); err != nil {
			logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" user AddFriend, error: %v", err))
			return
		}
	}()
}
