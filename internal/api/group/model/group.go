package model

import "time"

const InviteGroupExpire = 10 * time.Minute

type Group struct {
	Id                         int64  `gorm:"column:id;autoIncrement;uniqueIndex;" json:"id"`
	GroupId                    string `gorm:"column:group_id;size:20;uniqueIndex;primary_key;" json:"group_id"`
	Name                       string `gorm:"column:name;type:varchar(255);index" json:"name"`
	FaceUrl                    string `gorm:"column:face_url;type:varchar(255)" json:"face_url"`
	BigFaceUrl                 string `gorm:"column:big_face_url;type:varchar(255)" json:"big_face_url"`
	MembersTotal               int    `gorm:"column:members_total;type:int(11)" json:"members_total"`
	RobotTotal                 int    `gorm:"column:robot_total;type:int(11)" json:"robot_total"`
	AdminsTotal                int    `gorm:"column:admins_total;type:int(11)" json:"admins_total"`
	Notification               string `gorm:"column:notification;type:varchar(1500)"  json:"notification"`
	Introduction               string `gorm:"column:introduction;type:varchar(500)" json:"introduction"`
	CreateTime                 int64  `gorm:"column:create_time" json:"create_time"`
	UpdatedAt                  int64  `gorm:"column:updated_at" json:"updated_at"`
	LastVersion                int    `gorm:"last_version;type:int(11)" json:"last_version"`
	LastMemberVersion          int    `gorm:"last_member_version;type:int(11)" json:"last_member_version"`
	CreateUserId               string `gorm:"column:create_user_id;type:varchar(255)" json:"create_user_id"`
	Status                     int    `gorm:"column:status;default:1;type:tinyint(3);" json:"status"`
	NoShowNormalMember         int    `gorm:"column:no_show_normal_member;default:2;type:tinyint(3)" json:"no_show_normal_member"`
	NoShowAllMember            int    `gorm:"column:no_show_all_member;default:2;type:tinyint(3)" json:"no_show_all_member"`
	ShowQrcodeByNormalMember   int    `gorm:"column:show_qrcode_by_normal_member;default:1;type:tinyint(3);" json:"show_qrcode_by_normal_member"`
	ShowQrcodeByNormalMemberV2 int    `gorm:"column:show_qrcode_by_normal_member_v2;default:1;type:tinyint(3);" json:"show_qrcode_by_normal_member_v2"`
	JoinNeedApply              int    `gorm:"column:join_need_apply;default:1;type:tinyint(3)" json:"join_need_apply"`
	BanRemoveByNormalMember    int    `gorm:"column:ban_remove_by_normal_member;default:2;type:tinyint(3)" json:"ban_remove_by_normal"`
	MuteAllMember              int    `gorm:"column:mute_all_member;default:2;type:tinyint(3)" json:"mute_all_member"`
	IsDefault                  int    `gorm:"column:is_default;default:2;type:tinyint(3)" json:"is_default"`
	IsTopannocuncement         int    `gorm:"column:is_topannocuncement;default:2;type:tinyint(3)" json:"is_topannocuncement"`
	IsOpenAdminList            int    `gorm:"column:is_open_admin_list;default:1;type:tinyint(3)" json:"is_open_admin_list"`
	IsOpenAdminIcon            int    `gorm:"column:is_open_admin_icon;default:2;type:tinyint(3)" json:"is_open_admin_icon"`
	IsOpenGroupId              int    `gorm:"column:is_open_group_id;default:2;type:tinyint(3)" json:"is_open_group_id"`
	GroupSendLimit             int    `gorm:"column:group_send_limit;default:0;type:int(11)" json:"group_send_limit"`
	IsDisplayNicknameOpen      int    `gorm:"column:is_display_nickname_open;default:1;type:int(11);" json:"is_display_nickname_open"`
	MuteAllPeriod              string `gorm:"column:mute_all_period;type:varchar(255)" json:"mute_all_period"`
}

func (d *Group) TableName() string {
	return "groups"
}
