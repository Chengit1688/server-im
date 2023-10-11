package service

import (
	"encoding/json"
	"fmt"
	"github.com/u2takey/go-utils/rand"
	"gorm.io/gorm"
	"im/config"
	friendUseCase "im/internal/api/friend/usecase"
	groupUseCase "im/internal/api/group/usecase"
	settingModel "im/internal/api/setting/model"
	settingRepo "im/internal/api/setting/repo"
	userModel "im/internal/api/user/model"
	userRepo "im/internal/api/user/repo"
	"im/internal/api/user/usecase"
	"im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	configUsecase "im/internal/cms_api/config/usecase"
	discoverModel "im/internal/cms_api/discover/model"
	discoverRepo "im/internal/cms_api/discover/repo"
	operationModel "im/internal/cms_api/operation/model"
	operationRepo "im/internal/cms_api/operation/repo"
	"im/internal/cms_api/user/repo"
	cmsWalletRepo "im/internal/cms_api/wallet/repo"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var UserService = new(userService)

type userService struct{}

func (s *userService) RegisterV1(ctx *gin.Context) {
	var (
		err        error
		req        userModel.RegisterReq
		resp       userModel.RegisterResp
		user       *userModel.User
		userResult *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json , error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "req", req)

	if user, err = s.registerRule(&req, lang); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("register, error: %v", err))
		http.Failed(ctx, err)
		return
	}

	opt := userRepo.WhereOption{
		UserId:      user.UserID,
		Account:     user.Account,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
	}
	if b := userRepo.UserRepo.OrExists(opt); b {
		if req.AccountType == constant.PhoneNumberStr {
			http.Failed(ctx, response.GetError(response.ErrPhoneNumberExist, lang))
			return
		}
		http.Failed(ctx, response.GetError(response.ErrAccountExist, lang))
		return
	}

	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "info", "emqx create auth username")
	if err = usecase.AuthUseCase.CreateAuthUsername(user.UserID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("register error, create auth username error, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrRegisterFailed, lang))
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "info", "emqx create auth username success")

	user.Salt = util.RandString(6)
	user.Password = util.GetPassword(req.Password, user.Salt)
	user.InviteCode = req.InviteCode
	user.ImSite = constant.ImSiteIm
	if userResult, err = userRepo.UserRepo.Create(user); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error creating registered user, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrRegisterFailed, lang))
		return
	}
	if resp.Token, err = s.registerToken(userResult); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error creating registered token, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrRegisterFailed, lang))
		return
	}

	resp.UserId = userResult.UserID
	s.defaultInviteHandle(req.OperationID, req.InviteCode, userResult.UserID)

	http.Success(ctx, resp)
}

func (s *userService) Register(ctx *gin.Context) {
	var (
		err        error
		req        userModel.RegisterReq
		resp       userModel.RegisterResp
		userResult *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json , error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "req", req)
	if userResult, err = usecase.NewRegisterUseCase(ctx, friendUseCase.FriendUseCase, groupUseCase.GroupUseCase).Register(req); err != nil {
		http.Failed(ctx, err)
		return
	}
	if resp.Token, err = s.registerToken(userResult); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error creating registered token, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrRegisterFailed, lang))
		return
	}
	resp.UserId = userResult.UserID
	ip := ctx.ClientIP()

	err = userRepo.UserRepo.RecordUserLoginHistory(userResult.UserID, userResult.DeviceId, ip, req.Brand, req.Platform)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("RecordUserLoginHistory, error: %v", err))
	}
	http.Success(ctx, resp)
}

func (s *userService) registerToken(user *userModel.User) (string, error) {
	var (
		userInfo *userModel.User
		err      error
		token    string
	)
	opt := userRepo.WhereOption{
		Id: user.ID,
	}
	if userInfo, err = userRepo.UserRepo.GetByUserID(opt); err != nil {
		return "", err
	}

	if token, err = userRepo.UserCache.AddTokenCache(userInfo); err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) registerRule(req *userModel.RegisterReq, lang string) (*userModel.User, error) {
	var (
		err error
		c   *settingModel.SettingConfig
		reg settingModel.RegisterConfigInfo
	)
	if req.AccountType == constant.PhoneNumberStr {
		if req.PhoneNumber == "" {
			return nil, response.GetError(response.ErrBadPhoneNumber, lang)
		}

		if req.CountryCode == "" && !strings.HasPrefix(req.PhoneNumber, constant.PhoneNumberPrefix) {
			return nil, response.GetError(response.ErrBadPhoneNumber, lang)
		}
		if strings.HasPrefix(req.PhoneNumber, constant.PhoneNumberPrefix) {
			req.CountryCode = ""
		}
	}
	if req.AccountType == constant.AccountStr && req.Account == "" {
		return nil, response.GetError(response.ErrBadAccount, lang)
	}
	opt := settingRepo.WhereOption{
		ConfigType: []string{"register_config"},
	}

	if c, err = settingRepo.SettingRepo.GetInfo(opt); c != nil {
		if c.Content != "" {
			if err = json.Unmarshal([]byte(c.Content), &reg); err != nil {
				return nil, err
			}
		}
	}
	if err != nil {
		logger.Sugar.Warnw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetInfo, error: %v", err))
	}
	u := &userModel.User{
		PhoneNumber: req.PhoneNumber,
		Account:     req.Account,
		UserID:      util.RandID(db.UserIDSize),
		CountryCode: req.CountryCode,
		NickName:    util.RandString(10),
	}
	if req.AccountType == constant.AccountStr {
		u.Account = strings.ToLower(req.Account)
		u.NickName = u.Account
	}
	if req.AccountType == constant.PhoneNumberStr && reg.IsSmsCode == constant.SwitchOn {
		if err = s.smsVerify(req.PhoneNumber, req.SmsCode, lang); err != nil {
			return nil, err
		}
	}
	if reg.IsInviteCode == constant.SwitchOn {
		if req.InviteCode == "" && reg.CheckInviteCode == constant.SwitchOn {
			return nil, response.GetError(response.ErrInviteCode, lang)
		}
		u.InviteCode = req.InviteCode
		if !s.inviteCodeVerify(req.InviteCode) {
			u.InviteCode = ""
			if reg.CheckInviteCode == constant.SwitchOn {
				return nil, response.GetError(response.ErrBadInviteCode, lang)
			}
		}
	}
	if reg.IsVerificationCode == constant.SwitchOn {
		if req.VerificationToken == "" || req.VerificationPoint == "" || req.CaptchaType == "" {
			return nil, response.GetError(response.ErrVerificationCode, lang)
		}
		if err = usecase.GetCaptchaFactory().GetService(req.CaptchaType).Verification(req.VerificationToken, req.VerificationPoint); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (s *userService) smsVerify(PhoneNumber string, smsCode string, lang string) error {
	if smsCode == config.Config.Sms.SuperCode {
		return nil
	}
	var (
		cacheCode string
		err       error
	)
	if cacheCode, err = settingRepo.SettingCache.GetAccountCode(PhoneNumber); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("smsVerify GetAccountCode, error: %v", err))
		return response.GetError(response.ErrBadCode, lang)
	}
	if smsCode != cacheCode {
		return response.GetError(response.ErrBadCode, lang)
	}
	_, _ = settingRepo.SettingCache.DelAccountCode(PhoneNumber)

	return nil
}

func (s *userService) inviteCodeVerify(inviteCode string) bool {
	var (
		count int64
		err   error
	)
	if inviteCode == "" {
		return true
	}
	opt := configRepo.WhereOptionForInvite{
		InviteCode:   inviteCode,
		Status:       1,
		DeleteStatus: 2,
	}
	if count, err = configRepo.InviteCode.Exists(opt); count > 0 {
		return true
	}
	if err != nil {
		logger.Sugar.Warnw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user invite code error, error: %v", err))
	}

	return false
}

func (s *userService) defaultInviteHandle(operationID string, inviteCode string, userId string) bool {
	var (
		cInfo *model.InviteCode
		err   error
	)
	if inviteCode != "" {
		opt := configRepo.WhereOptionForInvite{
			InviteCode:   inviteCode,
			Status:       1,
			DeleteStatus: 2,
		}
		cInfo, err = configRepo.InviteCode.GetInfo(opt)
		if err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user invite code, error: %v", err))
		}
	}
	s.joinDefaultGroups(operationID, cInfo, userId)
	s.addDefaultFriends(operationID, cInfo, userId)

	return true
}

func (s *userService) joinDefaultGroups(operationID string, inviteCode *model.InviteCode, userId string) {
	var err error
	go func() {

		if inviteCode != nil {
			groupList := strings.Split(inviteCode.DefaultGroups, ",")
			for _, groupId := range groupList {
				if err = groupUseCase.GroupUseCase.JoinGroup(operationID, groupId, userId); err != nil {
					logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user joinDefaultGroups, error: %v", err))
					continue
				}
			}
			return
		}

		groupList := groupUseCase.GroupUseCase.GetDefaultGroups()
		for _, g := range groupList {
			if err = groupUseCase.GroupUseCase.JoinGroup(operationID, g.GroupId, userId); err != nil {
				logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user joinDefaultGroups, error: %v", err))
				continue
			}
		}
	}()
}

func (s *userService) addDefaultFriends(operationID string, inviteCode *model.InviteCode, userId string) {
	var (
		err error
		fin *model.DefaultFriend
	)
	go func() {

		if inviteCode != nil {
			friendList := strings.Split(inviteCode.DefaultFriends, ",")
			for _, friendId := range friendList {

				if err = friendUseCase.FriendUseCase.AddFriend(operationID, userId, friendId, "", inviteCode.GreetMsg, false); err != nil {
					logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user addDefaultFriends, error: %v", err))
					continue
				}
			}
			return
		}

		if fin, err = configRepo.DefaultFriendRepo.RandFriend(); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user addDefaultFriends, error: %v", err))
			return
		}

		if err = friendUseCase.FriendUseCase.AddFriend(operationID, userId, fin.UserId, "", fin.GreetMsg, false); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user addDefaultFriends, error: %v", err))
			return
		}
	}()
}

func (s *userService) UpdateInfo(ctx *gin.Context) {
	var (
		err      error
		req      userModel.UserInfoUpdateReq
		userInfo *userModel.User
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	loginUser := ctx.GetString("user_id")
	lang := ctx.GetHeader("Locale")

	if req.NickName != nil {
		nickName := *req.NickName
		optsNickname := new(userRepo.WhereOption)
		optsNickname.NickName = nickName
		dbUser, err := userRepo.UserRepo.GetInfo(*optsNickname)
		if err == nil {
			if dbUser.UserID != loginUser {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "昵称已存在", "nickname", optsNickname.NickName)
				http.Failed(ctx, response.GetError(response.ErrNickNameUsed, lang))
				return
			}
		}
	}

	if err = userRepo.UserRepo.UserInfoUpdate(loginUser, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update user error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}

	opt := userRepo.WhereOption{
		UserId: loginUser,
	}

	userRepo.UserCache.DelUserInfoOnCache(opt.UserId)
	if userInfo, err = userRepo.UserRepo.GetByUserID(opt); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetByUserID error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
		return
	}
	updateU := userModel.UserBaseInfo{
		UserId:      userInfo.UserID,
		Account:     userInfo.Account,
		FaceURL:     userInfo.FaceURL,
		BigFaceURL:  userInfo.BigFaceURL,
		Gender:      userInfo.Gender,
		NickName:    userInfo.NickName,
		Signatures:  userInfo.Signatures,
		Age:         userInfo.Age,
		IsPrivilege: userInfo.IsPrivilege,
		PhoneNumber: userInfo.PhoneNumber,
		CountryCode: userInfo.CountryCode,
	}

	if err = mqtt.SendMessageToUsers(req.OperationID, common.UserInfoPush, updateU, userInfo.UserID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("SendMessageToOnlineUsers error: %v", err))
	}
	s.runRoutineByUserId(req.OperationID, userInfo.UserID)

	http.Success(ctx)
}

func (s *userService) runRoutineByUserId(OperationID, userId string) {
	rm := []usecase.FuncForUserId{
		groupUseCase.GroupMemberUseCase.UserInfoChange,
		friendUseCase.FriendUseCase.UserInfoChange,
	}
	usecase.UserUseCase.DoRoutineByUserId(OperationID, userId, rm...)
}

func (s *userService) SignToday(ctx *gin.Context) {
	var (
		err error
		req userModel.SignTodayReq
		g   map[string]int64
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json  error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if g, err = s.getSignConfig(); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("SignToday error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}
	if g["sign"] == constant.SwitchOff {
		http.Failed(ctx, response.GetError(response.ErrUserPermissions, lang))
		return
	}

	resp, err := configRepo.ConfigRepo.GetSystemConfig()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetParameterConfig error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}
	userId := ctx.GetString("user_id")
	award := 0

	user, err := userRepo.UserRepo.GetInfo(userRepo.WhereOption{UserId: userId})
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrUserNotFound, lang))
		return
	}

	if resp.SignAward > 0 {
		SignRewardKey := fmt.Sprintf(cmsWalletRepo.SignReward, userId)
		SignRewardLock := util.NewLock(db.RedisCli, SignRewardKey)
		if err = SignRewardLock.Lock(); err != nil {
			return
		}
		defer SignRewardLock.Unlock()

		withdraw, err := configUsecase.ConfigUseCase.GetWithdrawConfig()
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
			return
		}
		tempAward := withdraw.Min - user.Balance
		if tempAward >= resp.SignAward {
			award = rand.Intn(int(resp.SignAward))
		} else if tempAward < resp.SignAward && tempAward > 100 {
			award = rand.Intn(int(tempAward))
		} else {
			award = 1
		}
	}

	if err = repo.SignLogRepo.CreateSignLog(ctx.GetString("user_id"), user.Balance, int64(award)); err != nil {
		http.Failed(ctx, err)
		return
	}
	if err = userRepo.UserRepo.UpdateWallet(userId, "+", int64(award)); err != nil {
		http.Failed(ctx, err)
		return
	}
	http.Success(ctx, award)
}

func (s *userService) getSignConfig() (map[string]int64, error) {
	var (
		err error
		c   *settingModel.SettingConfig
		g   map[string]int64
	)
	opt := settingRepo.WhereOption{
		ConfigType: []string{"sign_log_config"},
	}
	c, _ = settingRepo.SettingRepo.GetInfo(opt)
	if c == nil {
		c = &settingModel.SettingConfig{
			Content: `{"sign":1}`,
		}
	}
	if err = json.Unmarshal([]byte(c.Content), &g); err != nil {
		return nil, err
	}

	return g, nil
}

func (s *userService) SignList(ctx *gin.Context) {
	var (
		err         error
		OperationID string
		g           map[string]int64
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&OperationID); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	total, list := repo.SignLogRepo.GetMonthSignLogByUserId(ctx.GetString("user_id"))
	resp := &userModel.GetUserSignInfoResp{
		Total:    total,
		Today:    false,
		SignOpen: false,
		Days:     list,
	}
	for _, v := range list {
		if v == time.Now().Day() {
			resp.Today = true
			break
		}
	}
	if g, err = s.getSignConfig(); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("SignList, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}
	if g["sign"] == constant.SwitchOn {
		resp.SignOpen = true
	}

	http.Success(ctx, resp)
}

func (s *userService) SignListWeek(ctx *gin.Context) {
	var (
		err         error
		OperationID string
		g           map[string]int64
		user        *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&OperationID); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if user, err = userRepo.UserRepo.GetByUserID(userRepo.WhereOption{UserId: ctx.GetString("user_id")}); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}

	p, err1 := configRepo.ConfigRepo.GetSystemConfig()
	if err1 != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", err1)
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}
	total, list := repo.SignLogRepo.GetMonthSignLogByUserIdV2(ctx.GetString("user_id"))
	resp := &userModel.GetUserSignInfoV2Resp{
		Total:     total,
		Today:     false,
		SignOpen:  false,
		Days:      list,
		Balance:   user.Balance,
		SignAward: p.SignAward,
	}
	for _, v := range list {
		if v == time.Now().Format("2006-01-02") {
			resp.Today = true
			break
		}
	}
	if g, err = s.getSignConfig(); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}
	if g["sign"] == constant.SwitchOn {
		resp.SignOpen = true
	}

	http.Success(ctx, resp)
}

func (s *userService) Info(ctx *gin.Context) {
	var (
		err  error
		req  userModel.UserInfoReq
		resp *userModel.UserInfoResp
		user *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	opt := userRepo.WhereOption{
		UserId: ctx.GetString("user_id"),
	}
	if user, err = userRepo.UserRepo.GetByUserID(opt); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user login, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
		return
	}
	if err = util.Copy(user, &resp); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("copy struct, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}

	http.Success(ctx, resp)
}

func (s *userService) DeviceList(ctx *gin.Context) {
	var (
		err         error
		OperationID string
		resp        []userModel.GetDeviceInfo
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.BindQuery(&OperationID); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	for i, v := range constant.PlatformID2Name {
		resp = append(resp, userModel.GetDeviceInfo{
			PlatformClass: constant.PlatformNameToClass(v),
			DeviceId:      int64(i),
			DeviceName:    v,
		})
	}

	http.Success(ctx, resp)
}

func (s *userService) ForgotPassword(ctx *gin.Context) {
	var (
		err  error
		req  userModel.ForgotPasswordReq
		user *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	opt := userRepo.WhereOption{
		PhoneNumber: req.PhoneNumber,
		CountryCode: req.CountryCode,
		Status:      constant.SwitchOn,
	}
	if user, err = userRepo.UserRepo.GetByUserID(opt); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetByUserID user password, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
		return
	}

	data := &userModel.User{}
	opt = userRepo.WhereOption{
		Id: user.ID,
	}
	data.Password = util.GetPassword(req.NewPassword, user.Salt)
	if _, err = userRepo.UserRepo.UpdateBy(opt, data); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("UpdateBy user update password, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
		return
	}

	http.Success(ctx)
}

func (s *userService) VerifyPhoneCode(ctx *gin.Context) {
	var (
		err error
		req userModel.VerifyPhoneCodeReq
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if err = s.smsVerify(req.PhoneNumber, req.SmsCode, lang); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("sms error, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadCode, lang))
		return
	}

	http.Success(ctx)
}

func (s *userService) PasswordSecure(ctx *gin.Context) {
	var (
		err  error
		req  userModel.PasswordSecureReq
		user *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	opt := userRepo.WhereOption{
		UserId: ctx.GetString("user_id"),
	}
	if user, err = userRepo.UserRepo.GetByUserID(opt); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetByUserID user update password, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
		return
	}
	data := &userModel.User{}
	switch req.PasswordType {
	case 1:
		if req.OriginalPassword != "" {
			if !util.CheckPassword(user.Password, req.OriginalPassword, user.Salt) {
				http.Failed(ctx, response.GetError(response.ErrWrongPassword, lang))
				return
			}
		}
		data.Password = util.GetPassword(req.NewPassword, user.Salt)
	case 2:
		if req.OriginalPassword == "" {
			http.Failed(ctx, response.GetError(response.ErrWrongPassword, lang))
			return
		}
		if !util.CheckPassword(user.PayPassword, req.OriginalPassword, user.Salt) {
			http.Failed(ctx, response.GetError(response.ErrWrongPassword, lang))
			return
		}
		data.PayPassword = util.GetPassword(req.NewPassword, user.Salt)
	case 3:
		if req.OriginalPassword == "" {
			http.Failed(ctx, response.GetError(response.ErrWrongPassword, lang))
			return
		}
		if !util.CheckPassword(user.Password, req.OriginalPassword, user.Salt) {
			http.Failed(ctx, response.GetError(response.ErrWrongPassword, lang))
			return
		}
		data.PayPassword = util.GetPassword(req.NewPassword, user.Salt)
	default:
		http.Failed(ctx, response.GetError(response.ErrConfigPassword, lang))
		return
	}
	opt = userRepo.WhereOption{
		UserId: user.UserID,
	}
	if _, err = userRepo.UserRepo.UpdateBy(opt, data); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("UpdateBy, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUpdateAccount, lang))
		return
	}

	http.Success(ctx)
}

func (s *userService) GetUserBaseInfo(ctx *gin.Context) {
	var (
		err      error
		req      userModel.GetUserInfoReq
		userInfo *userModel.UserBaseInfo
		resp     *userModel.UserBaseInfoResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if userInfo, err = userRepo.UserCache.GetBaseUserInfo(req.UserId, userRepo.UserRepo.GetBaseInfoByUserId); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("from cache, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrLoadUserInfo, lang))
		return
	}
	if err = util.Copy(userInfo, &resp); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("Copy, error: %v", err))
	}
	http.Success(ctx, resp)
}

func (s *userService) GetUserInfo(ctx *gin.Context) {
	var (
		err  error
		req  userModel.GetUserInfoReq
		resp *userModel.UserInfoResp
		user *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if user, err = userRepo.UserRepo.GetByUserID(userRepo.WhereOption{UserId: req.UserId}); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetByUserID, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if err = util.Copy(user, &resp); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("copy, error: %v", err))
	}

	http.Success(ctx, resp)
}

func (s *userService) GetServerVersion(ctx *gin.Context) {
	var (
		err  error
		req  userModel.GetServerVersionReq
		resp userModel.GetServerVersionRsp
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	user_id := ctx.GetString("user_id")

	groups := groupUseCase.GroupUseCase.GetGroupVersions(user_id)
	resp.GroupsVerison = make([]userModel.GroupVersionInfo, 0)
	for _, info := range groups {
		group := userModel.GroupVersionInfo{}
		util.CopyStructFields(&group, info)
		resp.GroupsVerison = append(resp.GroupsVerison, group)
	}

	http.Success(ctx, resp)
}

func (s *userService) UserConfigHandle(ctx *gin.Context) {
	var (
		err        error
		req        userModel.UserConfigHandleReq
		userConfig userModel.UserConfig
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	userConfig = userModel.UserConfig{
		Content: req.Content,
		UserId:  ctx.GetString("user_id"),
	}
	if userConfig, err = userRepo.UserConfigRepo.CreateOrUpdate(userConfig); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("CreateOrUpdate, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}
	pushConfig := userModel.UserConfigPushResp{
		Version: userConfig.Version,
		Content: userConfig.Content,
	}

	if err = mqtt.SendMessageToUsers(req.OperationID, common.ConfigUserChangePush, pushConfig, userConfig.UserId); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" user config push error: %v", err))
	}

	http.Success(ctx)
}

func (s *userService) GetUserConfig(ctx *gin.Context) {
	var (
		err         error
		OperationID string
		resp        userModel.GetUserConfigResp
		uf          userModel.UserConfig
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&OperationID); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	opt := userRepo.WhereOptionForUserConfig{
		UserId: ctx.GetString("user_id"),
	}
	if uf, err = userRepo.UserConfigRepo.GetUserConfig(opt); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetUserConfig, error: %v", err))

		http.Success(ctx, resp)
		return
	}
	if err = util.CopyStructFields(&resp, uf); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("copy, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUnknown, lang))
		return
	}

	http.Success(ctx, resp)
}

func (s *userService) GetUserOnlineStatus(ctx *gin.Context) {
	var (
		err  error
		req  userModel.GetUserOnlineStatusReq
		resp userModel.GetUserOnlineStatusResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.BindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}

	resp.UserId = req.UserId
	username := fmt.Sprintf("%s_%s", config.Config.Station, req.UserId)
	onlineClients, err := mqtt.GetClients(username, "", "", mqtt.ConnStateTypeConnected, 1, 1)
	if err != nil {
		user, err := usecase.UserUseCase.GetInfo(req.UserId)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user GetInfo, error: %v", err))
			http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
			return
		}
		resp.Online = false
		timeLayout := "2006/01/02 15:04:05"
		timeStr := time.Unix(user.LatestLoginTime, 0).Format(timeLayout)
		resp.OfflineInfo = fmt.Sprintf("离线 %s", timeStr)
	} else {
		user, _ := usecase.UserUseCase.GetInfo(req.UserId)
		for i := len(onlineClients) - 1; i >= 0; i-- {
			v := onlineClients[i]
			temp := strings.Split(v.ClientID, "_")
			if len(temp) < 2 {
				continue
			}
			resp.Online = true

			resp.Ip = user.LoginIp
			resp.IpAddress, _ = util.QueryIpRegion(resp.Ip)

			break
		}
		if !resp.Online {
			user, err := usecase.UserUseCase.GetInfo(req.UserId)
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user GetInfo, error: %v", err))
				http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
				return
			}
			resp.Online = false
			timeLayout := "2006/01/02 15:04:05"
			timeStr := time.Unix(user.LatestLoginTime, 0).Format(timeLayout)
			resp.OfflineInfo = fmt.Sprintf("离线 %s", timeStr)
		}
	}

	http.Success(ctx, resp)
}

func (s *userService) Suggestion(ctx *gin.Context) {
	var (
		err  error
		req  userModel.SuggestionReq
		user *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	opt := userRepo.WhereOption{
		UserId: ctx.GetString("user_id"),
	}
	if user, err = userRepo.UserRepo.GetInfo(opt); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetInfo, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrFailRequest, lang))
		return
	}
	suggestion := operationModel.Suggestion{
		UserID:     user.UserID,
		Account:    user.Account,
		NickName:   user.NickName,
		Content:    req.Content,
		Brand:      req.Brand,
		Platform:   req.Platform,
		AppVersion: req.AppVersion,
	}
	if _, err = operationRepo.SuggestionRepo.Create(suggestion); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("Create, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrFailRequest, lang))
		return
	}

	http.Success(ctx)
}

func (s *userService) BindPhone(ctx *gin.Context) {
	var (
		err      error
		req      userModel.VerifyPhoneCodeReq
		userInfo *userModel.User
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if err = s.smsVerify(req.PhoneNumber, req.SmsCode, lang); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("sms error, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadCode, lang))
		return
	}

	opt := userRepo.WhereOption{
		UserId: ctx.GetString("user_id"),
	}
	d := &userModel.User{
		PhoneNumber: req.PhoneNumber,
		CountryCode: req.CountryCode,
	}
	if _, err = userRepo.UserRepo.UpdateBy(opt, d); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update user error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userRepo.UserCache.DelUserInfoOnCache(opt.UserId)
	if userInfo, err = userRepo.UserRepo.GetByUserID(opt); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetByUserID error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrUserIdNotExist, lang))
		return
	}
	updateU := userModel.UserBaseInfo{
		UserId:      userInfo.UserID,
		Account:     userInfo.Account,
		FaceURL:     userInfo.FaceURL,
		BigFaceURL:  userInfo.BigFaceURL,
		Gender:      userInfo.Gender,
		NickName:    userInfo.NickName,
		Signatures:  userInfo.Signatures,
		Age:         userInfo.Age,
		IsPrivilege: userInfo.IsPrivilege,
		PhoneNumber: userInfo.PhoneNumber,
		CountryCode: userInfo.CountryCode,
	}

	if err = mqtt.SendMessageToUsers(req.OperationID, common.UserInfoPush, updateU, userInfo.UserID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("SendMessageToOnlineUsers error: %v", err))
	}

	http.Success(ctx)
}

func (s *userService) FavoriteImagePaging(ctx *gin.Context) {
	var (
		err    error
		req    userModel.GetFavoriteImageReq
		resp   userModel.GetFavoriteImageResp
		images []userModel.FavoriteImage
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := ctx.GetString("user_id")
	req.Check()
	images, resp.Count, err = userRepo.UserRepo.FavoriteImagePaging(userID, req.Offset, req.Limit)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db query, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	util.CopyStructFields(&resp.List, images)
	if resp.List == nil {
		resp.List = make([]userModel.FavoriteImageItem, 0)
	}
	http.Success(ctx, resp)
}

func (s *userService) FavoriteImageAdd(ctx *gin.Context) {
	var (
		err   error
		req   userModel.AddFavoriteImageReq
		resp  userModel.AddFavoriteImageResp
		image userModel.FavoriteImage
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	userID := ctx.GetString("user_id")
	check := userModel.FavoriteImage{
		UserID: userID,
		UUID:   req.UUID,
	}
	exist := userRepo.UserRepo.FavoriteImageExist(check)
	if exist {
		http.Failed(ctx, response.GetError(response.ErrUserFavImageExist, lang))
		return
	}
	image = userModel.FavoriteImage{
		UserID:         userID,
		UUID:           req.UUID,
		ImageUrl:       req.ImageUrl,
		ImageThumbnail: req.ImageThumbnail,
		ImageWidth:     &req.ImageWidth,
		ImageHeight:    &req.ImageHeight,
	}

	err = userRepo.UserRepo.FavoriteImageAdd(image)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db insert, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}

	if err = mqtt.SendMessageToUsers(req.OperationID, common.UserFavoriteImagePush, userID, userID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("SendMessageToOnlineUsers error: %v", err))
	}

	util.CopyStructFields(&resp, image)
	http.Success(ctx, resp)
}

func (s *userService) FavoriteImageDel(ctx *gin.Context) {
	var (
		err error
		req userModel.DelFavoriteImageReq
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	userID := ctx.GetString("user_id")

	err = userRepo.UserRepo.FavoriteImageDel(userModel.FavoriteImage{
		UserID: userID,
		UUID:   req.UUID,
	})
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db delete, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
	}

	if err = mqtt.SendMessageToUsers(req.OperationID, common.UserFavoriteImagePush, userID, userID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("SendMessageToOnlineUsers error: %v", err))
	}
	http.Success(ctx, "")
}

func (s *userService) SetPrivacy(ctx *gin.Context) {
	var (
		err error
		req userModel.SetPrivacyReq
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	userID := ctx.GetString("user_id")
	jsonBytes, _ := json.Marshal(req.Data)
	if _, err = userRepo.UserRepo.UpdateBy(userRepo.WhereOption{UserId: userID}, &userModel.User{
		Privacy: string(jsonBytes),
	}); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db delete, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
	}

	http.Success(ctx)
}

func (s *userService) GetPrivacy(ctx *gin.Context) {
	var (
		err  error
		user *userModel.User
		resp userModel.Privacy
	)
	lang := ctx.GetHeader("Locale")
	userID := ctx.GetString("user_id")
	defaultByte, _ := util.MarshalJSONByDefault(&userModel.Privacy{}, true)
	privacy := string(defaultByte.([]byte))
	if user, err = userRepo.UserRepo.GetByUserID(userRepo.WhereOption{UserId: userID}); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
	}
	if user.Privacy != "" {
		privacy = user.Privacy
	}
	if err = json.Unmarshal([]byte(privacy), &resp); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("make json, error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
	}

	http.Success(ctx, resp)
}

func (s *userService) RealName(ctx *gin.Context) {
	var (
		err  error
		user userModel.User
		req  userModel.RealNameReq
		resp userModel.RealNameResp
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	lang := ctx.GetHeader("Locale")
	userID := ctx.GetString("user_id")
	if user, err = userRepo.UserRepo.UpdateRealName(userID, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}

	resp = userModel.RealNameResp{
		RealName:    user.RealName,
		IDNo:        user.IDNo,
		IDFrontImg:  user.IDFrontImg,
		IDBackImg:   user.IDBackImg,
		RealAuthMsg: user.RealAuthMsg,
		IsRealAuth:  user.IsRealAuth,
	}
	http.Success(ctx, resp)
}

func (s *userService) RealNameInfo(ctx *gin.Context) {
	var (
		err  error
		req  userModel.RealNameInfoReq
		user userModel.User
		resp userModel.RealNameResp
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	lang := ctx.GetHeader("Locale")
	userID := ctx.GetString("user_id")
	if user, err = userRepo.UserRepo.GetRealName(userID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}

	resp = userModel.RealNameResp{
		RealName:    user.RealName,
		IDNo:        user.IDNo,
		IDFrontImg:  user.IDFrontImg,
		IDBackImg:   user.IDBackImg,
		RealAuthMsg: user.RealAuthMsg,
		IsRealAuth:  user.IsRealAuth,
	}
	http.Success(ctx, resp)
}

func (s *userService) RedeemPrize(ctx *gin.Context) {
	var (
		err   error
		req   userModel.RedeemPrizeReq
		user  *userModel.User
		prize discoverModel.PrizeList
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	userID := ctx.GetString("user_id")
	if user, err = userRepo.UserRepo.GetByUserID(userRepo.WhereOption{UserId: userID}); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrUserNotFound, lang))
		return
	}

	key := userRepo.UserCache.GetPrizeKey(userID)
	l := util.NewLock(db.RedisCli, key)
	if err = l.Lock(); err != nil {
		return
	}
	defer l.Unlock()

	if user.Balance <= 0 {
		http.Failed(ctx, response.GetError(response.ErrBalanceNotEnough, lang))
		return
	}
	if prize, err = discoverRepo.DiscoverRepo.FetchPrize(req.PrizeID); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if prize.IsFictitious == 2 && (req.UserName == "" || req.Address == "") {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if user.Balance < prize.Cost {
		http.Failed(ctx, response.GetError(response.ErrBalanceNotEnough, lang))
		return
	}
	if err = discoverRepo.DiscoverRepo.AddRedeemPrize(userID, req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}

	http.Success(ctx)
}

func (s *userService) RedeemPrizeList(ctx *gin.Context) {
	var (
		err  error
		req  userModel.RedeemPrizeListReq
		resp discoverModel.RedeemPrizeLogResp
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	req.UserID = ctx.GetString("user_id")
	resp.List, resp.Count, err = discoverRepo.DiscoverRepo.RedeemPrizeLogForAPI(req)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, code.ErrUnknown)
			return
		}
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(ctx, resp)
}

func (s *userService) PrizeList(c *gin.Context) {
	req := new(userModel.GetPrizeListReq)
	resp := new(discoverModel.PrizeListResp)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	req.Status = 1
	listReq := discoverModel.PrizeListReq{}
	_ = util.CopyStructFields(&listReq, &req)
	resp.List, resp.Count, err = discoverRepo.DiscoverRepo.PrizeList(listReq)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrUnknown)
			return
		}
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}
