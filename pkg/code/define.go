package code

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 用户从10100 开始  好友从10200开始 群从10300 聊天从10400开始, 通用错误码就从500开始
var (
	ErrUnknown           = newError(500, "服务器异常")
	ErrDB                = newError(501, "数据库错误")
	ErrUserPermissions   = newError(502, "权限不足")
	ErrImSitePermissions = newError(503, "此站点权限不足")
	ErrImSiteParmLimit   = newError(504, "配置限制超过上限")
	ErrTaskBusy          = newError(505, "当前任务正在进行中")
	ErrDataNotExists     = newError(506, "查询的数据不存在")

	ErrUnauthorized    = newError(10000, "请重新登录")
	ErrBadRequest      = newError(10001, "请求参数错误")
	ErrNotLogin        = newError(10002, "用户未登录,请登录")
	ErrNoPermission    = newError(10004, "无权访问")
	ErrFailRequest     = newError(10005, "请求失败，请重试")
	ErrWrongPassword   = newError(10006, "密码错误")
	ErrIpBlocked       = newError(10007, "网络IP异常,请联系接待")
	ErrRegisterBlocked = newError(10008, "注册入口已关闭,请联系接待")
	ErrDeviceIDBlocked = newError(10009, "设备ID异常,请联系接待")
	ErrNickNameUsed    = newError(10010, "当前昵称已经被使用")
	ErrMutePeriod      = newError(10011, "请输入正确的时间段,以-分割")

	//用户
	ErrLoginDenied           = newError(10101, "你无法登录，因为没有登录权限")
	ErrLoadUserInfo          = newError(10102, "获取数据失败，请重试")
	ErrBadPassword           = newError(10103, "密码长度需在6～16之间")
	ErrBadAccount            = newError(10104, "请正确输入账号")
	ErrBadPhoneNumPwd        = newError(10105, "手机号或密码错误")
	ErrBadPhoneNumber        = newError(10106, "请正确输入手机号")
	ErrUserIdNotExist        = newError(10107, "帐号不存在")
	ErrRegisterFailed        = newError(10108, "注册失败，请重试")
	ErrConfigPassword        = newError(10109, "设置密码失败，请重试")
	ErrUserFreeze            = newError(10110, "该账号已被冻结")
	ErrInviteCode            = newError(10111, "请输入邀请码")
	ErrVerificationCode      = newError(10112, "请验证图形验证码")
	ErrSmsCode               = newError(10113, "短信验证码为必输项")
	ErrSettingNotExist       = newError(10114, "配置信息不存在")
	ErrBadCode               = newError(10115, "验证码错误")
	ErrAccountExist          = newError(10116, "账号已注册，不能重复注册")
	ErrPhoneNumberExist      = newError(10116, "手机号已注册，不能重复注册")
	ErrUpdateAccount         = newError(10117, "设置密码失败，请重试")
	ErrVersionNotExist       = newError(10118, "未设置可用版本号")
	ErrVersionExist          = newError(10119, "版本号已经存在，不可再添加")
	ErrInviteCodeExist       = newError(10120, "邀请码已经存在，不可再添加")
	ErrVersionRepeat         = newError(10121, "该平台已经存在一个开启的版本")
	ErrBadAccountPwd         = newError(10122, "账号或密码错误")
	ErrBadInviteCode         = newError(10123, "请输入正确的邀请码")
	ErrUserNotFound          = newError(10124, "用户找不到")
	ErrDefaultFriendNotFound = newError(10125, "默认好友不存在")
	ErrBadCaptcha            = newError(10126, "图形验证码获取失败，请重试")
	ErrBadSmsFactory         = newError(10127, "获取短信厂商失败，请重试")
	ErrBadSmsSend            = newError(10128, "一分钟内验证码不能重复发送")
	ErrSmsSend               = newError(10129, "发送验证码失败")
	ErrSignLog               = newError(10130, "今日已签到")
	ErrUserConfig            = newError(10131, "用户配置不存在")
	ErrBadSignLog            = newError(10132, "签到记录不存在")
	ErrBadRequestType        = newError(10133, "请求参数的类型错误")
	ErrUserIdExist           = newError(10134, "帐号已经存在")
	ErrIPBlackListExist      = newError(10135, "IP黑名单已经存在")
	ErrDefaultFriendFound    = newError(10136, "默认好友已存在")
	ErrRoleNameExist         = newError(10137, "此角色描述已存在")
	ErrRoleKeyExist          = newError(10138, "此角色名称已存在")
	ErrSameAccount           = newError(10139, "重复账号无法创建")
	ErrShieldExist           = newError(10140, "敏感词已经存在，不可再添加")
	ErrPrivilegeUserExist    = newError(10141, "该用户已经是特权")
	ErrRegisterLimit         = newError(10142, "该设备已达到注册上限")
	ErrRegisterTimeLimit     = newError(10143, "该IP已达到注册上限")
	ErrSuggestion            = newError(10144, "投诉建议不存在")
	ErrBalanceNotEnough      = newError(10145, "账户余额不足")
	ErrRedPackSingleDisable  = newError(10146, "禁止发送个人红包")
	ErrPayPasswdWrong        = newError(10147, "支付密码错误")
	ErrRedPackTimeout        = newError(10148, "该红包已超过24小时。如已领取，可在钱包明细中查看")
	ErrPayPasswdAlreadySet   = newError(10149, "支付密码已设置，忘记请找回支付密码")
	ErrInvalidAccount        = newError(10150, "账号中不能包含特殊字符")
	ErrSameNickname          = newError(10151, "重复昵称无法创建")
	ErrIPWhiteListExist      = newError(10152, "IP白名单已经存在")
	ErrIPNotInWhiteList      = newError(10153, "IP不在白名单里")
	ErrUserFavImageExist     = newError(10154, "用户图片收藏已存在")
	ErrRedPackPreAmount      = newError(10155, "普通红包需设置单个红包金额")
	ErrRedPackReceive        = newError(10156, "红包领取失败")
	ErrRedPackReceiveFinish  = newError(10157, "红包已被领完")
	ErrRedPackReceiveRepeat  = newError(10158, "红包已领取，不能再次领")
	ErrCustomerUserExist     = newError(10159, "该用户已经是客服人员")
	ErrRegisterLen           = newError(10160, "账号不安全,不得少于四个字符")
	//好友
	ErrAlreadyIsFriend      = newError(10200, "对方已经是好友了")
	ErrFriendNotExist       = newError(10201, "好友不存在")
	ErrFriendCanNotSelf     = newError(10202, "不能添加自己为好友")
	ErrFriendSearchNotExist = newError(10203, "用户不存在，请重试")
	ErrFriendNormalNotAdd   = newError(10204, "不允许加好友")
	ErrFriendIsMax          = newError(10205, "好友数量达到上限")
	ErrFriendNormalNotDel   = newError(10206, "不允许删除好友")
	ErrFriendLabelExist     = newError(10207, "好友分组已经存在")
	ErrFriendLabelForbiden  = newError(10208, "默认分组不可修改")
	ErrFriendLabelDelete    = newError(10209, "默认分组不可删除")
	ErrFriendLabelLimit     = newError(10210, "分组数量超过最大限制")
	ErrFriendLabelNotExist  = newError(10211, "当前分组不存在")

	//群组
	ErrNotInGroup           = newError(10300, "用户没有在群组中")
	ErrGroupNotExist        = newError(10301, "群组不存在")
	ErrGroupNotMember       = newError(10302, "群内没有成员")
	ErrApplyNotFound        = newError(10303, "申请不存在")
	ErrApplyDone            = newError(10304, "申请已被处理，不能重复审核")
	ErrAlreadyInGroup       = newError(10305, "用户已经在群组中")
	ErrOwnerCanNotQuit      = newError(10306, "群主不能退群")
	ErrSetOwnerOnlyAdmin    = newError(10307, "只能设置管理员为群主")
	ErrFuncRunning          = newError(10308, "其他管理员正在执行该操作，请稍后再试！")
	ErrCloseQuit            = newError(10309, "当前群禁止退群")
	ErrGroupMemberMax       = newError(10340, "群成员已满")
	ErrGroupMemberOutMax    = newError(10341, "群成员将超过最大数量")
	ErrGroupNormalNotJoin   = newError(10342, "普通用户不允许加群")
	ErrGroupIsMax           = newError(10343, "创建的群达到上限")
	ErrGroupNormalNotCreate = newError(10344, "普通用户不允许创建群聊")
	ErrGroupMuteAll         = newError(10345, "当前群体禁言")
	ErrGroupMuteAllPeriod   = newError(10345, "当前在禁言时间段")
	ErrGroupMuteUser        = newError(10345, "您已被禁言")

	//聊天
	ErrMsgSendRunning    = newError(10400, "其他管理员正在执行该操作，请稍后再试！")
	ErrChatRTCNotFound   = newError(10401, "通话不存在")
	ErrChatRTCStatus     = newError(10402, "通话状态异常")
	ErrChatRTCBusy       = newError(10403, "通话忙")
	ErrChatRTCTargetBusy = newError(10404, "对方通话忙")
	ErrChatRTCUser       = newError(10405, "通话用户不合法")
	ErrChatRTCDevice     = newError(10406, "通话设备不合法")
	ErrChatRTCNetworkBad = newError(10407, "通话网络不佳")
	ErrChatRTCAbort      = newError(10408, "通话中断")
	ErrChatRevoke        = newError(10409, "该消息无法撤销")
	ErrChatTooMuch       = newError(10410, "发送消息太频繁")

	//联系人
	ErrTagExists          = newError(20410, "tag已经存在")
	ErrTagNotExists       = newError(20411, "tag不存在")    //tag不存在
	ErrTagFriendNotExists = newError(20412, "好友不在tag中")  //好友不在tag中
	ErrTagFriendExists    = newError(20413, "好友已经在tag中") //好友已经在tag中

	//店铺
	ErrShopExists           = newError(30001, "shop已经存在")
	ErrShopNotExists        = newError(30001, "shop不存在")
	ErrInviteCodeExists     = newError(30002, "inviteCode已经存在")
	ErrInviteCodeNotExists  = newError(30003, "inviteCode不存在")
	ErrShopTeamMemberExists = newError(30004, "团队成员已存在，无需重复加入")
	ErrShopExSelf           = newError(30005, "不能加入自己创建的团队")
)

//friend error message

func newError(code int, message string) error {
	return status.New(codes.Code(code), message).Err()
}

func NewErrorAs(code int, message string) *errorAsMsg {
	var e = new(errorAsMsg)
	e.Err = newError(code, message)
	return e
}

type errorAsMsg struct {
	code int
	Err  error
}

func (e *errorAsMsg) SetUserMsg(format string, a ...interface{}) error {
	message := fmt.Sprintf(format, a...)
	return newError(e.code, message)
}
