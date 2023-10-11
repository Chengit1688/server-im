package model

import "im/pkg/pagination"

type Op struct {
	OperationID string `json:"operation_id" binding:"required"`
}

type ListReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GroupInfo struct {
	Group
	Role          RoleType `json:"role"`
	GroupNickName string   `json:"group_nick_name"`
}

type MaxGroupSeq struct {
	ConversationID string `json:"conversation_id"`
	MaxGroupSeq    int64  `json:"max_group_seq`
}

type GroupUpdateReq struct {
	Op
	GroupId                    string `json:"group_id"  binding:"required"`
	Name                       string `json:"name"`
	FaceUrl                    string `json:"face_url"`
	Notification               string `json:"notification"`
	Introduction               string `json:"introduction"`
	Status                     int    `json:"status"`
	NoShowNormalMember         int    `json:"no_show_normal_member"`
	NoShowAllMember            int    `json:"no_show_all_member"`
	ShowQrcodeByNormalMember   int    `json:"show_qrcode_by_normal"`
	ShowQrcodeByNormalMemberV2 int    `json:"show_qrcode_by_normal_member_v2"`
	JoinNeedApply              int    `json:"join_need_apply"`
	BanRemoveByNormalMember    int    `json:"ban_remove_by_normal"`
	MuteAllMember              int    `json:"mute_all_member"`
	IsTopannocuncement         int    `json:"is_topannocuncement"`
	IsOpenAdminList            int    `json:"is_open_admin_list"`
	IsOpenAdminIcon            int    `json:"is_open_admin_icon"`
	IsOpenGroupId              int    `json:"is_open_group_id"`
	GroupSendLimit             int    `json:"group_send_limit"`
	IsDisplayNicknameOpen      int    `json:"is_display_nickname_open"`
}

type GroupUpdateAvatarReq struct {
	Op
	GroupId string `json:"group_id"  binding:"required"`
	FaceUrl string `json:"face_url"`
	Name    string `json:"name"`
}
type GroupUpdateResp GroupInfo

type GroupUpdateAvatarResp GroupInfo

type GroupMemberInfo struct {
	Id            int64    `json:"id"`
	GroupId       string   `json:"group_id"`
	UserId        string   `json:"user_id"`
	GroupNickName string   `json:"group_nick_name"`
	Role          RoleType `json:"role"`
	MuteEndTime   int64    `json:"mute_end_time"`
	NickName      string   `json:"nick_name"`
	FaceUrl       string   `json:"face_url"`
	BigFaceUrl    string   `json:"big_face_url"`
	Version       int      `json:"version"`
	Status        int      `json:"status"`
	Account       string   `json:"account"`
}

type GetAdminOwnerResp struct {
	List []GroupMemberInfo
}

type GroupInformationInfo struct {
	Name               *string `json:"name"`
	FaceUrl            *string `json:"face_url"`
	BigFaceUrl         *string `json:"big_face_url"`
	Notification       *string `json:"notification"`
	Introduction       *string `json:"introduction"`
	IsTopannocuncement *int    `json:"is_topannocuncement"`
}

type GroupManageInfo struct {
	NoShowNormalMember       *int `json:"no_show_normal_member"`
	NoShowAllMember          *int `json:"no_show_all_member"`
	ShowQrcodeByNormalMember *int `json:"show_qrcode_by_normal_member"`
	JoinNeedApply            *int `json:"join_need_apply"`
	BanRemoveByNormalMember  *int `json:"ban_remove_by_normal_member"`
	IsDefault                *int `json:"is_default"`
	IsOpenAdminList          *int `json:"is_open_admin_list"`
	IsOpenAdminIcon          *int `json:"is_open_admin_icon"`
	IsOpenGroupId            *int `json:"is_open_group_id"`
	GroupSendLimit           *int `json:"group_send_limit"`
}

type GroupSearchReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	pagination.Pagination
	Keyword string `json:"keyword" binding:"required"`
}

type GroupSearchResp struct {
	pagination.Pagination
	List []GroupInfo `json:"list"`
}

type CreateGroupReq struct {
	Op
	Name         string `json:"name"  binding:"required"`
	FaceUrl      string `json:"face_url"`
	Introduction string `json:"introduction"`
	Notification string `json:"notification"`
}

type CreateGroupResp GroupInfo

type Face2FaceAddReq struct {
	Op
	GroupNumber string `json:"group_number"  binding:"required,min=6,max=6"`
}
type Face2FaceAddResp struct {
	Group
}

type UserInfo struct {
	UserID     string `json:"user_id"`
	Account    string `json:"account"`
	NickName   string `json:"nick_name"`
	FaceURL    string `json:"face_url"`
	BigFaceURL string `json:"big_face_url"`
}

type Face2FaceInviteResp struct {
	Users []UserInfo `json:"users"`
}

type JoinGroupApplyReq struct {
	Op
	GroupId string `json:"group_id"  binding:"required"`
	Remark  string `json:"remark"`
}

type JoinGroupApplyResp struct {
}

type JoinGroupVerifyReq struct {
	Op
	ApplyId int64 `json:"apply_id"  binding:"required"`
	Status  int   `json:"status"  binding:"required"`
}

type JoinGroupVerifyResp struct {
}

type JoindGroupListReq struct {
	Op
	ListReq
	Version int `json:"version"`
}

type JoinApplyListReq struct {
	Op
	ListReq
}

type GroupListSyncReq struct {
	Op
	GroupIdList []string `json:"group_id_list"`
}

type GroupListSyncResp struct {
	Count    int         `json:"count"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     []GroupInfo `json:"list"`
}

type ApplyInfo struct {
	Id             int64  `gorm:"column:id;primary_key;" json:"id"`
	GroupName      string `json:"group_name"`
	GroupId        string `gorm:"column:group_id;index;type:varchar(255);" json:"group_id"`
	UserId         string `gorm:"column:user_id;index;type:varchar(255);" json:"user_id"`
	Remark         string `gorm:"column:remark;type:varchar(255)" json:"remark"`
	CreateTime     int64  `gorm:"column:create_time" json:"create_time"`
	Status         int    `gorm:"column:status;default:0;type:tinyint(3)" json:"status"`
	OperationTime  int64  `gorm:"operator_time" json:"operation_time"`
	OperatorUserId string `gorm:"operator_user_id;type:varchar(255)" json:"operator_user_id"`
	NickName       string `json:"nick_name"`
	FaceUrl        string `json:"face_url"`
	BigFaceUrl     string `json:"big_face_url"`
}

type JoinApplyListResp struct {
	Count    int         `json:"count"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     []ApplyInfo `json:"list"`
}

type JoindGroupListResp struct {
	Count    int         `json:"count"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     []GroupInfo `json:"list"`
}

type MyGroupListReq struct {
	Op
	ListReq
	Version int `json:"version"`
}
type MyGroupListMaxSeqResp struct {
	Count    int           `json:"count"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	List     []MaxGroupSeq `json:"list"`
}

type MyGroupListResq struct {
	Count    int         `json:"count"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     []GroupInfo `json:"list"`
}

type GroupInfoReq struct {
	Op
	GroupId string `json:"group_id"  binding:"required"`
}

type GroupInfoResp GroupInfo

type GroupInformationReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	GroupID     string `json:"group_id"`
	GroupInformationInfo
}

type GroupInformationResp struct {
	Group
}

type GroupManageReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	GroupID     string `json:"group_id"`
	GroupManageInfo
}

type GroupManageResp struct {
	Group
}

type GroupRemoveReq struct {
	OperationID string `json:"operation_id" binding:"required"`
	GroupID     string `json:"group_id"`
}

type GroupRemoveResp struct {
	GroupID string `json:"group_id"`
}

type GroupMemberListReq struct {
	Op
	ListReq
	GroupId   string `json:"group_id"`
	SearchKey string `json:"search_key"`
	IsMute    int    `json:"is_mute"`
}

type UpdateGroupMemberReq struct {
	Op
	GroupID       string `json:"group_id"  binding:"required"`
	UserId        string `json:"user_id"`
	GroupNickName string `json:"group_nick_name"`
}

type GroupMemberListResp struct {
	Count    int               `json:"count"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	List     []GroupMemberInfo `json:"list"`
}

type QuitGroupReq struct {
	Op
	GroupId string `json:"group_id"  binding:"required"`
}

type QuitGroupResp struct{}

type RemoveGroupMemberReq struct {
	Op
	GroupId    string   `json:"group_id"  binding:"required"`
	UserIdList []string `json:"user_id_list"  binding:"required"`
}

type RemoveGroupMemberResp struct {
}

type InviteGroupMemberReq struct {
	Op
	GroupId    string   `json:"group_id"  binding:"required"`
	UserIdList []string `json:"user_id_list"  binding:"required"`
}

type InviteGroupMemberResp struct {
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Count    int               `json:"count"`
	List     []GroupMemberInfo `json:"list"`
}

type GroupSyncReq struct {
	Op
	GroupId      string `json:"group_id"  binding:"required"`
	LocalVersion int    `json:"local_version"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
}

type GroupSyncResp GroupMemberListResp

type GroupSetAdminReq struct {
	Op
	GroupId string `json:"group_id"  binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
	Status  int    `json:"status" binding:"required"`
}

type GroupSetAdminResp struct{}

type GroupSetOwnerReq struct {
	Op
	GroupId string `json:"group_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

type GetOwnerAdminReq struct {
	Op
	GroupId string `json:"group_id" binding:"required"`
}

type GroupSetOwnerResp struct{}

type GroupMuteMemberReq struct {
	Op
	GroupId string `json:"group_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
	MuteSec int    `json:"mute_sec"`
}

type GroupMuteMemberResp struct{}

type GroupsVersion struct {
	GroupId       string `json:"group_id"`
	GroupVersion  int    `json:"group_version"`
	MemberVersion int    `json:"member_version"`
}

type GroupMuteAllReq struct {
	Op
	GroupId       string `json:"group_id"  binding:"required"`
	MuteAllMember int    `json:"mute_all_member"`
	MuteAllPeriod string `json:"mute_all_period"`
}

type GroupMuteAllResp struct{}

type GroupNickNameReq struct {
	Op
	GroupId       string `json:"group_id"  binding:"required"`
	GroupNickName string `json:"group_nick_name"`
}

type GroupNickNameResp struct{}
