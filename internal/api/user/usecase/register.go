package usecase

import (
	"encoding/json"
	"fmt"
	"im/config"
	groupModel "im/internal/api/group/model"
	settingModel "im/internal/api/setting/model"
	settingRepo "im/internal/api/setting/repo"
	userModel "im/internal/api/user/model"
	userRepo "im/internal/api/user/repo"
	"im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterUseCase struct {
	ctx        *gin.Context
	FriendCase IAuthFriend
	GroupCase  IAuthGroup
}
type IAuthFriend interface {
	AddFriend(operationID string, userID, friendID string, remark string, greeting string, isCustomer bool) error
	CreateFriendLabel(userID, labelID, labelName string) error
}

type IAuthGroup interface {
	JoinGroup(operationID string, groupId, userId string) error
	GetDefaultGroups() (res []groupModel.Group)
}

func NewRegisterUseCase(ctx *gin.Context, f IAuthFriend, g IAuthGroup) *RegisterUseCase {
	return &RegisterUseCase{
		ctx:        ctx,
		FriendCase: f,
		GroupCase:  g,
	}
}

func (s *RegisterUseCase) Register(req userModel.RegisterReq) (*userModel.User, error) {
	var (
		err        error
		user       *userModel.User
		userResult *userModel.User
		sysConfig  *model.ParameterConfigResp
	)
	if sysConfig, err = userRepo.UserCache.GetSystemConfigInfo(configRepo.ConfigRepo.GetSystemConfig); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("register system config, error: %v", err))
		return nil, code.ErrFailRequest
	}
	if sysConfig.IsOpenRegister == constant.SwitchOff {
		return nil, code.ErrRegisterBlocked
	}

	if user, err = s.registerRule(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("register, error: %v", err))
		return nil, err
	}

	opt := userRepo.WhereOption{
		UserId:      user.UserID,
		Account:     user.Account,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
	}
	if b := userRepo.UserRepo.OrExists(opt); b {
		if req.AccountType == constant.PhoneNumberStr {
			return nil, code.ErrPhoneNumberExist
		}
		return nil, code.ErrAccountExist
	}

	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "info", "emqx create auth username")
	if err = AuthUseCase.CreateAuthUsername(user.UserID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("register error, create auth username error, error: %v", err))
		return nil, code.ErrRegisterFailed
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "info", "emqx create auth username success")

	user.Salt = util.RandString(6)
	user.Password = util.GetPassword(req.Password, user.Salt)
	user.InviteCode = req.InviteCode
	user.DeviceId = req.DeviceId
	user.RegisterDeviceId = req.DeviceId
	user.RegisterIp = s.ctx.ClientIP()
	user.LoginIp = user.RegisterIp
	user.ImSite = req.ImSite
	user.LatestLoginTime = time.Now().Unix()

	if userResult, err = userRepo.UserRepo.Create(user); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error creating registered user, error: %v", err))
		return nil, code.ErrRegisterFailed
	}

	s.FriendCase.CreateFriendLabel(userResult.UserID, userResult.UserID, "我的好友")

	s.defaultInviteHandle(req.OperationID, req.InviteCode, userResult.UserID)

	err = userRepo.UserRepo.RecordUserDeviceAndIP(userResult.UserID, userResult.DeviceId, userResult.RegisterIp, req.Platform)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("RecordUserDeviceAndIP, error: %v", err))
	}
	userRepo.UserCache.RecordRegisterDeviceIDCount(userResult.DeviceId)
	userRepo.UserCache.RecordRegisterIPCount(userResult.RegisterIp)
	return userResult, nil
}

func (s *RegisterUseCase) registerRule(req *userModel.RegisterReq) (*userModel.User, error) {
	var (
		err error
		c   *settingModel.SettingConfig
		reg settingModel.RegisterConfigInfo
	)

	if err = s.SystemConfigRule(req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("systemConfigRule, error: %v", err))
		return nil, err
	}

	opt := settingRepo.WhereOption{
		ConfigType: []string{configRepo.ConfigRegisterConfig},
	}

	if c, err = settingRepo.SettingRepo.GetInfo(opt); c != nil {
		if c.Content != "" {
			if err = json.Unmarshal([]byte(c.Content), &reg); err != nil {
				return nil, err
			}
		}
	}

	if req.AccountType == constant.AccountStr {
		if len(req.Account) < constant.AccountLen {
			return nil, code.ErrBadAccount
		}

		if reg.IsAllAccount == constant.SwitchOn {
			if !util.IsAlphaNumeric(req.Account) {
				return nil, code.ErrInvalidAccount
			}
		}

		if reg.IsAllAccount == constant.SwitchOff {
			if !util.IsAlphaNumericChinese(req.Account) {
				return nil, code.ErrInvalidAccount
			}
		}
	}

	if req.AccountType == constant.PhoneNumberStr {
		if req.PhoneNumber == "" {
			return nil, code.ErrBadPhoneNumber
		}

		if req.CountryCode == "" && !strings.HasPrefix(req.PhoneNumber, constant.PhoneNumberPrefix) {
			return nil, code.ErrBadPhoneNumber
		}
		if strings.HasPrefix(req.PhoneNumber, constant.PhoneNumberPrefix) {
			req.CountryCode = ""
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
		if req.SmsCode == "" {
			return nil, code.ErrSmsCode
		}
		if err = s.smsVerify(req.PhoneNumber, req.SmsCode); err != nil {
			return nil, err
		}
	}
	if reg.IsInviteCode == constant.SwitchOn {
		if req.InviteCode == "" && reg.CheckInviteCode == constant.SwitchOn {
			return nil, code.ErrInviteCode
		}
		u.InviteCode = req.InviteCode
		if !s.inviteCodeVerify(req.InviteCode) {
			u.InviteCode = ""
			if reg.CheckInviteCode == constant.SwitchOn {
				return nil, code.ErrBadInviteCode
			}
		}
	}
	if reg.IsVerificationCode == constant.SwitchOn {
		if req.VerificationToken == "" || req.VerificationPoint == "" || req.CaptchaType == "" {
			return nil, code.ErrVerificationCode
		}
		if err = GetCaptchaFactory().GetService(req.CaptchaType).Verification(req.VerificationToken, req.VerificationPoint); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (s *RegisterUseCase) SystemConfigRule(req *userModel.RegisterReq) error {
	var (
		err         error
		paramConfig *model.ParameterConfigResp
	)
	if paramConfig, err = userRepo.UserCache.GetSystemConfigInfo(configRepo.ConfigRepo.GetSystemConfig); err != nil {
		return code.ErrFailRequest
	}
	regIp := s.ctx.ClientIP()
	if paramConfig.IpRegLimitTime == 0 {
		userRepo.UserCache.ClearAllIPLimit()
	} else {
		isLimit := userRepo.UserCache.CheckIpLimit(regIp, paramConfig.IpRegLimitCount, paramConfig.IpRegLimitTime)
		if isLimit {
			return code.ErrRegisterTimeLimit
		}
	}

	if paramConfig.DeviceRegLimit > 0 {
		isLimit := userRepo.UserCache.CheckDeviceLimit(req.DeviceId, paramConfig.DeviceRegLimit)
		if isLimit {
			return code.ErrRegisterLimit
		}
	}

	return nil
}

func (s *RegisterUseCase) smsVerify(PhoneNumber string, smsCode string) error {
	if smsCode == config.Config.Sms.SuperCode {
		return nil
	}
	var (
		cacheCode string
		err       error
	)
	if cacheCode, err = settingRepo.SettingCache.GetAccountCode(PhoneNumber); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("smsVerify GetAccountCode, error: %v", err))
		return code.ErrBadCode
	}
	if smsCode != cacheCode {
		return code.ErrBadCode
	}
	_, _ = settingRepo.SettingCache.DelAccountCode(PhoneNumber)

	return nil
}

func (s *RegisterUseCase) inviteCodeVerify(inviteCode string) bool {
	var (
		count int64
		err   error
	)
	if inviteCode == "" {
		return true
	}
	opt := configRepo.WhereOptionForInvite{
		InviteCode:   inviteCode,
		Status:       constant.SwitchOn,
		DeleteStatus: constant.SwitchOff,
	}
	if count, err = configRepo.InviteCode.Exists(opt); count > 0 {
		return true
	}
	if err != nil {
		logger.Sugar.Warnw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user invite code error, error: %v", err))
	}

	return false
}

func (s *RegisterUseCase) defaultInviteHandle(operationID string, inviteCode string, userId string) bool {
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

func (s *RegisterUseCase) joinDefaultGroups(operationID string, inviteCode *model.InviteCode, userId string) {
	var err error
	go func() {

		if inviteCode != nil {
			if inviteCode.DefaultGroups == "" {
				return
			}
			groupList := strings.Split(inviteCode.DefaultGroups, ",")

			if inviteCode.IsOpenTurn == constant.SwitchOff {
				for _, groupId := range groupList {
					if err = s.GroupCase.JoinGroup(operationID, groupId, userId); err != nil {
						logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user joinDefaultGroups, error: %v", err))
						continue
					}
				}
			}

			if inviteCode.IsOpenTurn == constant.SwitchOn {
				if inviteCode.GroupIndex >= int64(len(groupList)) {
					inviteCode.GroupIndex = 0
				}
				groupId := groupList[inviteCode.GroupIndex]
				if err = s.GroupCase.JoinGroup(operationID, groupId, userId); err != nil {
					logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user joinDefaultGroups on turn, error: %v", err))
				}

				inviteCode.GroupIndex = inviteCode.GroupIndex + 1

				opt := configRepo.WhereOptionForInvite{
					Id: inviteCode.ID,
				}
				if err = configRepo.InviteCode.UpdateById(opt, inviteCode); err != nil {
					logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user joinDefaultGroups on turn, error: %v", err))
					return
				}
			}

			return
		}

		groupList := s.GroupCase.GetDefaultGroups()

		if isOpen, err := configRepo.ConfigRepo.GetDefaultIsOpen(); err == nil && isOpen == constant.SwitchOn {
			index, _ := configRepo.ConfigRepo.GetDefaultGroupIndex()
			if index >= len(groupList) {
				index = 0
			}
			group := groupList[index]
			if err = s.GroupCase.JoinGroup(operationID, group.GroupId, userId); err != nil {
				logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user joinDefaultGroups on turn, error: %v", err))
				return
			}

			index = index + 1
			configRepo.ConfigRepo.SetDefaultGroupIndex(index)
			return
		}

		for _, g := range groupList {
			if err = s.GroupCase.JoinGroup(operationID, g.GroupId, userId); err != nil {
				logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user joinDefaultGroups, error: %v", err))
				continue
			}
		}
	}()
}

func (s *RegisterUseCase) addDefaultFriends(operationID string, inviteCode *model.InviteCode, userId string) {
	var (
		err        error
		friendList []model.DefaultFriend
	)
	go func() {

		if inviteCode != nil {
			if inviteCode.DefaultFriends == "" {
				return
			}
			friendList := strings.Split(inviteCode.DefaultFriends, ",")

			if inviteCode.IsOpenTurn == constant.SwitchOff {
				for _, friendId := range friendList {
					if err = s.FriendCase.AddFriend(operationID, userId, friendId, "", inviteCode.GreetMsg, false); err != nil {
						logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user addDefaultFriends, error: %v", err))
						continue
					}
				}
			}

			if inviteCode.IsOpenTurn == constant.SwitchOn {
				if inviteCode.FriendIndex >= int64(len(friendList)) {
					inviteCode.FriendIndex = 0
				}
				friendId := friendList[inviteCode.FriendIndex]
				if err = s.FriendCase.AddFriend(operationID, userId, friendId, "", inviteCode.GreetMsg, false); err != nil {
					logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user addDefaultFriends on turn, error: %v", err))
				}

				inviteCode.FriendIndex = inviteCode.FriendIndex + 1

				opt := configRepo.WhereOptionForInvite{
					Id: inviteCode.ID,
				}
				if err = configRepo.InviteCode.UpdateById(opt, inviteCode); err != nil {
					logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user addDefaultFriends on turn, error: %v", err))
					return
				}
			}

			return
		}

		if friendList, err = configRepo.DefaultFriendRepo.GetAllDefalutFriendIdsAndGreetMsg(); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user addDefaultFriends, error: %v", err))
			return
		}

		if isOpen, err := configRepo.ConfigRepo.GetDefaultIsOpen(); err == nil && isOpen == constant.SwitchOn {
			index, _ := configRepo.ConfigRepo.GetDefaultFriendIndex()
			if index >= len(friendList) {
				index = 0
			}
			friendInfo := friendList[index]
			if err = s.FriendCase.AddFriend(operationID, userId, friendInfo.UserId, "", friendInfo.GreetMsg, false); err != nil {
				logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user joinDefaultGroups on turn, error: %v", err))
				return
			}

			index = index + 1
			configRepo.ConfigRepo.SetDefaultFriendIndex(index)
			return
		}

		for _, friendInfo := range friendList {
			if err = s.FriendCase.AddFriend(operationID, userId, friendInfo.UserId, "", friendInfo.GreetMsg, false); err != nil {
				logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user addDefaultFriends, error: %v", err))
				continue
			}
		}
	}()
}
