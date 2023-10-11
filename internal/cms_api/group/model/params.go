package model

type Op struct {
	OperationID string `json:"operation_id" binding:"required"`
}

type ListReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GroupInfo struct {
	Id                       int64  `json:"id"`
	GroupId                  string `json:"group_id"`
	Name                     string `json:"name"`
	FaceUrl                  string `json:"face_url"`
	MembersTotal             int    `json:"members_total"`
	RobotTotal               int    `json:"robot_total"`
	Notification             string `json:"notification"`
	Introduction             string `json:"introduction"`
	CreateTime               int64  `json:"create_time"`
	LastVersion              int    `json:"last_version"`
	CreateUserId             string `json:"create_user_id"`
	Status                   int    `json:"status"`
	NoShowNormalMember       int    `json:"no_show_normal_member"`
	NoShowAllMember          int    `json:"no_show_all_member"`
	ShowQrcodeByNormalMember int    `json:"show_qrcode_by_normal"`
	JoinNeedApply            int    `json:"join_need_apply"`
	BanRemoveByNormalMember  int    `json:"ban_remove_by_normal"`
	MuteAllMember            int    `json:"mute_all_member"`
	OwnerNickName            string `json:"owner_nick_name"`
	OwnerUserId              string `json:"owner_user_id"`
	IsDefault                int    `json:"is_default"`
	GroupSendLimit           int    `json:"group_send_limit"`
	MuteAllPeriod            string `json:"mute_all_period"`
	IsTopannocuncement       int    `json:"is_topannocuncement"`
}

type GroupMemberInfo struct {
	Id            int64  `json:"id"`
	GroupId       string `json:"group_id"`
	UserId        string `json:"user_id"`
	GroupNickName string `json:"group_nick_name"`
	Role          string `json:"role"`
	MuteEndTime   int64  `json:"mute_end_time"`
	NickName      string `json:"nick_name"`
	FaceUrl       string `json:"face_url"`
	Version       int    `json:"version"`
	Status        int    `json:"status"`
	Account       string `json:"account"`
}

type CreateGroupReq struct {
	Op
	Name               string `json:"name"  binding:"required"`
	FaceUrl            string `json:"face_url"`
	OwnerId            string `json:"owner_id" binding:"required"`
	Introduction       string `json:"introduction"`
	Notification       string `json:"notification"`
	IsTopannocuncement int    `json:"is_topannocuncement"`
}

type CreateGroupResp GroupInfo
type GroupListReq struct {
	Op
	ListReq
	GroupName string `json:"group_name"`
	OwnerName string `json:"owner_name"`
	IsDefault int    `json:"is_default"`
}

type GroupListResq struct {
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

type GroupUpdateReq struct {
	Op
	GroupId                  string `json:"group_id"  binding:"required"`
	Name                     string `json:"name"`
	FaceUrl                  string `json:"face_url"`
	Notification             string `json:"notification"`
	Introduction             string `json:"introduction"`
	Status                   int    `json:"status"`
	NoShowNormalMember       int    `json:"no_show_normal_member"`
	NoShowAllMember          int    `json:"no_show_all_member"`
	ShowQrcodeByNormalMember int    `json:"show_qrcode_by_normal"`
	JoinNeedApply            int    `json:"join_need_apply"`
	BanRemoveByNormalMember  int    `json:"ban_remove_by_normal_"`
	MuteAllMember            int    `json:"mute_all_member"`
	IsDefault                int    `json:"is_default"`
	GroupSendLimit           int    `json:"group_send_limit"`
	IsTopannocuncement       int    `json:"is_topannocuncement"`
	MuteAllPeriod            string `json:"mute_all_period"`
}

type GroupUpdateResp struct{}

type GroupRobotUpdateReq struct {
	Op
	GroupId    string `json:"group_id"  binding:"required"`
	RobotTotal int    `json:"robot_total"`
}

type GroupRobotUpdateResp struct{}

type GroupMemberListReq struct {
	Op
	ListReq
	GroupId   string `json:"group_id"`
	SearchKey string `json:"search_key"`
}

type GroupMemberListResp struct {
	Count    int               `json:"count"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	List     []GroupMemberInfo `json:"list"`
}

type RemoveGroupMemberReq struct {
	Op
	GroupId    string   `json:"group_id"  binding:"required"`
	UserIdList []string `json:"user_id_list"  binding:"required"`
}

type RemoveGroupMemberResp struct {
}

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

type GroupSetOwnerResp struct{}

type GroupMuteMemberReq struct {
	Op
	GroupId string `json:"group_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
	MuteSec int    `json:"mute_sec"`
}

type GroupMuteMemberResp struct{}

type AddGroupMembersReq struct {
	Op
	GroupId    string   `json:"group_id" binding:"required"`
	UserIdList []string `json:"user_id_list" binding:"required"`
}

type AddGroupMembersResp struct{}

type GroupMergeReq struct {
	Op
	FromGroupId string `json:"from_group_id" binding:"required"`
	ToGroupId   string `json:"to_group_id" binding:"required"`
}

type GroupMergeResp GroupInfo

type JoinGroupsReq struct {
	Op
	UserId string `json:"user_id" binding:"required"`
}

type JoinGroupsResp struct{}

type GroupSearchReq struct {
	OperationID string `json:"operation_id" form:"operation_id" binding:"required" msg:"操作ID不能为空"`
	Search      string `json:"search" form:"search" binding:"required,gte=1"`
}

type GroupSearchResq struct {
	Count int         `json:"count"`
	List  []GroupInfo `json:"list"`
}
