package constant

const (
	SwitchOn               = 1 //开关开启
	SwitchOff              = 2 //开关关闭
	DataNormal             = 1 //数据正常
	DataAbNormal           = 2 //数据不正常
	PhoneNumberStr         = 1 //电话号码帐号
	AccountStr             = 2 //非电话号码帐号
	GuestStr               = 3 //游客帐号
	MixedChar              = 3 //数字+字母
	UserStatusNormal       = 1 //正常
	UserStatusFreeze       = 2 //冻结
	ChinaCountryCode       = "+86"
	PhoneNumberPrefix      = "+"
	ImSiteZhaoCai          = "zhaocai" //其他站点名
	ImSiteIm               = "im"      //im站点名
	ImSiteHeaderStr        = "imsite"
	GuestUserModel         = 2    // 登录模式，1：正常注册的帐号，2：游客模式
	NormalUserModel        = 1    // 登录模式，1：正常注册的帐号，2：游客模式
	UserIdPrefix           = "用户" //user_id前缀
	PhoneNumberForRegister = 1    //电话验证码用来注册
	PhoneNumberForChangeWd = 2    //电话验证码用来找回密码
	PhoneNumberForBand     = 3    //电话验证码用来绑定手机号
	MuteMemberOpen         = 1    //禁言所有人
	MuteMemberClose        = 2    //不禁言
	MuteNewJoinMember      = 3    //新入群禁言
	MuteMemberPeriod       = 4    //时间段禁言
	AccountLen             = 4
)
