package response

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

const (
	ErrUnknown           = 500
	ErrDB                = 501
	ErrUserPermissions   = 502
	ErrImSitePermissions = 503
	ErrImSiteParmLimit   = 504
	ErrTaskBusy          = 505
	ErrDataNotExists     = 506

	ErrUnauthorized    = 10000
	ErrBadRequest      = 10001
	ErrNotLogin        = 10002
	ErrNoPermission    = 10004
	ErrFailRequest     = 10005
	ErrWrongPassword   = 10006
	ErrIpBlocked       = 10007
	ErrRegisterBlocked = 10008
	ErrDeviceIDBlocked = 10009
	ErrNickNameUsed    = 10010
	ErrMutePeriod      = 10011

	//用户
	ErrLoginDenied           = 10101
	ErrLoadUserInfo          = 10102
	ErrBadPassword           = 10103
	ErrBadAccount            = 10104
	ErrBadPhoneNumPwd        = 10105
	ErrBadPhoneNumber        = 10106
	ErrUserIdNotExist        = 10107
	ErrRegisterFailed        = 10108
	ErrConfigPassword        = 10109
	ErrUserFreeze            = 10110
	ErrInviteCode            = 10111
	ErrVerificationCode      = 10112
	ErrSmsCode               = 10113
	ErrSettingNotExist       = 10114
	ErrBadCode               = 10115
	ErrPhoneNumberExist      = 10116
	ErrUpdateAccount         = 10117
	ErrVersionNotExist       = 10118
	ErrVersionExist          = 10119
	ErrInviteCodeExist       = 10120
	ErrVersionRepeat         = 10121
	ErrBadAccountPwd         = 10122
	ErrBadInviteCode         = 10123
	ErrUserNotFound          = 10124
	ErrDefaultFriendNotFound = 10125
	ErrBadCaptcha            = 10126
	ErrBadSmsFactory         = 10127
	ErrBadSmsSend            = 10128
	ErrSmsSend               = 10129
	ErrSignLog               = 10130
	ErrUserConfig            = 10131
	ErrBadSignLog            = 10132
	ErrBadRequestType        = 10133
	ErrUserIdExist           = 10134
	ErrIPBlackListExist      = 10135
	ErrDefaultFriendFound    = 10136
	ErrRoleNameExist         = 10137
	ErrRoleKeyExist          = 10138
	ErrSameAccount           = 10139
	ErrShieldExist           = 10140
	ErrPrivilegeUserExist    = 10141
	ErrRegisterLimit         = 10142
	ErrRegisterTimeLimit     = 10143
	ErrSuggestion            = 10144
	ErrBalanceNotEnough      = 10145
	ErrRedPackSingleDisable  = 10146
	ErrPayPasswdWrong        = 10147
	ErrRedPackTimeout        = 10148
	ErrPayPasswdAlreadySet   = 10149
	ErrInvalidAccount        = 10150
	ErrSameNickname          = 10151
	ErrIPWhiteListExist      = 10152
	ErrIPNotInWhiteList      = 10153
	ErrUserFavImageExist     = 10154
	ErrRedPackPreAmount      = 10155
	ErrRedPackReceive        = 10156
	ErrRedPackReceiveFinish  = 10157
	ErrRedPackReceiveRepeat  = 10158
	ErrCustomerUserExist     = 10159
	ErrRegisterLen           = 10160
	ErrAccountExist          = 10161
	//好友
	ErrAlreadyIsFriend      = 10200
	ErrFriendNotExist       = 10201
	ErrFriendCanNotSelf     = 10202
	ErrFriendSearchNotExist = 10203
	ErrFriendNormalNotAdd   = 10204
	ErrFriendIsMax          = 10205
	ErrFriendNormalNotDel   = 10206
	ErrFriendLabelExist     = 10207
	ErrFriendLabelForbiden  = 10208
	ErrFriendLabelDelete    = 10209
	ErrFriendLabelLimit     = 10210
	ErrFriendLabelNotExist  = 10211
	ErrFriendInBlack        = 10212
	ErrFriendExistInBlack   = 10213

	//群组
	ErrNotInGroup           = 10300
	ErrGroupNotExist        = 10301
	ErrGroupNotMember       = 10302
	ErrApplyNotFound        = 10303
	ErrApplyDone            = 10304
	ErrAlreadyInGroup       = 10305
	ErrOwnerCanNotQuit      = 10306
	ErrSetOwnerOnlyAdmin    = 10307
	ErrFuncRunning          = 10308
	ErrCloseQuit            = 10309
	ErrGroupMemberMax       = 10340
	ErrGroupMemberOutMax    = 10341
	ErrGroupNormalNotJoin   = 10342
	ErrGroupIsMax           = 10343
	ErrGroupNormalNotCreate = 10344
	ErrGroupMuteAll         = 10345
	ErrGroupMuteAllPeriod   = 10346
	ErrGroupMuteUser        = 10347

	//聊天
	ErrMsgSendRunning    = 10400
	ErrChatRTCNotFound   = 10401
	ErrChatRTCStatus     = 10402
	ErrChatRTCBusy       = 10403
	ErrChatRTCTargetBusy = 10404
	ErrChatRTCUser       = 10405
	ErrChatRTCDevice     = 10406
	ErrChatRTCNetworkBad = 10407
	ErrChatRTCAbort      = 10408
	ErrChatRevoke        = 10409
	ErrChatTooMuch       = 10410

	//联系人
	ErrTagExists          = 20410 //tag已经存在
	ErrTagNotExists       = 20411 //tag不存在
	ErrTagFriendNotExists = 20412 //好友不在tag中
	ErrTagFriendExists    = 20413 //好友已在tag中

	//加盟商
	ErrShopExists           = 30001 //shop已经存在
	ErrInviteCodeExists     = 30002 //inviteCode已经存在
	ErrShopNotExists        = 30003 //shop不存在
	ErrInviteCodeNotExists  = 30004 //邀请码不存在
	ErrShopTeamMemberExists = 30005 //团队成员已存在
	ErrShopExSelf           = 30006 //不能加入自己创建的团队
)

func newError(code int, message string) error {
	return status.New(codes.Code(code), message).Err()
}

func SelectLang(lang string) *sync.Map {
	switch lang {
	case "en_US":
		return enResponses
	case "zh_CN":
		return cnResponses
	case "ja":
		return japResponses
	default:
		return cnResponses
	}
}

func GetError(code int, lang string) error {
	if data, ok := SelectLang(lang).Load(code); ok {
		message := data.(string)
		return newError(code, message)
	}
	return nil
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
