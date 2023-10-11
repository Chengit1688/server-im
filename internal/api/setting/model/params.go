package model

import "im/pkg/pagination"

type OperationIDReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,gte=1"  msg:"日志id必须传递"`
}

type RegisterConfigInfo struct {
	IsInviteCode       int64 `json:"is_invite_code" form:"is_invite_code"`
	CheckInviteCode    int64 `json:"check_invite_code" form:"check_invite_code"`
	IsVerificationCode int64 `json:"is_verification_code" form:"is_verification_code"`
	IsSmsCode          int64 `json:"is_sms_code" form:"is_sms_code"`
	IsAllAccount       int64 `json:"is_all_account" form:"is_all_account"`
}

type VersionReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"日志id必须传递"`
	Platform    int64  `json:"platform" form:"platform" binding:"required,oneof=1 2 3 4 5 6 7 8 9"`
}

type VersionResp struct {
	Version     string `json:"version" form:"version"  binding:"required,gte=1"`
	IsForce     int64  `json:"is_force" form:"is_force"  binding:"required,oneof=1 2"`
	Title       string `json:"title" form:"title"  binding:"required,min=1,max=50"`
	DownloadUrl string `json:"download_url" form:"download_url"  binding:"required,url,min=1"`
	UpdateDesc  string `json:"update_desc" form:"update_desc"  binding:"omitempty,max=200"`
}

type VersionConfig struct {
	Platform    int64  `json:"platform" form:"platform"`
	Version     string `json:"version" form:"version"`
	IsForce     int64  `json:"is_force" form:"is_force"`
	Title       string `json:"title" form:"title"`
	DownloadUrl string `json:"download_url" form:"download_url"`
	UpdateDesc  string `json:"update_desc" form:"update_desc"`
}

type AboutUsReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"日志id必须传递"`
}

type AboutUsResp struct {
	Content string `json:"content" form:"content"`
}

type imgSize struct {
	Width  string
	Height string
}

type CaptchaReq struct {
	OperationID string  `json:"operation_id" form:"operation_id"  binding:"required,min=1" msg:"日志id必须传递"`
	CaptchaType string  `json:"captcha_type" form:"captcha_type"  binding:"required,oneof=blockPuzzle clickWord" msg:"验证码类型必须传递"`
	ClientUid   string  `json:"client_uid" form:"captcha_type"  binding:"omitempty" msg:"请传递client_uid传递"`
	Mode        string  `json:"mode" form:"mode"`
	VSpace      string  `json:"v_space" form:"v_space"`
	Explain     string  `json:"explain" form:"explain"`
	ImgSize     imgSize `json:"imgSize" form:"img_size"`
}

type CheckCaptchaReq struct {
	OperationID string `json:"operation_id" form:"operation_id"  binding:"required,min=1" msg:"日志id必须传递"`
	CaptchaType string `json:"captcha_type" form:"captcha_type"  binding:"required,oneof=blockPuzzle clickWord" msg:"验证码类型必须传递"`
	PointJson   string `json:"point_json" form:"point_json"  binding:"required,gte=1" msg:"point_json必须传递"`
	Token       string `json:"token" form:"token"  binding:"required,gte=1" msg:"token必须传递"`
}

type SmsReq struct {
	OperationID       string `json:"operation_id" form:"operation_id"  binding:"required,gte=1" msg:"日志id必须传递"`
	PhoneNumber       string `json:"phone_number" form:"phone_number" binding:"required,min=3,max=11"  msg:"请输入正确的手机号"`
	CountryCode       string `json:"country_code" form:"country_code" binding:"required,min=2"  msg:"输入选择正确的国家编码"`
	UsedFor           int64  `json:"used_for" form:"used_for" binding:"required,min=1,oneof=1 2 3"  msg:"输入used_for"`
	VerificationToken string `json:"verification_token"  form:"verification_token" binding:"omitempty,min=1"  msg:"请选择验证码"`
	VerificationPoint string `json:"verification_point"  form:"verification_point" binding:"omitempty,min=1"  msg:"请选择验证码"`
	CaptchaType       string `json:"captcha_type"  form:"captcha_type" binding:"omitempty,oneof=blockPuzzle clickWord"  msg:"请选择验证码"`
}

type GetDiscoverInfoResp struct {
	IsOpen int            `json:"is_open"`
	List   []DiscoverInfo `json:"list"`
}

type DiscoverInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int    `json:"sort"`
	Url  string `json:"url"`
}

type ShieldListReq struct {
	OperationID string `json:"operation_id"  binding:"required,min=1"`
	pagination.Pagination
}

type ShieldListResp struct {
	List []string `json:"list"`
	pagination.Pagination
}
