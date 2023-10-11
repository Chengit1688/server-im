package model

import "im/pkg/pagination"

type LoginReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"   msg:"日志id必须传递"`
	Account     string `json:"account" form:"account" binding:"omitempty,lte=16,gte=2" msg:"请正确输入帐号"`
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"omitempty,min=3,max=15"  msg:"输入正确的手机号和密码"`
	Password    string `json:"password" form:"password" binding:"omitempty,min=6,max=16" msg:"请正确输入密码"`
	Platform    int64  `json:"platform" form:"platform" binding:"required,oneof=1 2 3 4 5 6 7 8 9"  msg:"不支持此平台登录"`
	DeviceId    string `json:"device_id" form:"device_id" binding:"required,min=1" msg:"设备id必须传递"`
	LoginType   int64  `json:"login_type" form:"login_type" binding:"required,oneof=1 2 3"  msg:"只支持1/2/3"`
	LoginIp     string `json:"login_ip" form:"login_ip"`
	CountryCode string `json:"country_code"  form:"country_code" binding:"omitempty,min=2"  msg:"请选择国际编码"`
	Brand       string `json:"brand" form:"brand" msg:"请输入正确的品牌"`
}

type LoginResp struct {
	Token  string `json:"token" form:"token" `
	UserId string `json:"user_id"`
}

type RegisterReq struct {
	OperationID       string `json:"operation_id" form:"operation_id" binding:"required,gte=2"   msg:"日志id必须传递"`
	Account           string `json:"account" form:"account" binding:"omitempty,alphanum,min=2,max=16" msg:"请输入长度为2-16位的限英文数字账号"`
	PhoneNumber       string `json:"phone_number" form:"phone_number"  binding:"omitempty,min=3,max=15"  msg:"请正确输入手机号"`
	Password          string `json:"password"  form:"password" binding:"min=6,max=16" msg:"请输入密码，长度6-16位"`
	Platform          int64  `json:"platform" form:"platform"  binding:"required,oneof=1 2 3 4 5 6 7 8 9"`
	InviteCode        string `json:"invite_code" form:"invite_code" binding:"omitempty,min=1" msg:"请输入邀请码"`
	VerificationCode  string `json:"verification_code" form:"verification_code" binding:"omitempty,min=1" msg:"请输入验证码"`
	SmsCode           string `json:"sms_code" form:"sms_code"  binding:"omitempty,min=1" msg:"请输入短信验证码"`
	DeviceId          string `json:"device_id" form:"device_id" binding:"required,min=1"   msg:"设备id必须传递"`
	AccountType       int64  `json:"account_type"  form:"account_type" binding:"required,oneof=1 2"  msg:"只支持1/2"`
	VerificationToken string `json:"verification_token"  form:"verification_token" binding:"omitempty,min=1"  msg:"请验证图形验证码"`
	VerificationPoint string `json:"verification_point"  form:"verification_point" binding:"omitempty,min=1"  msg:"请验证图形验证码"`
	CaptchaType       string `json:"captcha_type"  form:"captcha_type" binding:"omitempty,oneof=blockPuzzle clickWord"  msg:"请验证图形验证码"`
	CountryCode       string `json:"country_code"  form:"country_code" binding:"omitempty,min=2"  msg:"请选择国际编码"`
	ImSite            string `json:"im_site"  form:"im_site"`
	Brand             string `json:"brand" form:"brand" msg:"请输入正确的品牌"`
}

type RegisterResp struct {
	LoginResp
}

type UserInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"   msg:"日志id必须传递"`
}

type UserInfoResp struct {
	UserId      string `json:"user_id" form:"user_id"`
	Account     string `json:"account" form:"account"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	FaceURL     string `json:"face_url" form:"face_url"`
	BigFaceURL  string `json:"big_face_url" form:"big_face_url"`
	Gender      int64  `json:"gender" form:"gender"`
	Platform    int64  `json:"platform" form:"platform"`
	DeviceId    string `json:"device_id" form:"device_id"`
	NickName    string `json:"nick_name" form:"nick_name"`
	Signatures  string `json:"signatures" form:"signatures"`
	Age         int64  `json:"age" form:"age"`
	IsPrivilege int64  `json:"is_privilege" form:"is_privilege"`
	InviteCode  string `json:"invite_code"`
}

type UserInfoUpdateReq struct {
	OperationID string  `json:"operation_id" form:"operation_id" binding:"required,min=1"  msg:"日志id必须传递"`
	NickName    *string `json:"nick_name" form:"nick_name" binding:"omitempty,min=1,max=16" msg:"请正确输入昵称"`
	Signatures  *string `json:"signatures" form:"signatures" binding:"omitempty,max=200" msg:"请输入长度小于200的签名"`
	Gender      *int64  `json:"gender" form:"gender" binding:"omitempty,oneof=1 2" msg:"请选择性别"`
	Age         *int64  `json:"age" form:"age" binding:"omitempty,number" msg:"请正确输入年龄"`
	FaceURL     *string `json:"face_url" form:"face_url" binding:"omitempty" msg:"请正确上传头像"`
	BigFaceURL  *string `json:"big_face_url" form:"big_face_url" binding:"omitempty" msg:"请正确上传头像"`
}

type PasswordSecureReq struct {
	OperationID      string `json:"operation_id" form:"operation_id" binding:"required,min=1"  msg:"日志id必须传递"`
	OriginalPassword string `json:"original_password" form:"original_password" binding:"omitempty,min=6,max=16"  msg:"请正确输入密码"`
	NewPassword      string `json:"new_password" form:"new_password" binding:"required,eqfield=ConfirmPassword,min=6,max=16"  msg:"请正确输入新密码"`
	ConfirmPassword  string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=NewPassword,min=6,max=16"  msg:"请再次确认新密码是否正确"`
	PasswordType     int64  `json:"password_type" form:"password_type" binding:"required,oneof=1 2 3" msg:"不支持此修改密码类型"`
}

type ForgotPasswordReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"   msg:"日志id必须传递"`
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required,min=3,max=15" msg:"请正确输入手机号"`
	CountryCode string `json:"country_code"  form:"country_code" binding:"required,min=2"  msg:"请选择国家编码"`

	NewPassword     string `json:"new_password" form:"new_password" binding:"required,eqfield=ConfirmPassword,min=6,max=16" msg:"请正确输入密码"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=NewPassword,min=6,max=16" msg:"请再次确认新密码"`
}

type VerifyPhoneCodeReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"   msg:"日志id必须传递"`
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required,min=3,max=15" msg:"请正确输入手机号"`
	CountryCode string `json:"country_code"  form:"country_code" binding:"required,min=2"  msg:"请选择国家编码"`
	SmsCode     string `json:"sms_code" form:"sms_code"  binding:"required,number" msg:"请正确输入短信验证码"`
}

type CurrentTokenReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"  msg:"日志id必须传递"`
	UserId      string `json:"user_id" form:"user_id"  binding:"required,number" msg:"user_id必传"`
}

type CurrentTokenResp struct {
	Token string `json:"token" form:"token" `
}

type UserBaseInfoResp struct {
	UserBaseInfo
}

type UserBaseInfo struct {
	UserId      string `json:"user_id"`
	Account     string `json:"account" form:"account"`
	FaceURL     string `json:"face_url"`
	BigFaceURL  string `json:"big_face_url"`
	Gender      int64  `json:"gender" form:"gender"`
	NickName    string `json:"nick_name" form:"nick_name"`
	Signatures  string `json:"signatures" form:"signatures"`
	Age         int64  `json:"age" form:"age"`
	IsPrivilege int64  `json:"is_privilege" form:"is_privilege"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	CountryCode string `json:"country_code" form:"country_code"`
}

type GetUserInfoReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"  msg:"日志id必须传递"`
	UserId      string `json:"user_id" binding:"required,number" msg:"user_id必传"`
}

type GetDeviceInfo struct {
	PlatformClass string `json:"platform_class"`
	DeviceName    string `json:"device_name"`
	DeviceId      int64  `json:"device_id"`
}

type GetServerVersionReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"  msg:"日志id必须传递"`
}

type GetServerVersionRsp struct {
	FriendVersion       int32              `json:"friend_version"`
	ConversationVersion int64              `json:"conversation_version"`
	GroupsVerison       []GroupVersionInfo `json:"groups_version"`
}

type GroupVersionInfo struct {
	GroupId       string `json:"group_id"`
	GroupVersion  int    `json:"group_version"`
	MemberVersion int    `json:"member_version"`
}

type GetUserSignInfoResp struct {
	Total    int   `json:"total"`
	Today    bool  `json:"today"`
	SignOpen bool  `json:"sign_open"`
	Days     []int `json:"days"`
}

type GetUserSignInfoV2Resp struct {
	Total     int64    `json:"total"`
	Today     bool     `json:"today"`
	SignOpen  bool     `json:"sign_open"`
	Balance   int64    `json:"balance"`
	SignAward int64    `json:"sign_award"`
	Days      []string `json:"days"`
}

type SignTodayReq struct {
	SignDate    string `json:"sign_date" form:"sign_date"  binding:"omitempty,min=1" msg:"请传递签到日期"`
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"  msg:"日志id必须传递"`
}

type UserConfigHandleReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"  msg:"日志id必须传递"`
	Content     string `json:"content" form:"content" binding:"required,min=1"  msg:"配置内容必须传递"`
}

type GetUserConfigReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"  msg:"日志id必须传递"`
	Content     string `json:"content" form:"content" binding:"required,min=1"  msg:"配置内容必须传递"`
	Version     int64  `json:"version" form:"version" binding:"required,number"  msg:"版本号必须传递"`
}
type GetUserConfigResp struct {
	Content string `json:"content"`
	Version int64  `json:"version"`
}

type GetUserOnlineStatusReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"  msg:"日志id必须传递"`
	UserId      string `json:"user_id" from:"user_id" binding:"required,min=1"  msg:"用户id必须传递"`
}

type GetUserOnlineStatusResp struct {
	UserId      string `json:"user_id"`
	Online      bool   `json:"online"`
	Ip          string `json:"ip"`
	IpAddress   string `json:"ip_address"`
	OfflineInfo string `json:"offline_info"`
}

type UserConfigPushResp struct {
	Version int64  `json:"version"`
	Content string `json:"content"`
}

type SuggestionReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,min=1"  msg:"日志id必须传递"`
	Content     string `json:"content" form:"content" binding:"required,min=1"  msg:"请输入投诉建议内容"`
	Brand       string `json:"brand" form:"brand" binding:"omitempty,min=1"  msg:"请输入正确的品牌"`
	Platform    int64  `json:"platform" form:"platform" binding:"required,min=1"  msg:"请传递正确的客户端类型"`
	AppVersion  string `json:"app_version"  form:"app_version" binding:"omitempty,min=1"  msg:"请传递正确的app版本"`
}

type GetFavoriteImageReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"   msg:"日志id必须传递"`
	pagination.Pagination
}

type GetFavoriteImageResp struct {
	List []FavoriteImageItem `json:"list"`
	pagination.Pagination
}

type FavoriteImageItem struct {
	UUID           string `json:"uuid"`
	ImageUrl       string `json:"image_url"`
	ImageThumbnail string `json:"image_thumbnail"`
	ImageWidth     *int   `json:"image_width"`
	ImageHeight    *int   `json:"image_height"`
}

type AddFavoriteImageReq struct {
	OperationID    string `json:"operation_id" form:"operation_id" binding:"required,gte=2"   msg:"日志id必须传递"`
	UUID           string `json:"uuid" binding:"required"`
	ImageUrl       string `json:"image_url" binding:"required"`
	ImageThumbnail string `json:"image_thumbnail" binding:"required"`
	ImageWidth     int    `json:"image_width" binding:"required"`
	ImageHeight    int    `json:"image_height" binding:"required"`
}

type AddFavoriteImageResp struct {
	FavoriteImageItem
}

type DelFavoriteImageReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required,gte=2"   msg:"日志id必须传递"`
	UUID        string `json:"uuid" binding:"required"`
}

type SetPrivacyReq struct {
	OperationID string  `json:"operation_id"  binding:"required,gte=2"`
	Data        Privacy `json:"data"`
}

type Privacy struct {
	IsShowOnlineTime  int64 `json:"is_show_online_time" default:"1" binding:"required,oneof=1 2"`
	IsFriendVerify    int64 `json:"is_friend_verify" default:"1"  binding:"required,oneof=1 2"`
	IsMobileSearch    int64 `json:"is_mobile_search" default:"1"  binding:"required,oneof=1 2"`
	IsIDSearch        int64 `json:"is_id_search" default:"1"  binding:"required,oneof=1 2"`
	IsFromGroupFriend int64 `json:"is_from_group_friend"  default:"1" binding:"required,oneof=1 2"`
	IsNicknameSearch  int64 `json:"is_nickname_search" default:"1"  binding:"required,oneof=1 2"`
	IsMessageVibrate  int64 `json:"is_message_vibrate" default:"1"  binding:"required,oneof=1 2"`
}

type SignInReq struct {
	OperationID string `json:"operation_id"  binding:"required,gte=2"`
}

type RedeemPrizeReq struct {
	OperationID string `json:"operation_id"  binding:"required,gte=2"`
	PrizeID     int64  `json:"prize_id"  binding:"required,number"`
	UserName    string `json:"user_name"`
	Address     string `json:"address"`
	Mobile      string `json:"mobile"  binding:"required,gte=2"`
}

type RedeemPrizeListReq struct {
	pagination.Pagination
	OperationID   string `json:"operation_id"  binding:"required,gte=2"`
	Key           string `json:"key"`
	UserID        string `json:"-"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	ExpressNumber string `json:"express_number"`
}

type GetPrizeListReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	pagination.Pagination
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Name      string `json:"name"`
	Cost      int64  `json:"cost"`
	Describe  string `json:"describe"`
	Status    int64  `json:"-"`
}

type RealNameReq struct {
	OperationID string `json:"operation_id" binding:"required" msg:"操作ID不能为空"`
	RealName    string `json:"real_name" binding:"required"`
	IDNo        string `json:"id_no" binding:"required"`
	IDFrontImg  string `json:"id_front_img" binding:"required"`
	IDBackImg   string `json:"id_back_img" binding:"required"`
}

type RealNameResp struct {
	RealName    string `json:"real_name"`
	IDNo        string `json:"id_no"`
	IDFrontImg  string `json:"id_front_img"`
	IDBackImg   string `json:"id_back_img"`
	RealAuthMsg string `json:"real_auth_msg"`
	IsRealAuth  int    `json:"is_real_auth"`
}

type RealNameInfoReq struct {
	OperationID string `json:"operation_id"  binding:"required,gte=1"`
}
