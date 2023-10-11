package model

type ReasonType int

const (
	ReasonTypeNone ReasonType = 0
	ReasonTypeQuit ReasonType = 1
	ReasonTypeKick ReasonType = 2
)

type RoleType string

const (
	RoleTypeOwner RoleType = "owner"
	RoleTypeAdmin RoleType = "admin"
	RoleTypeStaff RoleType = "staff"
	RoleTypeUser  RoleType = "user"
)

type GroupMember struct {
	Id            int64    `gorm:"column:id;primary_key;" json:"id"`
	GroupId       string   `gorm:"column:group_id;index;size:20;" json:"group_id"`
	UserId        string   `gorm:"column:user_id;index;size:20;" json:"user_id"`
	GroupNickName string   `gorm:"column:group_nick_name;type:varchar(255)" json:"group_nick_name"`
	Role          RoleType `gorm:"column:role;type:varchar(255)" json:"role"`
	RoleIndex     int      `gorm:"column:role_index" json:"role_index"`
	MuteEndTime   int64    `gorm:"column:mute_end_time" json:"mute_end_time"`
	JoinType      string   `gorm:"column:join_type;type:varchar(255)" json:"join_type"`
	JoinByUserId  string   `gorm:"column:join_by_user_id;type:varchar(255)" json:"join_by_user_id"`
	CreateTime    int64    `gorm:"column:create_time" json:"create_time"`
	Version       int      `gorm:"column:version;index;type:bigint" json:"version"`
	Status        int      `gorm:"column:status;default:1" json:"status"`
}

type GroupMemberApply struct {
	Id             int64  `gorm:"column:id;primary_key;" json:"id"`
	GroupId        string `gorm:"column:group_id;index;size:20;" json:"group_id"`
	UserId         string `gorm:"column:user_id;index;size:20;" json:"user_id"`
	Remark         string `gorm:"column:remark;type:varchar(300)" json:"remark"`
	CreateTime     int64  `gorm:"column:create_time" json:"create_time"`
	Status         int    `gorm:"column:status;default:0;type:tinyint(3)" json:"status"`
	OperationTime  int64  `gorm:"operator_time" json:"operation_time"`
	OperatorUserId string `gorm:"operator_user_id;type:varchar(255)" json:"operator_user_id"`
}
