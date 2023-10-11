package model

import (
	"im/pkg/db"
)

type OnlineStatusType string

const (
	OnlineStatusTypeOnline  OnlineStatusType = "Online"
	OnlineStatusTypeOffline OnlineStatusType = "Offline"
)

type User struct {
	db.CommonModel
	UserID           string           `gorm:"column:user_id;size:80;uniqueIndex;default:'';" json:"user_id"`
	Account          string           `gorm:"column:account;size:80;index;default:'';" json:"account"`
	Password         string           `gorm:"column:password;size:255;default:'';" json:"password"`
	PhoneNumber      string           `gorm:"column:phone_number;index;size:20;default:'';" json:"phone_number"`
	CountryCode      string           `gorm:"column:country_code;size:10;default:'';" json:"country_code"`
	FaceURL          string           `gorm:"column:face_url;size:255;default:'';" json:"face_url"`
	BigFaceURL       string           `gorm:"column:big_face_url;size:255;default:'';" json:"big_face_url"`
	Gender           int64            `gorm:"column:gender;default:1;" json:"gender"`
	Platform         int64            `gorm:"column:platform;size:255;default:2;" json:"platform"`
	DeviceId         string           `gorm:"column:device_id;size:255;default:'';" json:"device_id"`
	NickName         string           `gorm:"column:nick_name;size:255;unique;default:'';" json:"nick_name"`
	Signatures       string           `gorm:"column:signatures;size:255;default:'';" json:"signatures"`
	Age              int64            `gorm:"column:age;default:18;" json:"age"`
	LoginIp          string           `gorm:"column:login_ip;size:50;default:'';" json:"login_ip"`
	UserModel        int64            `gorm:"column:user_model;default:1;" json:"user_model"`
	PayPassword      string           `gorm:"column:pay_password;size:80;default:'';" json:"pay_password"`
	Salt             string           `gorm:"column:salt;size:50;default:'';" json:"salt"`
	InviteCode       string           `gorm:"column:invite_code;size:80;" json:"invite_code"`
	Status           int64            `gorm:"column:status;default:1;" json:"status"`
	Balance          int64            `gorm:"column:balance;" json:"balance"`
	OnlineStatus     OnlineStatusType `gorm:"column:online_status;size:20;" json:"online_status"`
	LatestLoginTime  int64            `gorm:"column:latest_login_time;" json:"latest_login_time"`
	ImSite           string           `gorm:"column:im_site;size:20;default:'im';" json:"im_site"`
	IsPrivilege      int64            `gorm:"column:is_privilege;size:20;default:2;" json:"is_privilege"`
	IsCustomer       int64            `gorm:"column:is_customer;size:20;default:2;" json:"is_customer"`
	RegisterIp       string           `gorm:"column:register_ip;size:50;default:'';" json:"register_ip"`
	RegisterDeviceId string           `gorm:"column:register_device_id;size:255;default:'';" json:"register_device_id"`
	UserDevices      []UserDevice     `gorm:"foreignKey:UserID;references:UserID"`
	UserIps          []UserIp         `gorm:"foreignKey:UserID;references:UserID"`
	Privacy          string           `gorm:"column:privacy;size:1800;default:'';" json:"privacy"`
	RealName         string           `gorm:"column:real_name;size:80;default:'';" json:"real_name"`
	IDNo             string           `gorm:"column:id_no;size:80;default:'';" json:"id_no"`
	IDFrontImg       string           `gorm:"column:id_front_img;size:500;default:'';" json:"id_front_img"`
	IDBackImg        string           `gorm:"column:id_back_img;size:500;default:'';" json:"id_back_img"`
	IsRealAuth       int              `gorm:"column:is_real_auth;default:1;" json:"is_real_auth"`
	RealAuthMsg      string           `gorm:"column:real_auth_msg;default:'';" json:"real_auth_msg"`
}

func (d *User) TableName() string {
	return "users"
}

type UserDevice struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   string `gorm:"column:user_id;size:32" json:"user_id"`
	DeviceID string `gorm:"column:device_id;size:48;" json:"device_id"`
	Platform int64  `gorm:"column:platform;" json:"platform"`
	Status   uint   `gorm:"column:status" json:"status"`
}

func (d *UserDevice) TableName() string {
	return "user_devices"
}

type UserIp struct {
	ID     uint   `gorm:"primaryKey"`
	UserID string `gorm:"column:user_id;size:32" json:"user_id"`
	Ip     string `gorm:"column:ip;size:50;default:'';" json:"ip"`
	Status uint   `gorm:"column:status" json:"status"`
}

func (d *UserIp) TableName() string {
	return "user_ips"
}

type LoginHistory struct {
	ID        uint   `gorm:"primaryKey"`
	CreatedAt int64  `gorm:"column:created_at;index:idx_created_at" json:"created_at"`
	UserID    string `gorm:"column:user_id;size:32;index:idx_user_id" json:"user_id"`
	Ip        string `gorm:"column:ip;size:50;default:'';" json:"ip"`
	DeviceID  string `gorm:"column:device_id;size:48;" json:"device_id"`
	Platform  int64  `gorm:"column:platform;" json:"platform"`
	Brand     string `gorm:"column:brand;" json:"brand"`
}

func (d *LoginHistory) TableName() string {
	return "login_history"
}

type SignInV2 struct {
	db.CommonModel
	UserID       string `gorm:"column:user_id;size:32;uniqueIndex;default:0;" json:"user_id"`
	Ip           string `gorm:"column:ip;size:50;default:'';" json:"ip"`
	LastTime     int64  `gorm:"column:last_time;default:0;" json:"last_time"`
	ContinueDays int64  `gorm:"column:continue_days;default:0;" json:"continue_days"`
}

func (d *SignInV2) TableName() string {
	return "sign_in_v2"
}

type SignGifts struct {
	db.CommonModel
	Day       int64 `gorm:"column:day;default:0;uniqueIndex;" json:"day"`
	GiftsCoin int64 `gorm:"column:gifts_coin;default:0;" json:"gifts_coin"`
	IsSeries  int64 `gorm:"column:is_series;default:2;" json:"is_series"`
}

func (d *SignGifts) TableName() string {
	return "sign_gifts"
}
