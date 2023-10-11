package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	UserID            string     `gorm:"column:user_id;size:64;index:user_id,unique" json:"user_id"`
	Nickname          string     `gorm:"column:nick_name;size:255" json:"nick_name"`
	Username          string     `gorm:"column:username;size:50;index:username,unique" json:"username"`
	PhoneNumber       string     `gorm:"column:phone_number;size:32" json:"phone_number"`
	Password          string     `gorm:"column:password;size:255" json:"password"`
	Google2fSecretKey string     `gorm:"column:google_2f_secret_key;size:255" json:"google_2f_secret_key"`
	Salt              string     `gorm:"column:salt" json:"salt"`
	Status            int        `gorm:"column:status;size:1;default:1;comment:'1 opened 2 banned'"`
	Role              int        `gorm:"column:role;size:1;default:1;" json:"role"`
	CreateUser        string     `gorm:"column:create_user"`
	CreateTime        int64      `gorm:"column:create_time"`
	UpdateUser        string     `gorm:"column:update_user"`
	UpdateTime        int64      `gorm:"column:update_time"`
	DeleteUser        string     `gorm:"column:delete_user"`
	DeleteTime        int64      `gorm:"column:delete_time;index:delete_time"`
	AppMangerLevel    int32      `gorm:"column:app_manger_level"`
	GlobalRecvMsgOpt  int32      `gorm:"column:global_recv_msg_opt"`
	TwoFactorEnabled  bool       `gorm:"column:two_factor_enabled;default:1"`
	User2FAuthEnable  bool       `gorm:"column:user_two_factor_control_status;default:1"`
	LastloginIp       string     `gorm:"column:last_login_ip" json:"last_login_ip"`
	LastloginTime     *time.Time `gorm:"column:last_login_time" json:"last_login_time"`
	Remark            string     `gorm:"column:remark"`
	RoleKey           string     `gorm:"-" json:"role_key"`
}

func (Admin) TableName() string {
	return "cms_admins"
}
