package common

var ContentKey = []byte("1122334455667788") // 消息内容加密key

type MessageType int32

/*
code区间定义
200-299 好友相关
300-399 群相关
400-499 聊天相关
*/
const (
	UserInfoPush          MessageType = 100 // 用户信息推送
	FreezeUserPush        MessageType = 150 // 冻结用户账号推送
	ConfigUserChangePush  MessageType = 101 //用户配置更新
	UserFavoriteImagePush MessageType = 102 //用户图片收藏更新
	UserStatusOffline     MessageType = 103 //下线通知
	UserStatusOnline      MessageType = 104 //上线通知

	FriendRequestPush    MessageType = 200 //好友请求
	AddFriendPush        MessageType = 201 //添加好友
	AddFriendAckPush     MessageType = 202 //添加好友回应
	RemoveFriendPush     MessageType = 203 //删除好友
	ChangeFriendPush     MessageType = 204 //变化好友
	FriendInfoLableChage MessageType = 205 //好友信息分组变化
	FriendLabelchage     MessageType = 206 //好友分组变化
	FriendInBlack        MessageType = 207 //好友在黑名单

	GroupAddPush          MessageType = 300 //新增的群
	GroupRemovePush       MessageType = 301 //移除群
	GroupChangePush       MessageType = 302 //群更新通知
	GroupMemberApplyPush  MessageType = 303 //入群申请
	GroupMemberVerifyPush MessageType = 304 //入群申请审核
	GroupMemberAddPush    MessageType = 305 //群用户新增
	GroupMemberRemovePush MessageType = 306 //群用户删除
	GroupMemberChangePush MessageType = 307 //群用户更新
	GroupMemberAtPush     MessageType = 308 //At用户

	ChatMessagePush           MessageType = 400 // 聊天消息推送
	ChatMessageAckSeqPush     MessageType = 401 // 聊天消息已读回执推送
	ChatMessageClearPush      MessageType = 402 // 聊天消息清空推送
	ChatMessageAdminClearPush MessageType = 403 // 聊天消息客户端删除推送
	ChatMessageRTCPush        MessageType = 404 // 聊天消息RTC推送
	ChatMessageReadSeqPush    MessageType = 405 // 聊天消息已读推送

	SysConfigBroadcast MessageType = 500 //系统信息推送
	ConfigShieldPush   MessageType = 501 // 铭感词推送

	RedpackSingleRecvPush MessageType = 600 //个人红包领取通知推送
	RedpackGroupRecvPush  MessageType = 601 //群红包领取通知推送

	MomentsInviteFriend  MessageType = 701 //朋友圈@提醒谁看
	MomentsLikeFriend    MessageType = 702 //朋友圈点赞
	MomentsCommentFriend MessageType = 703 //朋友圈评论

	LinkUpdatePush MessageType = 801 //链接更新成功推送内容

)

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}
