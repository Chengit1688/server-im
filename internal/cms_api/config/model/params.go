package model

import (
	"im/pkg/pagination"
	"time"
)

type MenuListRespItem struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Title      string     `json:"title"`
	Icon       string     `json:"icon"`
	Path       string     `json:"path"`
	Paths      string     `json:"paths"`
	Type       int        `json:"type"`
	Action     string     `json:"action"`
	Permission string     `json:"permission"`
	ParentId   int        `json:"parent_id"`
	NoCache    int        `json:"no_cache"`
	Component  string     `json:"component"`
	Sort       int        `json:"sort"`
	Visible    int        `json:"visible"`
	Hidden     int        `json:"hidden"`
	IsFrame    int        `json:"is_frame"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type GetMenuConfigReq struct {
	Timestamp int64 `json:"timestamp" form:"timestamp" binding:"required" msg:"时间戳不能为空"`
}

type GetMenuConfigRespData struct {
	Menus     []MenuListRespItem `json:"menus"`
	Timestamp int64              `json:"timestamp" form:"timestamp"`
}

type GetMenuConfigResp struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    GetMenuConfigRespData `json:"data"`
}

type VersionReq struct {
	Id          int64  `json:"id"  binding:"omitempty,number"`
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	Platform    int64  `json:"platform" binding:"required,oneof=1 2 3 4 5 6 7 8 9"`
	Version     string `json:"version" binding:"required,gte=1" msg:"请准确填写版本号"`
	IsForce     int64  `json:"is_force"  binding:"required,oneof=1 2"`
	Title       string `json:"title" binding:"required,min=1,max=50" msg:"请准确填写，最多50字"`
	DownloadUrl string `json:"download_url"  binding:"required,url,min=1" msg:"请准确填写下载地址"`
	UpdateDesc  string `json:"update_desc" binding:"required,max=200" msg:"请准确填写，更新说明最多200字"`
}

type VersionUpdateStatusReq struct {
	Id          int64  `json:"id"  binding:"omitempty,number"`
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	Platform    int64  `json:"platform" binding:"oneof=1 2 3 4 5 6 7 8 9"`
	Status      int64  `json:"status" binding:"required,number"`
}

type VersionDeleteReq struct {
	Id          int64  `json:"id"  binding:"omitempty,number"`
	OperationID string `json:"operation_id"  binding:"required,min=1"`
}

type VersionListReq struct {
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	Status      int64  `json:"status" binding:"omitempty,number"`
	Platform    int64  `json:"platform" binding:"omitempty,oneof=1 2 3 4 5 6 7 8 9"`
	BeginDate   int64  `json:"begin_date" binding:"omitempty,number"`
	EndDate     int64  `json:"end_date"  binding:"omitempty,number"`
	pagination.Pagination
}
type VersionListInfo struct {
	ID          int64  `json:"id"`
	Platform    int64  `json:"platform"`
	Version     string `json:"version"`
	IsForce     int64  `json:"is_force"`
	Title       string `json:"title"`
	DownloadUrl string `json:"download_url"`
	UpdateDesc  string `json:"update_desc"`
	Status      int64  `json:"status"`
	CreatedAt   int64  `json:"created_at"`
}

type VersionListResp struct {
	List     []VersionListInfo `json:"list"`
	Count    int64             `json:"count"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

type GetLoginConfigReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type GetLoginConfigResp struct {
	Pc     []int `json:"pc" binding:"required,min=1,max=2"`
	Mobile []int `json:"mobile" binding:"required,min=1,max=3"`
}

type GetRegisterConfigResp struct {
	CheckInviteCode    int `json:"check_invite_code"`
	IsInviteCode       int `json:"is_invite_code"`
	IsVerificationCode int `json:"is_verification_code"`
	IsSmsCode          int `json:"is_sms_code"`
	IsAllAccount       int `json:"is_all_account"`
}

type UpdateLoginConfigReq struct {
	OperationID string             `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Config      GetLoginConfigResp `json:"config" binding:"required" msg:"配置不能为空"`
}

type UpdateRegisterConfigReq struct {
	OperationID string                `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Config      GetRegisterConfigResp `json:"config" binding:"required" msg:"配置不能为空"`
}

type InviteCodeReq struct {
	Id             int64  `json:"id"  binding:"omitempty,number"`
	OperationID    string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	InviteCode     string `json:"invite_code" binding:"required,min=1,number" msg:"邀请码填写错误"`
	DefaultFriends string `json:"default_friends"  binding:"omitempty,gte=1" msg:"请选择默认好友"`
	DefaultGroups  string `json:"default_groups" binding:"omitempty,gte=1" msg:"请选择默认群组"`
	GreetMsg       string `json:"greet_msg" binding:"omitempty,min=1,max=200" msg:"请填写好友打招呼消息，最多200字"`
	Remarks        string `json:"remarks"  binding:"omitempty,min=1,max=300" msg:"请填写备注"`
}

type InviteUpdateReq struct {
	Id          int64   `json:"id"  binding:"omitempty,number"`
	OperationID string  `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	InviteCode  string  `json:"invite_code" binding:"required,min=1,number" msg:"邀请码填写错误"`
	GreetMsg    *string `json:"greet_msg" binding:"omitempty,max=200" msg:"请填写好友打招呼消息，最多200字"`
	Remarks     *string `json:"remarks"  binding:"omitempty,max=300" msg:"请填写备注,最多300字"`
}

type InviteUpdateFriendReq struct {
	Id             int64  `json:"id"  binding:"omitempty,number"`
	OperationID    string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	InviteCode     string `json:"invite_code" binding:"required,min=1,number" msg:"邀请码填写错误"`
	DefaultFriends string `json:"default_friends"  binding:"omitempty" msg:"请选择默认好友"`
}

type InviteUpdateGroupReq struct {
	Id            int64  `json:"id"  binding:"omitempty,number"`
	OperationID   string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	InviteCode    string `json:"invite_code" binding:"required,min=1,number" msg:"邀请码填写错误"`
	DefaultGroups string `json:"default_groups" binding:"omitempty" msg:"请选择默认群组"`
}

type InviteDeleteReq struct {
	Ids         []int64 `json:"ids"  binding:"required" msg:"ids必须传"`
	OperationID string  `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
}

type InviteUpdateStatusReq struct {
	Id          int64  `json:"id"  binding:"required,number" msg:"id必须传"`
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	Status      int64  `json:"status" binding:"required,number,oneof=1 2" msg:"请选择状态"`
	IsOpenTurn  int64  `json:"is_open_turn" binding:"required,number,oneof=1 2" msg:"请选择是否开启轮流开关"`
}

type InviteListReq struct {
	OperationID    string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	InviteCode     string `json:"invite_code" binding:"omitempty,min=1,number" msg:"邀请码填写错误"`
	DefaultGroups  string `json:"default_group" binding:"omitempty,min=6" msg:"请输入默认群id"`
	DefaultFriends string `json:"default_friend" binding:"omitempty,min=6" msg:"请输入默认好友id"`
	Remarks        string `json:"remarks"  binding:"omitempty,min=1" msg:"请输入备注"`
	OperationUser  string `json:"operation_user"  binding:"omitempty,min=1" msg:"请输入操作者"`
	BeginDate      int64  `json:"begin_date" binding:"omitempty,number" msg:"请输入开始时间"`
	EndDate        int64  `json:"end_date"  binding:"omitempty,number" msg:"请输入结束时间"`
	Status         int64  `json:"status" binding:"omitempty,number,oneof=1 2" msg:"请选择状态"`
	pagination.Pagination
}

type InviteListResp struct {
	List     []InviteListInfo `json:"list"`
	Count    int64            `json:"count"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

type InviteListInfo struct {
	ID             int64         `json:"id"`
	InviteCode     string        `json:"invite_code"`
	DefaultGroups  []interface{} `json:"default_groups"`
	DefaultFriends []interface{} `json:"default_friends"`
	GreetMsg       string        `json:"greet_msg"`
	Remarks        string        `json:"remarks"`
	Status         int64         `json:"status"`
	OperationUser  string        `json:"operation_user"`
	CreatedAt      int64         `json:"created_at"`
	UpdatedAt      int64         `json:"updated_at"`
	IsOpenTurn     int64         `json:"is_open_turn"`
}

type DefaultFriendReq struct {
	Id          int64  `json:"id"  binding:"omitempty,number"`
	OperationID string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	UserId      string `json:"user_id" binding:"required,min=1" msg:"user_id不能为空"`
	GreetMsg    string `json:"greet_msg" binding:"omitempty,min=1,max=200" msg:"请填写好友打招呼消息，最多200字"`
	Remarks     string `json:"remarks"  binding:"omitempty,min=1,max=300" msg:"请填写备注，最多300字"`
}

type DefaultFriendDeleteReq struct {
	Id          int64  `json:"id"  binding:"required,number" msg:"id必须传"`
	UserId      string `json:"user_id"  binding:"omitempty,min=1" msg:"请传user_id"`
	OperationID string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
}

type DefaultFriendListReq struct {
	OperationID string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	Account     string `json:"account"  binding:"omitempty,min=1" msg:"请输入帐号"`
	NickName    string `json:"nick_name"  binding:"omitempty,min=1" msg:"请输入昵称"`
	UserId      string `json:"user_id"  binding:"omitempty,min=1" msg:"请传user_id"`
	Remarks     string `json:"remarks"  binding:"omitempty,min=1,max=300" msg:"请填写备注，最多300字"`
	pagination.Pagination
}

type DefaultFriendListInfo struct {
	ID          int64  `json:"id"`
	Account     string `json:"account"`
	NickName    string `json:"nick_name"`
	Remarks     string `json:"remarks"`
	FriendTotal int64  `json:"friend_total"`
	GreetMsg    string `json:"greet_msg"`
	UserId      string `json:"user_id"`
}

type DefaultFriendListResp struct {
	List     []DefaultFriendListInfo `json:"list"`
	Count    int64                   `json:"count"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}

type SignLogConfigResp struct {
	Sign int `json:"sign" binding:"required,number"`
}

type SignLogConfigReq struct {
	OperationID string            `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Config      SignLogConfigResp `json:"config" binding:"required,min=3" msg:"配置不能为空"`
}

type GetSignLogConfigReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type SignLogListReq struct {
	OperationID string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	Id          int64  `json:"id" binding:"omitempty,gte=1,number" msg:"请输入id"`
	UserId      string `json:"user_id" binding:"omitempty,number,min=1" msg:"请输入user_id"`
	NickName    string `json:"nick_name" binding:"omitempty,min=1" msg:"请输入昵称"`
	BeginDate   int64  `json:"begin_date" binding:"omitempty,number" msg:"请输入开始时间"`
	EndDate     int64  `json:"end_date"  binding:"omitempty,number" msg:"请输入结束时间"`
	pagination.Pagination
}
type SignLogInfo struct {
	Id        int64  `json:"id"`
	NickName  string `json:"nick_name"`
	UserId    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	Reward    int64  `json:"reward"`
}

type SignLogListResp struct {
	List []SignLogInfo `json:"list"`
	pagination.Pagination
}

type GetCmsConfigReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type GetCmsConfigResp struct {
	GoogleCodeIsOpen int                 `json:"google_code_is_open"`
	UIInfo           CmsSiteInfoItemResp `json:"ui_info"`
}

type SetGoogleCodeIsOpenReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Status      int    `json:"status" binding:"required,min=1,max=2" msg:"status 1开 2关"`
}

type CmsSiteInfoItemResp struct {
	SiteName     string `json:"site_name"`
	LoginIcon    string `json:"login_icon"`
	LoginBackend string `json:"login_backend"`
	PageIcon     string `json:"page_icon"`
	MenuIcon     string `json:"menu_icon"`
}

type ShieldDeleteReq struct {
	Id          int64  `json:"id"  binding:"required,number"`
	OperationID string `json:"operation_id"  binding:"required,min=1"`
}

type ShieldListReq struct {
	OperationID   string `json:"operation_id"  binding:"required,min=1"`
	ShieldWords   string `json:"shield_words" binding:"omitempty,min=1"`
	OperationUser string `json:"operation_user" binding:"omitempty,min=1"`
	Status        int64  `json:"status" binding:"omitempty,oneof=1 2"`
	DeleteStatus  int64  `json:"delete_status" binding:"omitempty,oneof=1 2"`
	BeginDate     int64  `json:"begin_date" binding:"omitempty,number"`
	EndDate       int64  `json:"end_date"  binding:"omitempty,number"`
	pagination.Pagination
}

type ShieldListInfo struct {
	Id            int64  `json:"id"`
	ShieldWords   string `json:"shield_words"`
	OperationUser string `json:"operation_user"`
	CreatedAt     int64  `json:"created_at"`
}

type ShieldListResp struct {
	List []ShieldListInfo `json:"list"`
	pagination.Pagination
}
type ShieldReq struct {
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	ID          int64  `json:"id" binding:"omitempty,number"`
	ShieldWords string `json:"shield_words" binding:"required,min=1"`
	Status      int64  `json:"status" binding:"omitempty,oneof=1 2"`
}

type GetJPushReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"`
}

type GetJPushResp struct {
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

type SetJPushReq struct {
	OperationID  string `json:"operation_id" binding:"required,min=1"`
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

type SetJPushResp struct {
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

type ParameterConfigResp struct {
	IpRegLimitCount         int64   `json:"ip_reg_limit_count" default:"0"`
	IpRegLimitTime          float64 `json:"ip_reg_limit_time" default:"0"`
	DeviceRegLimit          int64   `json:"device_reg_limit" default:"10" zero:"0"`
	GroupLimit              int64   `json:"group_limit" default:"1000" zero:"0"`
	ContactsFriendLimit     int64   `json:"contacts_friend_limit" default:"1000" zero:"0"`
	CreateGroupLimit        int64   `json:"create_group_limit" default:"1000" zero:"0"`
	UserIdPrefix            string  `json:"user_id_prefix" default:"用户" zero:""`
	FileSizeLimit           int64   `json:"file_size_limit" default:"100" zero:"0"`
	HistoryTime             int64   `json:"history_time" default:"30" zero:"0"`
	RevokeTime              int64   `json:"revoke_time" default:"0"`
	IsMemberDelFriend       int64   `json:"is_member_del_friend" default:"2"`
	IsMemberAddGroup        int64   `json:"is_member_add_group" default:"2"`
	IsOpenRegister          int64   `json:"is_open_register" default:"1"`
	IsShowMemberStatus      int64   `json:"is_show_member_status" default:"1"`
	IsMsgReadStatus         int64   `json:"is_msg_read_status" default:"1"`
	IsOpenRedPack           int64   `json:"is_open_redpack" default:"1"`
	IsOpenRedPackSingle     int64   `json:"is_open_redpack_single" default:"1"`
	IsNormalSeeAddress      int64   `json:"is_normal_see_address" default:"2"`
	IsOpenVoiceCall         int64   `json:"is_open_voice_call" default:"1"`
	IsOpenCameraCall        int64   `json:"is_open_camera_call" default:"1"`
	IsNormalSeeId           int64   `json:"is_normal_see_id" default:"2"`
	IsMemberAddFriend       int64   `json:"is_member_add_friend" default:"2"`
	IsNormalAddPrivilege    int64   `json:"is_normal_add_privilege" default:"1"`
	IsAddNormalVerify       int64   `json:"is_add_normal_verify" default:"1"`
	IsAddPrivilegeVerify    int64   `json:"is_add_privilege_verify" default:"1"`
	IsPrivilegeAddVerify    int64   `json:"is_privilege_add_verify" default:"2"`
	IsNormalJoinGroup       int64   `json:"is_normal_join_group" default:"1"`
	IsNormalMulSend         int64   `json:"is_normal_mul_send" default:"2"`
	IsShowRevoke            int64   `json:"is_show_revoke" default:"1"`
	IsPictureOpen           int64   `json:"is_picture_open" default:"1"`
	IsDisplayNicknameOpen   int64   `json:"is_display_nickname_open" default:"1"`
	IsOssOpen               int64   `json:"is_oss_open" default:"1"`
	SignAward               int64   `json:"sign_award" default:"0"`
	ShowGroupHistoryChatNum int64   `json:"show_group_history_chat_num" default:"0"`
	IsOpenNews              int64   `json:"is_open_news" default:"1"`
	IsOpenSign              int64   `json:"is_open_sign" default:"1"`
	IsOpenPlay              int64   `json:"is_open_play" default:"1"`
	IsOpenOperator          int64   `json:"is_open_operator" default:"1"`
	IsOpenShop              int64   `json:"is_open_shop" default:"1"`
	IsOpenTopStories        int64   `json:"is_open_top_stories" default:"1"`
	IsOpenWallet            int64   `json:"is_open_wallet" default:"1"`
	IsOpenPaymentCode       int64   `json:"is_open_payment_code" default:"1"`
}

type ParameterConfigReq struct {
	OperationID string              `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Config      ParameterConfigResp `json:"config" binding:"required" msg:"配置不能为空"`
}

type PrivilegeUserReq struct {
	OperationID string `json:"operation_id"  binding:"required,min=1" msg:"日志id必须传"`
	Account     string `json:"account" binding:"omitempty,min=1" msg:"请输入帐号"`
	NickName    string `json:"nick_name" binding:"omitempty,min=1" msg:"请输入昵称"`
	UserId      string `json:"user_id" binding:"omitempty,min=1,number" msg:"请输入user_id"`
	pagination.Pagination
}
type PrivilegeUserListInfo struct {
	UserID      string `json:"user_id"`
	NickName    string `json:"nick_name"`
	Account     string `json:"account"`
	Status      int64  `json:"status"`
	PhoneNumber string `json:"phone_number"`
}
type PrivilegeUserListResp struct {
	List     []PrivilegeUserListInfo `json:"list"`
	IsFreeze int                     `json:"is_freeze"`
	pagination.Pagination
}

type GetDepositeReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"`
}

type GetDepositeResp struct {
	Html   string `json:"html"`
	Url    string `json:"url" binding:"omitempty,url" msg:"请输入有效链接"`
	Switch int    `json:"switch" binding:"required,oneof=1 2" msg:"二选一 自定义网页或充值地址"`
}

type SetDepositeReq struct {
	OperationID string `json:"operation_id" binding:"required,min=1"`
	Html        string `json:"html"`
	Url         string `json:"url" binding:"omitempty,url" msg:"请输入有效链接"`
	Switch      int    `json:"switch" binding:"required,oneof=1 2" msg:"二选一 自定义网页或充值地址"`
}

type GetWithdrawConfigReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"`
}

type WithdrawConfigResp struct {
	Min     int64                    `json:"min"`
	Max     int64                    `json:"max"`
	Columns WithdrawConfigColumnList `json:"columns"`
}

type WithdrawApigResp struct {
	Min     int64                    `json:"min"`
	Max     int64                    `json:"max"`
	Columns WithdrawConfigColumnList `json:"columns"`
}

type WithdrawConfigColumn struct {
	Name                 string `json:"name" default:"提现金额"`
	DefaultContent       string `json:"default_content"`
	Required             int    `json:"required" default:"1"`
	DefaultContentModify int    `json:"default_content_modify" default:"1"`
	Sort                 int    `json:"sort" default:"1"`
	Value                string `json:"value"`
}

type WithdrawConfigColumnList []WithdrawConfigColumn

func (m WithdrawConfigColumnList) Len() int {
	return len(m)
}

func (m WithdrawConfigColumnList) Less(i, j int) bool {
	return m[i].Sort > m[j].Sort
}

func (m WithdrawConfigColumnList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type SetWithdrawConfigReq struct {
	OperationID string `json:"operation_id" binding:"required,min=1"`
	WithdrawConfigResp
}

type GetFeihuReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"`
}

type GetFeihuResp struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

type SetFeihuReq struct {
	OperationID string `json:"operation_id" binding:"required,min=1"`
	AppKey      string `json:"app_key"`
	AppSecret   string `json:"app_secret"`
}

type SetFeihuResp struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

type GetAboutUsReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"`
}

type GetAboutUsResp struct {
	Content string `json:"content"`
}

type SetAboutUsReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"`
	Content     string `json:"content" form:"content" binding:"required"`
}

type SetIPWhiteListIsOpenReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Status      int    `json:"status" binding:"required,min=1,max=2" msg:"status 1开 2关"`
}

type GetIPWhiteListIsOpenReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}

type SetDefaultIsOpenReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Status      int    `json:"status" binding:"required,min=1,max=2" msg:"status 1开 2关"`
}

type GetDefaultIsOpenReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
}
