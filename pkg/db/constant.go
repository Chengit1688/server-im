package db

const (
	UserIDSize  = 10 // 用户ID随机长度
	GroupIDSize = 10 // 群ID随机长度
)

const (
	ConversationVersion  = "conversation_version"  // 会话版本
	ConversationReadSeq  = "conversation_read_seq" // 会话已读 seq
	MessageSeq           = "message_seq"           // 消息序列号
	MessageStatus        = "message_status"        // 消息状态
	RTC                  = "rtc"                   // rtc
	TokenKey             = "access_token:user_"
	UserInfoKey          = "user_info:user_"
	AccountTempCode      = "account_temp_code:"
	FriendInfo           = "friend_info"            // 好友信息
	GroupMemberInfo      = "group_member_info"      // 群成员
	GetMomentsMessage    = "Moments_info"           // 朋友圈
	UserIPRegisterKey    = "user_info:regip_"       // 用户注册ip
	SystemConfigKey      = "system_config:info"     // 系统配置信息
	SystemConfigLockKey  = "system_config:lock"     // 系统配置信息
	DefaultFriendVersion = "default_friend_version" // 默认好友计数器
	LimitRate            = "limit_rate"             //限制流量
)
